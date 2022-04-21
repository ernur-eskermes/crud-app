# CRUD for Book Management [Backend Application] ![GO][go-badge]

[go-badge]: https://img.shields.io/github/go-mod/go-version/p12s/furniture-store?style=plastic
[go-url]: https://github.com/p12s/furniture-store/blob/master/go.mod

## Build & Run (Locally)
### Prerequisites
- go 1.17
- postgres & golang-migrate
- docker & docker-compose
- [golangci-lint](https://github.com/golangci/golangci-lint) (<i>optional</i>, used to run code checks)
- [swag](https://github.com/swaggo/swag) (<i>optional</i>, used to re-generate swagger documentation)

Create .env file in root directory and add following values:
```dotenv
POSTGRES_URI=postgresql://postgres:qwerty123@ninja-db/postgres

JWT_SIGNING_KEY=JUHDSGYUiAUBUGIFHOJPJF*($#O@J(*FU*#!)(J#
PASSWORD_SALT=JUHDSGYUiAUBUGIFHOJPJF*($#O@J(*FU*#!)(J#
SESSION_SECRET=JUHDSGYUiAUBUGIFHOJPJF*($#O@J(*FU*#!)(J#

HTTP_HOST=localhost

APP_ENV=local

GRPC_AUDIT_URL=host_ip:9000
```

Use `make run` to build&run project, `make lint` to check code with linter, `make migrate` to apply the migration scheme.