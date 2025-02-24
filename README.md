# Terraform Provider Scaffolding (Terraform Plugin Framework)

_This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding). See [Which SDK Should I Use?](https://developer.hashicorp.com/terraform/plugin/framework-benefits) in the Terraform documentation for additional information._

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, do the following:

```shell
# start the backend service
docker compose -f docker_compose/docker-compose.yml up -d
# check if the service is ready
curl -s -X GET "http://localhost:8080/ready"
# run acceptance tests once the service is ready
make testacc
```
