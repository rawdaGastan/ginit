# Ginit

Ginit is a Golang version of the popular [Foreman](https://ddollar.github.io/foreman/) tool.
It only supports GNU/Linux

> Foreman is a manager for Procfile-based applications. Its aim is to abstract away the details of the Procfile format, and allow you to either run your application directly or export it to some other process management format.

## How to use

- Get the binary

> Download the latest from the [releases page](https://github.com/rawdaGastan/ginit/releases)

- Extract the downloaded ginit

After downloading the binary, run the following command inside the extracted folder:

```bash
sudo cp ginit /usr/local/bin
```

- Run the following command wherever your project exists given `Procfile` and `.env` file

```bash
ginit start -f Procfile -e .env`
```

## Run the demo using the installed ginit

- Run cmd `ginit start -f demo/Procfile -e demo/.env`

## Run the demo locally

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
