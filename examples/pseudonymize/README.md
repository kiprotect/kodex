# Pseudonymization Demo

This demo shows how to pseudonymize structured data with KIProtect.
To pseudonymize the example data in `input.json`, simply run

    kiprotect run pseudonymize

To depseudonymize it again, simply run

    kiprotect run depseudonymize

# Pseudonymization with custom key

The two examples above will use the parameter management of KIProtect to
store pseudonymization keys. If you want to specify the keys yourself you
can do so as well:

    kiprotect run pseudonymize-with-key

and

    kiprotect run depseudonymize-with-key

KIProtect will ask you to enter a key and use that key to pseudonymize
and depseudonymize the items, without reyling on the parameter management.