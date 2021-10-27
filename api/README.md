# KIProtect - KIProtect API

KIProtect provides an API to the KIProtect tool, allowing users to
transform data via a RESTful API.

# Building

To build the API, run the following commands in the parent directory:

    make
    make install

# Running the tests

Run the following in the parent directory:

    make test-api

# Running the API command

First, define the `KODEX_SETTINGS` environment variable. Then, run

    kodex api run
