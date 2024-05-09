# zula-service

Written in golang, uses a swagger file to generate code for the API. All logic underneath is found in /service folder. Logging, env variables and db access is found in the infrastructure.

## Development

`docker-compose up service` runs it, it uses air for hot reloading

## How to generate code from swagger

- `go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen -package api -generate types,server,spec docs/swagger.yaml > api/api.gen.go`

## Data

- User information
- Notes
- User activity
