# 3. Entities from model

Date: 2024-11-11

## Status

Accepted

## Context

Since this is a coding challenge, there isn’t a business team or specific business logic behind it.
However, I want to stay consistent with the business language used in the challenge and reflect that in the model.

## Decision

The main business concepts I've identified are Transaction and AccountSummary.
Both of them will be kept in the `model` package.

## Consequences

Instructions suggest all processes are bounded to a single file, with a list of transactions. It could be
imaginable that several files could belong to the same account or transactions could be added forward.
