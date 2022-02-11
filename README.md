# gocovsh ![Test](https://github.com/orlangure/gnomock/workflows/Test/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/orlangure/gocovsh)](https://goreportcard.com/report/github.com/orlangure/gocovsh) [![codecov](https://codecov.io/gh/orlangure/gocovsh/branch/master/graph/badge.svg?token=F0XYPSEIMK)](https://codecov.io/gh/orlangure/gocovsh)

`gocovsh` is a tool for exploring [Go Coverage
reports](https://go.dev/blog/cover) from the command line.

[Screenshots](#screenshots) below. Don't skip the [Giving
back](#giving-back) section!

## Installation

### Using [Homebrew](https://brew.sh/)

```bash
# make sure brew is installed
brew tap orlangure/tap
brew update
brew install gocovsh
```

### Pre-built binary

Grab your pre-built binary from the
[Releases](https://github.com/orlangure/gocovsh/releases) page.

### From source

```bash
# install latest, probably unreleased version
go install github.com/orlangure/gocovsh@latest

# or install a specific version
go install github.com/orlangure/gocovsh@v0.1.0
```

## Usage

1. Generate Go coverage report at your project's root with
    ```bash
    go test -cover -coverprofile coverage.out
    ```

   For more information about generating Go coverage reports, see [my blog
   post](https://fedorov.dev/posts/2020-06-27-golang-end-to-end-test-coverage/).

2. Run `gocovsh` at the same folder with `coverage.out` report and `go.mod`
   file (`go.mod` is required).

   ```bash
   gocovsh
   gocovsh --profile profile.out # for other coverage profile names
   ```

3. Use `j/k/enter/esc` keys to explore the report. See built-in help for more
   key-bindings.

## Giving back

This is a free and open source project that hopefully helps its users, at least
a little. Even though I don't need donations to support it, I understand that
there are people that wish to give back anyway. If you are one of them, I
encourage you to [plant some trees with Tree
Nation](https://tree-nation.com/plant/offer) 🌲 🌳 🌴

If you want me to know about your contribution, make sure to use
`orlangure+gocovsh@gmail.com` as the recipient email.

Thank you!

## Screenshots

<img width="900" alt="image" src="https://user-images.githubusercontent.com/10244414/151678881-74b52fe5-0dea-4411-aa65-2343d71b8516.png">
<img width="900" alt="image" src="https://user-images.githubusercontent.com/10244414/151678915-e323a185-679f-48ff-9582-63c48edd09c0.png">

