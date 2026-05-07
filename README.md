# Password Generator

A small Go web application for generating secure passwords with adjustable length, character set toggles, batch generation, and entropy calculation.

## Features

- Set password length up to 128 characters
- Generate up to 50 passwords at once
- Toggle lowercase, uppercase, numbers, and symbols
- Calculate entropy in bits from the selected character pool
- Copy individual passwords or the full generated batch
- Dependency-free Go backend with embedded static assets

## Run

Install Go 1.22 or newer, then run:

```sh
go run .
```

Open [http://localhost:8080](http://localhost:8080).

## Test

```sh
go test ./...
```
