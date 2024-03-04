# 🐸 grog

A lightweight node package manager written in go.

## Installation
- Install [Go](https://go.dev/doc/install)
- Clone the repository
- Run `go build && go install`

## Usage

`grog install [package]`

`grog install [package]@[version]`

In order to properly use grog, you must use the `--preserve-symlinks` flag when running `node yourfile.js`. 

## How fast is grog?

**CLEAN INSTALLATION**

Benchmark of [bun](https://bun.sh/) vs. grog using [hyperfine](https://github.com/sharkdp/hyperfine)

<img width="1419" alt="Grog vs Bun Benchmark React" src="https://github.com/LOTaher/grog/assets/86690869/9f43fc1c-07c7-49dd-8f0d-b12abdf5f2b0">

*Due to the nature of HTTP, it is hard to give an accurate answer on who is "faster", as there are plenty of times bun is faster than grog. Bun is currently faster at loading cached modules.*

Benchmark of [npm](https://www.npmjs.com/) vs. grog using [hyperfine](https://github.com/sharkdp/hyperfine)

<img width="1419" alt="Grog vs NPM Benchmark React" src="https://github.com/LOTaher/grog/assets/86690869/f61547f0-12c4-404b-b46c-ca076a2d2c36">

## Features

- `grog install`: Installs a package, and caches the specific version in the `$HOME/.grog/cache` directory.
- `grog clear`: Clears the cache.

## Coming Soon

- Terminal user interface.
- The generation of package locks for each installed package to avoid the re-retrieval of dependencies.
- Creation and maintainence of a `package.json` in the working directory
- Creation and maintainence of a `package-lock.json` in the project directory 
- `grog uninstall`: Uninstalls a package.
