# fastreplace

[<img src="https://img.shields.io/github/license/esrrhs/fastreplace">](https://github.com/esrrhs/fastreplace)
[<img src="https://img.shields.io/github/languages/top/esrrhs/fastreplace">](https://github.com/esrrhs/fastreplace)
[<img src="https://img.shields.io/github/v/release/esrrhs/fastreplace">](https://github.com/esrrhs/fastreplace/releases)
[<img src="https://img.shields.io/github/downloads/esrrhs/fastreplace/total">](https://github.com/esrrhs/fastreplace/releases)
[<img src="https://img.shields.io/docker/pulls/esrrhs/fastreplace">](https://hub.docker.com/repository/docker/esrrhs/fastreplace)
[<img src="https://img.shields.io/github/actions/workflow/status/esrrhs/fastreplace/go.yml?branch=master">](https://github.com/esrrhs/fastreplace/actions)

> **A fast multithreaded CLI tool that replaces text across files under a directory.**

[English](README.md) | [Chinese](README_CN.md)

---

## Overview

**fastreplace** walks a path (recursively), optionally filters by file suffix, and replaces all occurrences of a string with another — using a bounded worker pool for speed.

## Install

Download a prebuilt binary from the **[Releases](https://github.com/esrrhs/fastreplace/releases)** page, or build from source:

```bash
git clone https://github.com/esrrhs/fastreplace.git
cd fastreplace
go build -o fastreplace .
```

Docker:

```bash
docker pull esrrhs/fastreplace
```

## Usage

Replace `aaa` with `bbb` in every `.txt` file under the current directory:

```bash
./fastreplace -path ./ -file .txt -from aaa -to bbb
```

Show version:

```bash
./fastreplace -version
```

More flags:

```text
Usage of ./fastreplace:
  -file string
        file format (suffix filter, e.g. .txt)
  -from string
        old string
  -path string
        replace path (default "./")
  -thread int
        replace thread (default 32)
  -to string
        new string
  -v    show info
  -version
        show version
```

## Build release packages

```bash
./pack.sh
```

Cross-platform archives are written to `pack/`. The **Release** workflow watches `version.go`: when its version string changes on `master` (and the tag does not exist yet), CI runs `./pack.sh` and publishes a GitHub Release.

## License

MIT — see [LICENSE](LICENSE).
