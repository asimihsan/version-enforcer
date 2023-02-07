# version-enforcer

This is a tool to enforce the versions of tools used in a runtime environment.

## Installation

```sh
go install github.com/asimihsan/version-enforcer@0.0.2
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
