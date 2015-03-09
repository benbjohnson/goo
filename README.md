goo 
===

Goo is a simple wrapper around the Go toolchain. This code is in an alpha
stage and will change.


## Usage

You can install `goo` by using `go get`:

```sh
$ go install github.com/benbjohnson/goo
```

Any command you execute with `goo` will have its arguments redirected to `go`:

```sh
$ goo build
$ goo test -v -run=TestMyFunc
```
