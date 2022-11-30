# Firetail AppSync Lambda

[![License: LGPL v3](https://img.shields.io/badge/License-LGPL_v3-blue.svg)](https://www.gnu.org/licenses/lgpl-3.0) [![Test and coverage](https://github.com/FireTail-io/firetail-appsync-lambda/actions/workflows/codecov.yml/badge.svg?branch=main)] [![codecov](https://codecov.io/gh/FireTail-io/firetail-appsync-lambda/branch/main/graph/badge.svg?token=GEPKMSC5ID)](https://codecov.io/gh/FireTail-io/firetail-appsync-lambda)

The Firetail AppSync Lambda is intended to recieve Cloudwatch log events regarding AppSync APIs, and forward request and response data to the Firetail logs API.



## Tests

Tests can be run with the standard `go test` command. [stretchr/testify](https://github.com/stretchr/testify) has been used for shorthand assertions. You may use the provided [Makefile](./Makefile) to generate a coverage report and view it in your browser using `go tool cover`:

```bash
make test
go tool cover -html coverage.out
```

Test coverate is uploaded to Codecov via [this GitHub action](./.github/workflows/codecov.yml).



## Deployment

The Firetail AppSync Lambda is written in Go, and can be built using the standard `go build` command. A [Makefile](./Makefile) is provided which will build the binary and place it in a `bin` directory as per the [serverless.yml](./serverless.yml)'s `handler` property for the `logs-handler` function.

The provided [serverless.yml](./serverless.yml) has two parameters:

1. `firetail-api-key`, your API key for the Firetail Logs API.
2. `cloudwatch-log-group`, the log group for your AppSync API in Cloudwatch.

Once you have these values, and the binary is built, you should be able to do:

```bash
sls deploy --param="firetail-api-key=YOUR_FIRETAIL_API_KEY" --param="cloudwatch-log-group=YOUR_CLOUDWATCH_LOG_GROUP"
```

