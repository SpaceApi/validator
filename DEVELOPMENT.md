# Development

## Dependencies

Install golang

Then:

    go get -d ./...
    go generate

## Run from code

Start the server:

    go run ./...

## Build a binary

Install the binary to `$GOPATH/bin/validator`:

    go install ./...

...then run it:

    validator

## Testing

To run tests:

    go test ./...

## Using local validator library

If you work on this project, you'll probably make changes to the
"go-spaceapi-validator" library. To use the local version instead of the one
published on GitHub, add the following to your `go.mod` (with the appropriate
path):

    replace "github.com/spaceapi-community/go-spaceapi-validator" => "../go-spaceapi-validator"
