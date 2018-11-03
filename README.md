go-libre-storage
=========

Cloud storage and file share platform writen in [Go](https://golang.org) (Golang).

## Prerequisites

* [go](https://golang.org) compiler
* [TDM-GCC-64](http://tdm-gcc.tdragon.net) or [mingw64](https://sourceforge.net/projects/mingw-w64/) - for windows only

For RaspberryPi 3 (BCM2837) builds:
* [armv7 hf gcc toolchain](http://gnutoolchains.com/raspberry/)

## Build

Windows
```
> mingw32-make && mingw32-make build
```

Unix
```
> make && make build
```

Cross-compile for RaspberryPi 3 on Windows
```
> mingw32-make && mingw32-make build-arm7hf
```

## Run development instance

```
> HOST="localhost" ./bin/go-libre-storage
```

## License

[MIT](http://opensource.org/licenses/MIT)

Copyright (c) 2018-present, towa48