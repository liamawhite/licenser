# Licenser

Verify and append licenses to your GitHub repositories.

## Supported Licenses & Languages

Licenses:

- Apache 2.0

Languages:

- Bazel
- C/C++
- Dockerfile
- Golang
- Javascript
- Make
- Protobuf
- Python
- Shell
- TypeScript
- YAML

Licenser will also automatically ignore the following files:

- `*.md`, `*.golden`
- `.gitignore`
- Files that should be ignored according to `.gitignore` (experimental)
- `.licenserignore`
- Files that should be ignored according to `.licenserignore` (experimental)

## Install

To install on macOS, use Homebrew.

```sh
brew install liamawhite/licenser/licenser
```

To install on Ubuntu/Debian, use `wget`.

```sh
wget -c https://github.com/liamawhite/licenser/releases/download/v${VERSION}/licenser_${VERSION}_Linux_x86_64.tar.gz -O - | sudo tar -xz -C /usr/bin
```

To install on other platforms, download from the [releases section](https://github.com/liamawhite/licenser/releases).

## Verifying Licenses in your Files

To verify that licenses are present in all files in a repository, run the `verify` command at the root, with the `--recurse` flag.

```sh
licenser verify -r
```

## Apply Licenses to your Files

To prepend licenses to all files in a repository, run the `apply` command at the root, with the `--recurse` flag, passing in the copyright owner.

```sh
licenser apply -r "Copyright Owner"
```
