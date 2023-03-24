package ui

import (
	. "github.com/kiprotect/gospel"
)

func Tabs(tabs []TabConfig) Element {
	return Div(
		Class("bulma-tabs"),
		Class("active"),
		Span(
			Class("bulma-more"),
			Span("cm-tabs-more", "&or;"),
		),
		Ul(
			"tabs",
		),
	)
}

type TabConfig struct {
	Name string
}

func Tab(active bool) Element {
	return Li(
		Class("bulma-is-active"),
		"test",
	)
}

/*
export const Tab = ({ active, children, href, icon, params, onClick }) => (
    <li className={active ? 'bulma-is-active' : ''}>
        <A href={href} params={params} onClick={onClick}>
            {icon && <span className="icon is-small">{icon}</span>}
            {children}
        </A>
    </li>
);

Tab.propTypes = {
    active: PropTypes.bool,
    children: PropTypes.node.isRequired,
    href: PropTypes.string,
    icon: PropTypes.node,
    params: PropTypes.object,
    onClick: PropTypes.func,
};
*/
