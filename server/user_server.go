// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package main

import (
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"path"
	"strings"

	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-jira/server/utils/types"
)

type OAuth1aTemporaryCredentials struct {
	Token  string
	Secret string
}

func (p *Plugin) httpOAuth1aComplete(w http.ResponseWriter, r *http.Request, instanceID types.ID) (status int, err error) {
	// Prettify error output
	defer func() {
		if err == nil {
			return
		}

		errtext := err.Error()
		if len(errtext) > 0 {
			errtext = strings.ToUpper(errtext[:1]) + errtext[1:]
		}
		status, err = p.respondSpecialTemplate(w, "/other/message.html", status, "text/html", struct {
			Header  string
			Message string
		}{
			Header:  "Failed to connect to Jira.",
			Message: errtext,
		})
	}()

	instance, err := p.instanceStore.LoadInstance(instanceID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// APPLINKTODO: OAuth complete should support Jira server and Jira cloud.
	// In this commit, Jira server is not supported because of this line
	si, ok := instance.(*cloudInstance)
	if !ok {
		return http.StatusInternalServerError,
			errors.Errorf("Not supported for instance type %s", instance.Common().Type)
	}

	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
	if err != nil {
		return http.StatusInternalServerError,
			errors.WithMessage(err, "failed to parse callback request from Jira")
	}

	mattermostUserID := r.Header.Get("Mattermost-User-Id")
	if mattermostUserID == "" {
		return http.StatusUnauthorized, errors.New("not authorized")
	}
	mmuser, appErr := p.API.GetUser(mattermostUserID)
	if appErr != nil {
		return http.StatusInternalServerError,
			errors.WithMessage(appErr, "failed to load user "+mattermostUserID)
	}

	oauthTmpCredentials, err := p.otsStore.OneTimeLoadOauth1aTemporaryCredentials(mattermostUserID)
	if err != nil || oauthTmpCredentials == nil || oauthTmpCredentials.Token == "" {
		return http.StatusInternalServerError, errors.WithMessage(err, "failed to get temporary credentials for "+mattermostUserID)
	}

	if oauthTmpCredentials.Token != requestToken {
		return http.StatusUnauthorized, errors.New("request token mismatch")
	}

	// Although we pass the oauthTmpCredentials as required here. The JIRA server does not appar to validate it.
	// We perform the check above for reuse so this is irrelevant to the security from our end.
	accessToken, accessSecret, err := si.getOAuth1Config().AccessToken(requestToken, oauthTmpCredentials.Secret, verifier)
	if err != nil {
		return http.StatusInternalServerError,
			errors.WithMessage(err, "failed to obtain oauth1 access token")
	}

	connection := &Connection{
		PluginVersion:      manifest.Version,
		Oauth1AccessToken:  accessToken,
		Oauth1AccessSecret: accessSecret,
	}

	client, err := instance.GetClient(connection)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	juser, err := client.GetSelf()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	connection.User = *juser

	// Set default settings the first time a user connects
	connection.Settings = &ConnectionSettings{Notifications: true}

	err = p.connectUser(instance, types.ID(mattermostUserID), connection)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return p.respondTemplate(w, r, "text/html", struct {
		MattermostDisplayName string
		JiraDisplayName       string
		RevokeURL             string
	}{
		JiraDisplayName:       juser.DisplayName + " (" + juser.Name + ")",
		MattermostDisplayName: mmuser.GetDisplayName(model.SHOW_NICKNAME_FULLNAME),
		RevokeURL:             path.Join(p.GetPluginURLPath(), instancePath(routeUserDisconnect, instance.GetID())),
	})
}

func (p *Plugin) httpOAuth1aDisconnect(w http.ResponseWriter, r *http.Request, instanceID types.ID) (int, error) {
	if r.Method != http.MethodGet {
		return respondErr(w, http.StatusMethodNotAllowed,
			errors.New("method "+r.Method+" is not allowed, must be GET"))
	}

	mattermostUserID := r.Header.Get("Mattermost-User-Id")
	if mattermostUserID == "" {
		return respondErr(w, http.StatusUnauthorized, errors.New("not authorized"))
	}

	_, err := p.DisconnectUser(instanceID.String(), types.ID(mattermostUserID))
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}

	return p.respondSpecialTemplate(w, "/other/message.html", http.StatusOK,
		"text/html", struct {
			Header  string
			Message string
		}{
			Header:  "Disconnected from Jira.",
			Message: "It is now safe to close this browser window.",
		})
}

func publicKeyString(p *Plugin) ([]byte, error) {
	rsaKey := p.getConfig().rsaKey
	b, err := x509.MarshalPKIXPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to encode public key")
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}), nil
}
