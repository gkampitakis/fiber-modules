# Fiber Modules

TBD

## Install

TBD

## Usage


## API

## Example response

## Local Development

This is a general guide when you want to test some changes in a `go module` 
before publishing.

You want to test some changes in `./gofiber-custom-healthcheck` in a 
`./server`.

Inside `./server`

```bash
# Assuming that your path structure is
# ../gofiber-custom-healthcheck
# ../server
go mod edit -replace=github.com/gkampitakis/gofiber-custom-healthcheck@v0.0.1-beta=../gofiber-custom-healthcheck

go get -d github.com/gkampitakis/gofiber-custom-healthcheck@v0.0.1-beta
```

With this way we replace the module with the local instance. Hope this helps.

## License

MIT License