# Fiber Modules

A list of modules useful for using with [fiber](https://github.com/gofiber/fiber)

## Install

```bash
go get github.com/gkampitakis/fiber-modules
```

## Contents

- [healthcheck](./healthcheck/README.md)
- [gracefulshutdown](./gracefulshutdown/README.md)

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
go mod edit -replace=github.com/gkampitakis/fiber-modules@v0.0.1-beta=../gofiber-custom-healthcheck

go get -d github.com/gkampitakis/fiber-modules@v0.0.1-beta
```

With this way we replace the module with the local instance. Hope this helps.

## License

MIT License