## Building The Firetail AppSync Lambda From Source

The Firetail AppSync Lambda is written in Go, and can be built using the standard `go build` command. First, clone the repository and change directory into `logs-handler`, where the Lambda's source is located:

```bash
git clone git@github.com:FireTail-io/firetail-appsync-lambda.git
cd firetail-appsync-lambda/logs-hander
```

Before building the source into a binary, set `GOARCH` to `amd64` and `GOOS` to `linux` to ensure the binary will be compatible with the [Lambda Go runtime](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html):

```bash
export GOARCH=amd64 GOOS=linux
```

Next, build the binary and output it into a `bin` directory at the root of the repository. `-ldflags="-s -w"` can be used to marginally reduce the size of the binary:

```bash
go build -ldflags="-s -w" -o ../bin/logs-handler
```

The [serverless.yml](./serverless.yml) provided in the root of this repository can be used to deploy this binary to Lambda, and expects the binary to be found in a `bin` directory at the root of the repository, hence `-o ../bin/logs-handler`.
