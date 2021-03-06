go-libre-storage
=========

Cloud storage and file share platform writen in [Go](https://golang.org) (Golang).

## Prerequisites

* [go](https://golang.org) compiler
* [gopls](https://github.com/golang/tools/blob/master/gopls/README.md) Go language server
* [TDM-GCC-64](http://tdm-gcc.tdragon.net) or [mingw64](https://sourceforge.net/projects/mingw-w64/) - for windows only
* [NodeJS](https://nodejs.org)

For RaspberryPi 3 (BCM2837) builds:
* [armv7 hf gcc toolchain](http://gnutoolchains.com/raspberry/)

## Getting source

```
> mkdir -p ~/go/src/github.com/towa48 && cd "$_"
> git clone https://github.com/towa48/go-libre-storage.git
```

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

## CLI

Restore files metadata from filesystem (disable files sharing)
```
> ./bin/go-libre-storage --crawl
```

Add user
```
> ./bin/go-libre-storage --add-user user2
```

Share folder to user
```
> ./bin/go-libre-storage --share-folder 8 --to user2 --write
```

## TODO

* Migrate to go modules https://blog.golang.org/migrating-to-go-modules
* Upgrate go dependencies
* Fix issues https://github.com/towa48/go-libre-storage/security/dependabot
* Migrate to React

## License

[GNU Affero General Public License v3.0](https://www.gnu.org/licenses/agpl.txt)

Copyright (c) 2018-present, towa48