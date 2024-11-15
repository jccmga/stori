# 7. Testing the database

Date: 2024-11-14

## Status

Accepted

## Context

In the challenge, one of the bonus points asked us to persist information in a database.
To test the managed dependency, I'll need to run the database in a container.
I see two options:
1. Run a docker container through the Makefile and use a Suite for setup and teardown.
2. Use TestMain for control over the setup and teardown and run the container through code.

## Decision

I decided to use TestMain and to use the [Dockertest](https://github.com/ory/dockertest) library.
The library has a list of pre-built runnable examples for some popular databases, in particular the one that I'll use [Postgres](https://github.com/ory/dockertest/blob/v3/examples/PostgreSQL.md).

## Consequences

Pros:
- The integrations tests are easier to run as there is no need to run the container externally.

Cons:
- A new dependency is added to the stack.
