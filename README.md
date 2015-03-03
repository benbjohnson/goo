gofirst
=======

This program is a simple utility to highlight the first reported file and
line number piped in. This is useful when developing Go as it will let you
jump to that file/line more quickly.


## Usage

Currently, `go` toolchain output must be redirected to `STDOUT`:

```sh
$ go install github.com/benbjohnson/gofirst
$ go test 2>&1 | gofirst
```

Now your shell output should show the first `FILE:LINE` in bold and if you
are on Mac OS X then it will copy that to the clipboard.
