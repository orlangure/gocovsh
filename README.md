# gocovsh

`gocovsh` is a CLI viewer of Go test coverage reports.

## Installation

```
$ go install github.com/orlangure/gocovsh@latest
```

More installation options will follow.

## Usage

1. Generate Go coverage report at your project's root with
    ```bash
    $ go test -cover -coverprofile coverage.out
   ```

   For more information about generating Go coverage reports, see [my blog
   post](https://fedorov.dev/posts/2020-06-27-golang-end-to-end-test-coverage/).

2. Run `gocovsh` at the same folder with `coverage.out` report and `go.mod`
   file (`go.mod` is required).

   ```bash
   $ gocovsh
   $ gocovsh --profile profile.out # for other coverage profile names
   ```

3. Use `j/k/enter/esc` keys to explore the report. See built-in help for more
   key-bindings.

<img width="900" alt="image" src="https://user-images.githubusercontent.com/10244414/151678881-74b52fe5-0dea-4411-aa65-2343d71b8516.png">
<img width="900" alt="image" src="https://user-images.githubusercontent.com/10244414/151678915-e323a185-679f-48ff-9582-63c48edd09c0.png">

