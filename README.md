# Firetail AppSync Lambda

[![License: LGPL v3](https://img.shields.io/badge/License-LGPL_v3-blue.svg)](https://www.gnu.org/licenses/lgpl-3.0) ![Test and coverage](https://github.com/FireTail-io/firetail-appsync-lambda/actions/workflows/codecov.yml/badge.svg?branch=main) [![codecov](https://codecov.io/gh/FireTail-io/firetail-appsync-lambda/branch/main/graph/badge.svg?token=GEPKMSC5ID)](https://codecov.io/gh/FireTail-io/firetail-appsync-lambda)

The Firetail AppSync Lambda recieves Cloudwatch log events from AppSync APIs, and forwards request and response data to the Firetail SaaS.



## Tests

Tests can be run with the standard `go test` command. [stretchr/testify](https://github.com/stretchr/testify) has been used for shorthand assertions. 

A [Makefile](./Makefile) is provided which is used by [this GitHub action](./.github/workflows/codecov.yml) to upload coverage reports to Codecov. The Makefile can be used to generate the same coverage report locally, which you can then view in your browser using `go tool cover`:

```bash
make test
go tool cover -html coverage.out
```



## Deployment

In order to use the Firetail AppSync lambda, an token for the Firetail Logs API is required. Documentation on how to obtain a Firetail Logs API token from the Firetail SaaS can be found [here](). #TODO

The Firetail AppSync Lambda is written in Go, and can be built using the standard `go build` command. First, clone the repository and change directory into `logs-handler`, where the Lambda's source is located:

```bash
git clone git@github.com:FireTail-io/firetail-appsync-lambda.git
cd firetail-appsync-lambda/logs-hander
```

Before building the source into a binary, set `GOARCH` to `amd64` and `GOOS` to `linux` to ensure the binary will be compatible with the [Lambda go runtime](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html):

```bash
GOARCH=amd64 GOOS=linux
```

Next, build the binary and output it into a `bin` directory at the root of the repository. `-ldflags="-s -w"` can be used to marginally reduce the size of the binary:

```bash
go build -ldflags="-s -w" -o ../bin/logs-handler
```

A [serverless.yml](./serverless.yml) is provided in the root of this repository which can be used to deploy this binary to Lambda, and expects the binary to be found in a `bin` directory at the root of the repository, hence `-o ../bin/logs-handler`.

The provided [serverless.yml](./serverless.yml) has two parameters:

1. `firetail-api-token`, an API token for the Firetail Logs API.
2. `cloudwatch-log-group`, the log group for an AppSync API in Cloudwatch.

Given these two values, the Lambda can be deployed by running the following serverless command from the root of the repository:

```bash
cd ..
sls deploy --param="firetail-api-token=YOUR_FIRETAIL_API_TOKEN" --param="cloudwatch-log-group=YOUR_CLOUDWATCH_LOG_GROUP"
```

This serverless command may require additional flags depending upon the use case, for example to specify the region in which the Lambda should be deployed. See `sls deploy --help` for a list of available flags.
