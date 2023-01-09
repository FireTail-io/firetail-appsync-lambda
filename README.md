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



## Setup Guide

1. [Install prerequisites](#installing-prerequisites)
2. [Configure the AppSync app's Cloudwatch logs](#configuring-appsync)
3. [Generate a Firetail API token](#generating-a-firetail-api-token)
4. [Build the Firetail AppSync Lambda](#building-the-firetail-appsync-lambda)
5. [Deploy the Firetail AppSync Lambda with Serverless](#deploying-the-firetail-appsync-lambda-with-serverless)



### Installing prerequisites

Building and deploying the Firetail AppSync Lambda requires the following prerequisites:

- A Go installation - see the [Download and Install instructions at go.dev](https://go.dev/doc/install)
- The Serverless CLI - see [Setting Up Serverless Framework With AWS at serverless.com](https://www.serverless.com/framework/docs/getting-started)
- An account and organisation on the Firetail SaaS - visit [firetail.app](https://firetail.app) to get started



### Configuring AppSync

Before deploying the Firetail AppSync Lambda, ensure the AppSync app is configured to log to Cloudwatch. The corresponding [AWS documentation can be found here](https://docs.aws.amazon.com/appsync/latest/devguide/monitoring.html).

It is recommended to:

- Enable **Include verbose content**
- Set the **Field resolver log level** to **All**

📝 Take note of the name of the Cloudwatch log group for the AppSync app, as this will be required when deploying the Firetail AppSync Lambda. It should be of the format `/aws/appsync/apis/{graphql_api_id}`.



### Generating a Firetail API Token

In order to deploy the Firetail AppSync Lambda, an API token from the Firetail SaaS is required. Documentation on how to obtain a Firetail API token from the Firetail SaaS can be found [here](). #TODO

📝 Take note of the Firetail API token, as this will be required when deploying the Firetail AppSync Lambda.



### Building the Firetail AppSync Lambda

The Firetail AppSync Lambda is written in Go, and can be built using the standard `go build` command. First, clone the repository and change directory into `logs-handler`, where the Lambda's source is located:

```bash
git clone git@github.com:FireTail-io/firetail-appsync-lambda.git
cd firetail-appsync-lambda/logs-hander
```

Before building the source into a binary, set `GOARCH` to `amd64` and `GOOS` to `linux` to ensure the binary will be compatible with the [Lambda Go runtime](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html):

```bash
GOARCH=amd64 GOOS=linux
```

Next, build the binary and output it into a `bin` directory at the root of the repository. `-ldflags="-s -w"` can be used to marginally reduce the size of the binary:

```bash
go build -ldflags="-s -w" -o ../bin/logs-handler
```

The [serverless.yml](./serverless.yml) provided in the root of this repository can be used to deploy this binary to Lambda, and expects the binary to be found in a `bin` directory at the root of the repository, hence `-o ../bin/logs-handler`.



### Deploying the Firetail AppSync Lambda with Serverless

A [serverless.yml](./serverless.yml) is provided in the root of this repository, which has two parameters:

1. `cloudwatch-log-group`, the log group for an AppSync API in Cloudwatch (see [Configuring AppSync📝](#configuring-appsync))
2. `firetail-api-token`, an API token from the Firetail SaaS (see [Generating a Firetail API Token📝](#generating-a-firetail-api-token))

Given these two values, the Lambda can be deployed by running the following serverless command from the root of the repository:

```bash
sls deploy --param="cloudwatch-log-group=YOUR_CLOUDWATCH_LOG_GROUP" --param="firetail-api-token=YOUR_FIRETAIL_API_TOKEN"
```

This serverless command may require additional flags depending upon the use case, for example to specify the region in which the Lambda should be deployed. See `sls deploy --help` for a list of available flags.
