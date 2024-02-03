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
