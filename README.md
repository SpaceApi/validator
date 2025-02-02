# Validation server for SpaceAPI endpoints

OpenAPI spec: https://validator.spaceapi.io/openapi.json

[![CircleCI][circle-ci-badge]][circle-ci]
[![Docker Image][docker-image-badge]][docker-image]
[![Go Report Card][go-report-card-badge]][go-report-card]


# API

There are two main endpoints, to validate raw JSON and to validate URLs:

- https://validator.spaceapi.io/v2/validateURL
- https://validator.spaceapi.io/v2/validateJSON

The full API specification in OpenAPI format can be found at https://validator.spaceapi.io/openapi.json.

## Validating URLs

Use this if your endpoint is already online.

Example (curl):

    curl -X POST -H "Content-Type: application/json" \
        https://validator.spaceapi.io/v2/validateURL \
        -d'{"url": "https://status.crdmp.ch/"}'

Example (httpie):

    http post \
        https://validator.spaceapi.io/v2/validateURL \
        url=https://status.crdmp.ch/

Response:

    {
        "valid": true,
        "message": "",
        "isHttps": true,
        "httpsForward": false,
        "reachable": true,
        "cors": true,
        "contentType": true,
        "certValid": true,
        "validatedJson": { … },
        "schemaErrors": [ … ]
    }

## Validating JSON

If you want to validate JSON data directly, use this endpoint. However, in
contrast to the URL endpoint, only the content will be validated, but not the
server configuration (e.g. whether CORS is set up properly or whether a valid
certificate is being used).

Example (curl):

    curl -X POST -H "Content-Type: application/json" \
        https://validator.spaceapi.io/v2/validateJSON \
        -d @mydata.json

Example (httpie):

    cat mydata.json | http post https://validator.spaceapi.io/v2/validateJSON

Response:

    {
        "message": "",
        "valid": true,
        "validatedJson": { … },
        "schemaErrors": [ … ]
    }


# Dev setup

See `DEVELOPMENT.md`.


# License

Licensed under either of

 * Apache License, Version 2.0, ([LICENSE-APACHE](LICENSE-APACHE) or
   http://www.apache.org/licenses/LICENSE-2.0)
 * MIT license ([LICENSE-MIT](LICENSE-MIT) or
   http://opensource.org/licenses/MIT)

at your option.

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you, as defined in the Apache-2.0 license, shall
be dual licensed as above, without any additional terms or conditions.


<!-- Badges -->
[circle-ci]: https://circleci.com/gh/SpaceApi/validator/tree/master
[circle-ci-badge]: https://circleci.com/gh/SpaceApi/validator/tree/master.svg?style=shield
[docker-image]: https://github.com/SpaceApi/validator/pkgs/container/validator
[docker-image-badge]: https://img.shields.io/badge/container%20image-ghcr.io/spaceapi/validator-blue.svg
[go-report-card]: https://goreportcard.com/report/github.com/spaceapi/validator
[go-report-card-badge]: https://goreportcard.com/badge/github.com/spaceapi/validator
