# Validation server for SpaceAPI endpoints

https://validator.spaceapi.io/v2/

[![CircleCI][circle-ci-badge]][circle-ci]
[![Docker Image][docker-image-badge]][docker-image]


# Dev setup

## Dependencies

Install golang

Then:
```bash
go get -d  ./...
go generate
go install  ./...
```

## Starting the Server

Start the server:

    validator

## Testing

To run tests:

    go test ./...


# API

## Request

To send a validation request, send a POST request to `/v1/validate/` with
`Content-Type: application/json`. The payload (in JSON format) should look like
this:

```javascript
{
    "data": "..."
}
```

The `data` field should contain the SpaceAPI endpoint data as a JSON string.

Example (curl):

    curl \
        -X POST \
        -H "Content-Type: application/json" \
        https://validator.spaceapi.io/v1/validate/ \
        -d'{"data": "{\"api\": \"0.13\"}"}'

Example (httpie):

    http POST \
        https://validator.spaceapi.io/v1/validate/ \
        data='{"api": "0.13"}'

## Response

If the request is not malformed, the endpoint returns a HTTP 200 response with
`Content-Type: application/json`.

The success response looks like this:

```javascript
{
    "valid": true,
    "message": null
}
```

The error response looks like this:

```javascript
{
    "valid": false,
    "message": "Error details"
}
```

It is planned that more error details (like row/col) will be added in the future.


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
[docker-image]: https://hub.docker.com/r/spaceapi/validator/
[docker-image-badge]: https://img.shields.io/docker/pulls/spaceapi/validator.svg
