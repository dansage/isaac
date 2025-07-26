# ISAAC
[![Go Reference](https://pkg.go.dev/badge/go.dsage.org/isaac.svg)][3]  
ISAAC is an implementation of Bob Jenkins' [ISAAC CSPRNG][1] in pure Go, originally ported by [George Tankersley][2].

## Getting Started
Start by adding the package as a dependency to your project:

```shell
go get -u go.dsage.org/isaac
```

For an overview of how to use the library, see the [documentation][3].

## Requirements
There are no direct dependencies for this project, but testing only takes place on the following hosts:

- Linux
  - Arch Linux (manual)
  - Ubuntu (automated, via GitHub Actions: `ubuntu-latest`)
- macOS (automated, via GitHub Actions: `macos-latest`)
- Windows (automated, via GitHub Actions: `windows-latest`)

## Security Vulnerabilities
The automated tests verify the library outputs the same results as the official reference implementation, but no claims
or guarantees around the safety of the algorithm or library are made. This fork will be archived shortly after release.

There are many modern CSPRNG implementations in existence. Please use this project at your own risk.

## License
This library is released under the [MIT License](https://choosealicense.com/licenses/mit/) (see [`LICENSE`](LICENSE)).

[1]: http://www.burtleburtle.net/bob/rand/isaacafa.html
[2]: https://github.com/gtank
[3]: https://pkg.go.dev/go.dsage.org/isaac
