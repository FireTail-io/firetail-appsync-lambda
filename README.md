# Firetail AppSync Lambda

[![License: LGPL v3](https://img.shields.io/badge/License-LGPL_v3-blue.svg)](https://www.gnu.org/licenses/lgpl-3.0) ![Test and coverage](https://github.com/FireTail-io/firetail-appsync-lambda/actions/workflows/codecov.yml/badge.svg?branch=main) [![codecov](https://codecov.io/gh/FireTail-io/firetail-appsync-lambda/branch/main/graph/badge.svg?token=GEPKMSC5ID)](https://codecov.io/gh/FireTail-io/firetail-appsync-lambda)

The Firetail AppSync Lambda recieves Cloudwatch log events from AppSync APIs, and forwards request and response data to [Firetail.App](https://firetail.app/).



## Tests

Tests can be run with the standard `go test` command. [stretchr/testify](https://github.com/stretchr/testify) has been used for shorthand assertions. 

A [Makefile](./Makefile) is provided which is used by [this GitHub action](./.github/workflows/codecov.yml) to upload coverage reports to Codecov. The Makefile can be used to generate the same coverage report locally, which you can then view in your browser using `go tool cover`:

```bash
make test
go tool cover -html coverage.out
```



## Setup Guide

1. [Install dependencies](#installing-dependencies)
2. [Configure the AppSync app's Cloudwatch logs](#configuring-appsync)
3. [Generate a Firetail API token](#generating-a-firetail-api-token)
4. [Build the Firetail AppSync Lambda](#building-the-firetail-appsync-lambda)
5. [Deploy the Firetail AppSync Lambda with Serverless](#deploying-the-firetail-appsync-lambda-with-serverless)



### Installing Dependencies

Building and deploying the Firetail AppSync Lambda requires the following dependencies:

- A Go installation - see the [Download and Install instructions at go.dev](https://go.dev/doc/install)
- The Serverless CLI - see [Setting Up Serverless Framework With AWS at serverless.com](https://www.serverless.com/framework/docs/getting-started)



### Configuring AppSync

Before deploying the Firetail AppSync Lambda, ensure the AppSync app is configured to log to Cloudwatch. The corresponding [AWS documentation can be found here](https://docs.aws.amazon.com/appsync/latest/devguide/monitoring.html).

It is recommended to:

- Enable **Include verbose content**
- Set the **Field resolver log level** to **All**

📝 Take note of the name of the Cloudwatch log group for the AppSync app, as this will be required when deploying the Firetail AppSync Lambda. It should be of the format `/aws/appsync/apis/{graphql_api_id}`.



### Generating A Firetail API Token

In order to deploy the Firetail AppSync Lambda, an API token from the Firetail SaaS is required. 

Create an account at [Firetail.App](https://firetail.app/), then:

1. If you do not already have an organisation: go to [firetail.app/organisations](https://firetail.app/organisations), click **Create Organisation**, select a plan, and give the organisation a name and description.
2. On the organisation page visit the **Applications** tab\* and click **Create Application**. Give the application a name.
3. On the application page visit the **APIs** tab\*\* and click **Create API**. Give the API a name and set the **API Type** to **GraphQL**. 
4. On the API page visit the **Tokens** tab\*\*\* and click **Create Token**. Give the token a name, then take note of the **Token Secret** as this will be required when deploying the Firetail AppSync Lambda. 📝

\*`https://firetail.app/organisations/your-org-id/applications`

\*\*`https://firetail.app/organisations/your-org-id/applications/your-app-id/apis` 

\*\*\*`https://firetail.app/organisations/your-org-id/applications/your-app-id/apis/your-api-id/tokens`



### Building The Firetail AppSync Lambda

The process of building the Firetail AppSync Lambda binary can be performed using [the Makefile at the root of this repository](./Makefile), using the `build` target:

```bash
git clone git@github.com:FireTail-io/firetail-appsync-lambda.git
cd firetail-appsync-lambda
make build
```

A more in-depth explanation of how to build the Firetail AppSync Lambda from source can be found in [docs/build-from-src.md](./docs/build-from-src.md).



### Deploying The Firetail AppSync Lambda With Serverless

A [serverless.yml](./serverless.yml) is provided in the root of this repository, which has two parameters:

1. `cloudwatch-log-group`, the log group for an AppSync API in Cloudwatch (see [Configuring AppSync📝](#configuring-appsync))
2. `firetail-api-token`, an API token from the Firetail SaaS (see [Generating a Firetail API Token📝](#generating-a-firetail-api-token))

Given these two values, the Lambda can be deployed by running the following serverless command from the root of the repository:

```bash
sls deploy --param="cloudwatch-log-group=YOUR_CLOUDWATCH_LOG_GROUP" --param="firetail-api-token=YOUR_FIRETAIL_API_TOKEN"
```

This serverless command may require additional flags depending upon the use case, for example to specify the region in which the Lambda should be deployed. See `sls deploy --help` for a list of available flags.
