goo ![Version](http://img.shields.io/badge/version-beta-blue.png)
===

Goo is a simple wrapper around the Go toolchain. It copies and bolds the first 
file + line number error that occurs in the standard output and standard input.
This allows you switch back to your text editor and quickly jump to the line.

This code is in in beta and has been tested on Mac OS X only.


## Usage

You can install `goo` by using `go get`:

```sh
$ go get github.com/benbjohnson/goo
```

Any command you execute with `goo` will have its arguments redirected to `go`:

```sh
$ goo build
$ goo test -v -run=TestMyFunc
```

If an error occurs with the format `FILE:LINE` then it will be copied to your
clipboard and the entire line will be bolded. This allows you to jump to your
text editor and quickly jump to the file.
