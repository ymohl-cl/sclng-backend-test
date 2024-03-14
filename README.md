# Backend Technical Test at Scalingo

## Instructions

You can find the project description and the requirements on `JOIN`.

Need to use the public [api rest of Github](https://docs.github.com/en/rest?apiVersion=2022-11-28) to get some data.

The project need to be fast and usable by a lot of users simultaneously.

## Implementation

The github crawling logic is implemented in the gitcrawler package `pkg/gitcrawler`.

Git do not provide a way to search repositories by date. So to have the newest created repositories:

- we need to get the most recent event of creation repository
- Extract their identifier
- Substract to the identifier the max element response by the API (100)
- Get the public repository since this new calculated identifier

The repository enrichment is executed in parallel to satisfy the requirement of the project, but it could be limited by the rate limit of the github API and not recommended by their best practices.

## Execution

Set environments variables:

``` bash
export GITHUB_TOKEN=your_github_token
export PORT=5000
```

``` bash
docker compose up
```

Application will be then running on port `5000`

## Makefile

Makefile is provided to simplify the usage of the project.

``` bash
## Install the tooling needs for run the project (ci context):
make ci-tool
## Install tool needs for the development project after calling the ci-tool rule
make tool
## Run linter (golangci-lint) on the full code base
make lint
## Update the generated resources (mocks, swagger, ...)
make update
## Build the project
make build
```

## Test

``` bash
$ curl localhost:5000/ping
{ "status": "pong" }
```

``` bash
# Get repositories without filter
$ curl localhost:5000/repos
{ "repositories": [] }
```

``` bash
# Get repositories with filter
$ curl localhost:5000/repos?language=go&license=mit
{ "repositories": [] }
```

## To go further

Here, some ideas / suggections to have an undustrialized solution but not implemented in a technical test context

- Add unit tests
- Add swagger documentation (with swaggo for exemple)
- improve mocks and interfaces drivers
- Handle rate limit errors.

### Limitations

Github API implement a [rate limit](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28):

- For unauthenticated requests, the rate limit allows for up to 60 requests per hour.
- For requests using Basic Authentication or OAuth, the rate limit allows for up to 5000 requests per hour.
- For requests using a GitHub App installation, the rate limit allows for up to 5000 requests per hour, (15000 for installations with the Organization Plan).
- For requests using a GitHub token, the rate limit allows for up to 1000 requests per hour per repository.

Globaly, one response of the API use 1 request to get 100 repository and 2 requests by repository to get the details. So, approximately 200 requests used for one demand.
The rate limit is reached quickly.

It could be interesting to use the user token to increase the system rate limit.
 