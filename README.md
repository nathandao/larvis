# Larvis

Your trusted poker solver.

Larvis takes in 2 poker hands, each has 5 cards from this deck [2 3 4 5 6 7 8 9 T J K Q A], and it tell you the winning
hand, or if it's a tie.

Example:
```
larvis AATJ8 33322
# Hand 2
```

## Prerequisite

These dependencies would needed to be installed on your machine in order to develop and build the project:

- [Go 1.20](https://go.dev/dl/)
- [Docker](https://www.docker.com/)
- [Make](https://www.gnu.org/software/make/), or any other versions of make. We're using Makefile to provide shortcuts
  to build and test commands.
- (Optional) [golangci-lint](https://golangci-lint.run/): Only needed if you'd like to run the linter, but not required
  to build the project.

## Local build

To build larvis locally:

```shell
# Build larvis locally on your machine. The binary will be available at ./build/larvis
make build

# or without using make:
# go build -o ./build/larvis *.go
```

And to run the local build:

```shell
./build/larvis FIRST_HAND SECOND_HAND
```

## Docker build

To build the docker image:

```shell
# Build larvis locally on your machine. The binary will be available at ./build/larvis
make build-docker

# or without using make:
# docker build . -t larvis:local
```

And to run larvis in docker:

```shell
docker run -it larvis:local larvis FIRST_HAND SECOND_HAND

# or you can also ssh into the docker image,
# and `larvis` command should be available in the terminal session
docker run -it larvis:local /bin/sh
```

## Development

Some extra useful commands during development:

```shell
make test

make lint # requires golangci-lint to be installed
```

## Libaries used

- [testify](https://github.com/stretchr/testify): make use of the `assert` method that helps defining test assertions
  quicker and make tests more readable.
- [slices](https://pkg.go.dev/golang.org/x/exp@v0.0.0-20230203172020-98cc5a0785f9/slices): provides some nicer methods
  to sort slices. Although this is a part of the experimental golang repo, I thought this is a small enough project that
  it should be ok to try out the module here.
