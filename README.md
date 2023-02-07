<h1 align="center">
  version-enforcer
</h1>

<h4 align="center">Enforce software version compliance.</h4>

<p align="center">
  <a href="#installation-example">Installation, example</a> •
  <a href="#usage">Usage</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#todo">TODO</a> •
  <a href="#license">License</a>
</p>

`version-enforcer` is a tool that enforces software version compliance. It is a
simple tool that can be used in CI pipelines and developer hosts to ensure that
the correct versions of software are installed.

## Installation, example

This will install `version-enforcer` to your `GOPATH`, which by default is `~/go`.

```sh
go install -v github.com/asimihsan/version-enforcer@0.0.8
```

Then create a config file:

```sh
tee version-enforcer.hcl <<EOF > /dev/null
binary "git" {
  version = "~2"
}
binary "make" {
  version = "^4.2.1"
}
EOF
```

Finally, run `version-enforcer`:

```sh
version-enforcer --config version-enforcer.hcl
```

## Usage

```sh
$ version-enforcer --help

Enforce tool versions

Usage:
  enforce --config <config file> [flags]

Flags:
      --config string   config file (e.g. version-enforcer.hcl)
  -h, --help            help for enforce
  -v, --verbose         verbose output
```

For example, you could run:

```
$ version-enforcer --config version-enforcer.hcl
```

## Configuration

Here is an example configuration file that specifies that

- `make` must be exactly `4.2.1`, and
- `git` must be between `>= 2.0.0` and `< 3.0.0`.

```hcl
binary "git" {
  version = "~2"
}

binary "make" {
  version = "^4.2.1"
}
```

The requirement specifications follow
[https://doc.rust-lang.org/cargo/reference/specifying-dependencies.html](https://doc.rust-lang.org/cargo/reference/specifying-dependencies.html).

## TODO

- [ ] Add support for `library` requirements.
- [ ] Output binary path in error messages.

## License

This project is licensed under the Apache License, Version 2.0. See
[LICENSE](LICENSE) for details.
