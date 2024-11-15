# 4. Handling decimals

Date: 2024-11-11

## Status

Accepted

## Context

Transactions can include cents or bigger and we could argue the program should work with any currency.
This raises the question of how to handle decimals. As Golang doesn't have a standard type for handling decimals,
there are two approaches that I've considered:
1. Use a custom-made type with uint64 as the underlying type and keep track of the precision.
2. Use a third-party library to handle decimals.

Both come with their own pros and cons. A custom-made type would give me more control over the calculations, but it
will also increase the cost of maintenance. On the other hand, a third-party library will facilitate the implementation
but in a production environment, some governance policies regarding third-party libraries need to be followed.

## Decision

I've decided to use a third-party library to handle decimals. Several libraries are available:
1. [decimal](https://github.com/shopspring/decimal)
2. [apd](https://github.com/cockroachdb/apd)
3. [alpaca](https://github.com/alpacahq/alpacadecimal)
4. [govalues](https://github.com/govalues/decimal)

I've decided to use `govalues` as some other libraries are no heavily maintained, and benchmarks from this library are promising.                                                       

## Consequences

This inevitable creates a dependency on a third-party library that will be everywhere.
