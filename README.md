[![CircleCI](https://circleci.com/gh/marccarre/go-github-release/tree/master.svg?style=shield)](https://circleci.com/gh/marccarre/go-github-release/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/marccarre/go-github-release)](https://goreportcard.com/report/github.com/marccarre/go-github-release)
[![Coverage Status](https://coveralls.io/repos/github/marccarre/go-github-release/badge.svg)](https://coveralls.io/github/marccarre/go-github-release)
[![codecov](https://codecov.io/gh/marccarre/go-github-release/branch/master/graph/badge.svg)](https://codecov.io/gh/marccarre/go-github-release)
[![Docker Repository on Quay](https://quay.io/repository/marccarre/go-github-release/status)](https://quay.io/repository/marccarre/go-github-release)

# go-github-release

## Features

- `release`:
  - creates the GitHub release matching the provided tag,
  - signs the provided release assets,
  - uploads the provided release assets and their signatures to the GitHub release.

## Installation

1. Download the [latest version](https://github.com/marccarre/go-github-release/releases/latest) of `ghrelease` and its signature.
2. [Optional] Verify `ghrelease`'s signature:

    1. Import my CI/CD GPG key:

      ```console
      $ gpg --recv-keys 92A05461
      gpg: key 6A74FEA692A05461: public key "Marc CARRE (GitHub software releases) <carre.marc@gmail.com>" imported
      gpg: Total number processed: 1
      gpg:               imported: 1
      ```

    2. Check its fingerprint:

      ```console
      $ gpg --fingerprint 92A05461
      gpg: checking the trustdb
      gpg: marginals needed: 3  completes needed: 1  trust model: pgp
      gpg: depth: 0  valid:   1  signed:   0  trust: 0-, 0q, 0n, 0m, 0f, 1u
      gpg: next trustdb check due at 2020-11-13
      pub   rsa4096 2018-11-11 [SC] [expires: 2019-11-11]
            49A9 5DE0 562B 072A 82B4  F072 6A74 FEA6 92A0 5461
      uid           [ unknown] Marc CARRE (GitHub software releases) <carre.marc@gmail.com>
      sub   rsa4096 2018-11-11 [E] [expires: 2019-11-11]
      ```

    3. Import my [personal GPG key](https://keybase.io/marccarre):

      ```console
      $ gpg --recv-keys F69B8B32
      gpg: key 062658EFF69B8B32: public key "Marc CARRE <carre.marc@gmail.com>" imported
      gpg: Total number processed: 1
      gpg:               imported: 1
      ```

    4. Ensure my CI/CD key is signed by [me](https://keybase.io/marccarre):

      ```console
      $ gpg --list-signatures 92A05461
      pub   rsa4096 2018-11-11 [SC] [expires: 2019-11-11]
            49A95DE0562B072A82B4F0726A74FEA692A05461
      uid           [ unknown] Marc CARRE (GitHub software releases) <carre.marc@gmail.com>
      sig 3        6A74FEA692A05461 2018-11-11  Marc CARRE (GitHub software releases) <carre.marc@gmail.com>
      sig 3        062658EFF69B8B32 2018-11-11  Marc CARRE <carre.marc@gmail.com>
      sub   rsa4096 2018-11-11 [E] [expires: 2019-11-11]
      sig          6A74FEA692A05461 2018-11-11  Marc CARRE (GitHub software releases) <carre.marc@gmail.com>
      ```

    5. Check `ghrelease` against its detached signature, e.g.:

      ```console
      $ gpg --verify ghrelease-v1.0.0-linux.asc ghrelease-v1.0.0-linux
      gpg: Signature made Wed 14 Nov 12:41:17 2018 JST
      gpg:                using RSA key 6A74FEA692A05461
      gpg: Good signature from "Marc CARRE (GitHub software releases) <carre.marc@gmail.com>" [ultimate]
      ```

## Usage

```console
$ gpg --export-secret-keys <key-id> > /path/to/your/private/key.asc
$ export GPG_PASSWD="..."
$ export GITHUB_API_TOKEN="..."
$ ghrelease release --help
Sign and upload the provided release assets on GitHub under the release corresponding to the provided tag

Usage:
  ghrelease release [flags]

Flags:
  -d, --draft          Should the release be a draft release, default: true (default true)
  -h, --help           help for release
  -k, --key string     Path to the private GPG key to use to sign the release assets
  -o, --owner string   GitHub owner, e.g. marccarre in github.com/marccarre/go-github-release
  -r, --repo string    GitHub repository, e.g. go-github-release in github.com/marccarre/go-github-release
  -t, --tag string     Git tag corresponding to the release to perform, e.g. v1.0.0

$ ghrelease release -o marccarre -r go-github-release -t v1.0.0 -k /path/to/your/private/key.asc <your-binary> ...
{"level":"info","msg":"creating release","owner":"marccarre","repo":"go-github-release","tag":"v1.0.0","draft":true}
{"level":"info","msg":"successfully created release","owner":"marccarre","repo":"go-github-release","tag":"v1.0.0","draft":true}
{"level":"info","msg":"signing release asset","file":"<your-binary>"}
{"level":"info","msg":"successfully signed release asset","file":"<your-binary>"}
{"level":"info","msg":"uploading release asset","file":"<your-binary>","release":"v1.0.0"}
{"level":"info","msg":"successfully uploaded release asset","file":"<your-binary>","release":"v1.0.0","asset":"<your-binary>"}
{"level":"info","msg":"uploading release asset","file":"<your-binary>.asc","release":"v1.0.0"}
{"level":"info","msg":"successfully uploaded release asset","file":"<your-binary>.asc","release":"v1.0.0","asset":"<your-binary>.asc"}
```

## Development

### Setup

- Install [`docker`](https://store.docker.com/search?type=edition&offering=community)
- Install `make`

That's all folks!
All other tools are packaged in build Docker images (see `Dockerfile`) to ensure any machine can build easily, hence avoiding the "[_it works on my machine_](http://www.codinghorror.com/blog/2007/03/the-works-on-my-machine-certification-program.html)" syndrome.

### Build

```console
make
```

### Lint

```console
make lint
```

### Test

```console
make test
```

### Release

```console
git tag vX.Y.Z <commit-hash> -a -m vX.Y.Z
git push --tags
```
