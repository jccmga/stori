# 2. Folder Structure

Date: 2024-11-11

## Status

Accepted

## Context

Go is the chosen language for this project because it's simple and familiar.
There’s a lot of discussion about the "best" way to organize a Go project.

## Decision

I’ve decided to use the `cmd` folder to store the main package.
This project won’t include internal or pkg folders.
Instead, I’ll name the other folders clearly, so it’s easy to tell what they do.

## Consequences

By not using the `internal` folder, the code could be reused in other projects.
However, since this is a coding challenge, this shouldn't be an issue.
Additionally, this will make responsibilities from each package more intuitive.
