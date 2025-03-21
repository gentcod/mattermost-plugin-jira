// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {changeOpacity} from 'mattermost-redux/utils/theme_utils';
import {Theme} from 'mattermost-redux/types/preferences';

import {ProvidedStyle, StyleState} from 'types/model';

export const getBaseStyles = (theme: Theme) => {
    return {
        codeBlock: ({
            padding: '10px 12px',
            background: changeOpacity(theme.centerChannelColor, 0.08),
            borderRadius: '4px',
            marginTop: '8px',
            marginBottom: '8px',
            fontSize: '13px',
        }),
    };
};

export const getModalStyles = (theme: Theme) => ({
    modalBody: {
        padding: '2em 2em 3em',
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
    },
    modalFooter: {
        padding: '2rem 15px',
    },
    descriptionArea: {
        height: 'auto',
        width: '100%',
        color: '#000',
    },
});

export const getStyleForReactSelect = (theme: Theme) => {
    if (!theme) {
        return {};
    }

    return {
        menuPortal: (provided: ProvidedStyle) => ({
            ...provided,
            zIndex: 9999,
        }),
        control: (provided: ProvidedStyle, state: StyleState) => ({
            ...provided,
            color: theme.centerChannelColor,
            background: theme.centerChannelBg,

            // Overwrittes the different states of border
            borderColor: state.isFocused ? changeOpacity(theme.centerChannelColor, 0.25) : changeOpacity(theme.centerChannelColor, 0.2),
            padding: '2px 4px 2px 6px',

            // Removes weird border around container
            boxShadow: 'inset 0 1px 1px ' + changeOpacity(theme.centerChannelColor, 0.075),
            borderRadius: '4px',

            '&:hover': {
                borderColor: changeOpacity(theme.centerChannelColor, 0.25),
            },
        }),
        option: (provided: ProvidedStyle, state: StyleState) => ({
            ...provided,
            background: state.isFocused ? changeOpacity(theme.centerChannelColor, 0.12) : theme.centerChannelBg,
            color: theme.centerChannelColor,
            '&:hover': {
                background: changeOpacity(theme.centerChannelColor, 0.12),
            },
        }),
        clearIndicator: (provided: ProvidedStyle) => ({
            ...provided,
            width: '34px',
            color: changeOpacity(theme.centerChannelColor, 0.4),
            transform: 'scaleX(1.15)',
            marginRight: '-10px',
            '&:hover': {
                color: theme.centerChannelColor,
            },
        }),
        multiValue: (provided: ProvidedStyle) => ({
            ...provided,
            background: changeOpacity(theme.centerChannelColor, 0.15),
        }),
        multiValueLabel: (provided: ProvidedStyle) => ({
            ...provided,
            color: theme.centerChannelColor,
            paddingBottom: '4px',
            paddingLeft: '8px',
            fontSize: '90%',
        }),
        multiValueRemove: (provided: ProvidedStyle) => ({
            ...provided,
            transform: 'translateX(-2px) scaleX(1.15)',
            color: changeOpacity(theme.centerChannelColor, 0.4),
            '&:hover': {
                background: 'transparent',
            },
        }),
        menu: (provided: ProvidedStyle) => ({
            ...provided,
            color: theme.centerChannelColor,
            background: theme.centerChannelBg,
            border: '1px solid ' + changeOpacity(theme.centerChannelColor, 0.2),
            borderRadius: '0 0 2px 2px',
            boxShadow: changeOpacity(theme.centerChannelColor, 0.2) + ' 1px 3px 12px',
            marginTop: '4px',
        }),
        input: (provided: ProvidedStyle) => ({
            ...provided,
            color: theme.centerChannelColor,
        }),
        placeholder: (provided: ProvidedStyle) => ({
            ...provided,
            color: theme.centerChannelColor,
        }),
        dropdownIndicator: (provided: ProvidedStyle) => ({
            ...provided,

            '&:hover': {
                color: theme.centerChannelColor,
            },
        }),
        singleValue: (provided: ProvidedStyle) => ({
            ...provided,
            color: theme.centerChannelColor,
        }),
        indicatorSeparator: (provided: ProvidedStyle) => ({
            ...provided,
            display: 'none',
        }),
    };
};
