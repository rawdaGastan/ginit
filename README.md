# Ginit

Ginit is a Golang version of the popular [Foreman](https://ddollar.github.io/foreman/) tool.
It only supports GNU/Linux 

> Foreman is a manager for Procfile-based applications. Its aim is to abstract away the details of the Procfile format, and allow you to either run your application directly or export it to some other process management format.

## How to use

- There is a demo folder in the project you can use:

- build `task build` or `make build`

- Run cmd `./bin/ginit start -f demo/Procfile -e demo/.env`

## Testing

Use these command to run the tests

```bash
task test
```

```bash
make test
```

## Coverage

- create a `coverage` folder
- Use these command to see the coverage

```bash
task coverage
```

```bash
make coverage
```

- Open coverage/coverage.html to trace the coverage

## Benchmarks

Use these command to run the tests

```bash
task benchmarks
```

```bash
make benchmarks
```
