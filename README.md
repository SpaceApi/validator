# Validation server for SpaceAPI endpoints

Written in Python 3 with [Bottle](http://bottlepy.org/).

[![Travis CI][travis-ci-badge]][travis-ci]

# Dev setup

Create a virtual environment:

    python3 -m venv VIRTUAL

Enable virtualenv:

    source VIRTUAL/bin/activate

Install dependencies:

    pip install -r requirements.txt
    pip install -r requirements-dev.txt

Start the server:

    cd validator
    python service.py

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
[travis-ci]: https://travis-ci.org/spacedirectory/validator
[travis-ci-badge]: https://img.shields.io/travis/spacedirectory/validator/master.svg
