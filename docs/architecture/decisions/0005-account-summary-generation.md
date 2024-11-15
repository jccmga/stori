# 5. Account Summary Generation

Date: 2024-11-11

## Status

Accepted

## Context

In order to generate an account summary, we need to process all given transactions from the input file. Some limits
on the input file are still unknown. Depending on the limits, we may need to recur to batch processing in parallel.

## Decision

I've decided to process the transactions through a structure that will allow the work to be split in batches if the
limits of the file are significantly high.

## Consequences

The architecture will allow for extension in the future to allow "humongous" files to be processed later. Some thoughts
will be required to merge the work if that time comes.
