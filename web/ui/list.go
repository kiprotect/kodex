package ui

import (
	. "github.com/gospel-dev/gospel"
)

func List(args ...any) *HTMLElement {
	return Div(
		Class("kip-list"),
		args,
	)
}

func ListItem(args ...any) *HTMLElement {
	return Div(
		Class("kip-item", "kip-is-card"),
		args,
	)
}

func ListHeader(args ...any) *HTMLElement {
	return Div(
		Class("kip-item", "kip-is-header"),
		args,
	)
}

func ListColumn(size string, args ...any) *HTMLElement {
	return Div(
		Class("kip-col", Fmt("kip-is-%s", size)),
		args,
	)
}

/*
import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import './list.scss';

export const List = ({ children }) => (
    <div className="kip-list">{children}</div>
);

export const ListHeader = ({ children }) => (
    <div className="kip-item kip-is-header">{children}</div>
);

export const ListColumn = ({ children, size = 'md', wraps = false }) => (
    <div
        className={classnames(`kip-col kip-is-${size}`, { 'kip-wraps': wraps })}
    >
        {children}
    </div>
);

export const ListItem = ({ children, isCard = true, onClick }) => (
    <div
        // Make focusable with the keyboard, if a handler is available
        tabIndex={onClick ? 0 : -1}
        className={classnames('kip-item', {
            'kip-is-card': isCard,
            'kip-is-clickable': onClick,
        })}
        onClick={e => {
            e.preventDefault();
            if (onClick) onClick();
        }}
    >
        {children}
    </div>
);

ListItem.propTypes = {
    children: PropTypes.node,
    isCard: PropTypes.bool,
    onClick: PropTypes.func,
};
*/
