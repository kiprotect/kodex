# Welcome to KIProtect!

KIProtect is an open-source toolkit for **privacy and security engineering**. It
helps you to automate data security and data protection measures in your
data engineering workflows. It offers the following functionality:

- Read data items from a variety of sources such as files, databases or
  message queues.
- Protect these data items using various privacy- & security enhancing
  transformations, like de-identification, masking, pseudonymization,
  anonymization or encryption.
- Send the protected items to a variety of destinations.

With KIProtect, you can describe your data protection and data security
workflows using a simple, declarative configuration language: Just like DevOps
tools let you describe infrastructure as code, KIProtect is a **PrivacyOps** /
**SecurityOps** tool that let you describe **privacy and security measures
as code**.

KIProtect takes care of the boring and difficult aspects of privacy, such as

- **Key management**: KIProtect manages encryption and pseudonymization
  keys for you (if you want that).
- **Parameter management**: KIProtect keeps track of how every single data item
  was processed so you can prove the compliance of your data workflows
  and create an audit trail.
- **Data transformation**: KIProtect implements modern cryptographic and
  statistical techniques to protect your data.

# Getting started

To download and install KIProtect from source, simply run

    git clone https://github.com/kiprotect/kiprotect
    cd kiprotect

    make
    make install

# Documentation

You can find the official documentation at https://kiprotect.com/docs.

# Transforming data

KIProtect reads its configuration from so-called blueprints. To get an idea
of how this works, check out the `docs/examples` directory, which contains
example blueprints together with instructions on how to run them. For example,
the pseudonymization example can be run as follows:

    cd docs/examples/pseudonymize
    # pseudonymize the example data in input.json
    kiprotect run pseudonymize
    # depseudonymize the data again
    kiprotect run depseudonymize

That's it! KIProtect takes care of generating and storing cryptographic
parameters for the pseudonymization. If you want to manually enter a key instead
to generate parameters, you can do that too:

    # pseudonymize the data with a user-supplied key
    kiprotect run pseudonymize-with-key
    # depseudonymize with a key as well
    kiprotect run depseudonymize-with-key

# Running the tests

KIProtect comes with a suite of automated unit tests, which you can run with
Make:

    make test

# Running the benchmarks

KIProtect also comes with a number of benchmarks that you can run as follows:

    make bench

# Status & Roadmap

This is still an early version of KIProtect and does not contain many features
yet. We will progressively port more functionality from our Enterprise Edition
(EE). The following features are next up on our list:

- **Anonymization**: Anonymize streaming data using differentially private
  aggregations.
- **Discovery**: Discover sensitive and personal information in your structured
  and unstructured data.
- **Encryption**: Encrypt and decrypt structured data.
- **Data Mapping**: Analyze and map your data infrastructure.
- **Consent Management**: Manage and enforce processing purposes and
  user consent for all your data streams.

# Enterprise Edition

Our open-source work is made possible by commercially offering a **KIProtect
enterprise edition (EE)**, which extends the community edition (CE) with
functionality that supports a deployment of KIProtect in a professional
enterprise environment. It includes e.g. the following functionality:

- Advanced, SQL-based configuration & parameter management and storage.
- REST-based API to control all KIProtect functionality.
- Web interface to manage and monitor data streams.
- More advanced data transformations.
- Role-based access control mechanism.

Are you interested to learn more about KIProtect EE? Just visit
[our website](https://kiprotect.com) or [get in touch with us](ee@kiprotect.com)!

# License

KIProtect is currently released under the Affero GPL license. See the
[license file](LICENSE) for more information.

# Contact us

Do you have trouble getting KIProtect to work? Do you want to suggest a new
feature or report a bug? Please open an issue in this issue tracker. If
it's something that you'd like to discuss directly with us, please
[send us an e-mail](ce@kiprotect.com), we love to hear from you!

# Spread the word

Are you using KIProtect in your organization and like it? Please let the world
know! Spreading the word about it and giving us feedback helps us to improve
the software.
