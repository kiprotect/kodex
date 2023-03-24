package ui

import (
	. "github.com/kiprotect/gospel"
)

func Tabs(args ...any) Element {
	return Div(
		Class("bulma-tabs"),
		Class("active"),
		Span(
			Class("bulma-more"),
			Span("cm-tabs-more", "&or;"),
		),
		Ul(
			args,
		),
	)
}

func ActiveTab(active bool) Attribute {
	if active {
		return Class("bulma-is-active")
	}
	return nil
}

func Tab(args ...any) Element {
	return Li(
		args,
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
