# 6. Adapters and business logic borderline

Date: 2024-11-13

## Status

Accepted

## Context

The input/output (I/O) for the simulator will come from the command line.

## Decision

I chose to follow the Stable Dependency Principle (SDP) and use a
Ports-and-Adapters architecture (Hexagonal Architecture) to improve extensibility and maintainability.

## Consequences

The adapters are placed in their own package to keep concerns separated, and they implement interfaces the application depends on.
Dependencies flow towards the application and model, maintaining the SDP.
