# Prefix all classes:

Use the following regexes:

    /\.((?!bulma)[a-z]+(?:-[0-9a-z]+)*)(?!\w*"\n)/.bulma-$1/g
    /\$((?!bulma)[a-z]+[^\s]*)/$bulma-$1/g
    /\=([a-z\-0-9]+)/=bulma-$1/g
    /\+([a-z\-0-9]+)/+bulma-$1/g    

This will prepend a "bulma-" prefix to all classes, variables and mixins. Bliss :)