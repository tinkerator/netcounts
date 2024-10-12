# netcounts - a package for monitoring network activity

## Overview

The `netcounts` package provides a Go API for monitoring network
traffic. It works by parsing the command output of the `ifconfig`
program.

```
$ git clone https://github.com/tinkerator/netcounts.git
$ cd netcounts
$ go get
$ go build examples/watch.go
$ ./watch
```

The default output summarizes total packets and bytes transferred over
all of the non-loopback network devices.

## License info

The `netcounts` package is distributed with the same BSD 3-clause
license as that used by [golang](https://golang.org/LICENSE) itself.

## Reporting bugs and feature requests

The `netcounts` package has been developed purely out of self-interest
If you find a bug or want to suggest a feature addition, please use
the [bug tracker](https://github.com/tinkerator/netcounts/issues).
