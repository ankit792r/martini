# Martini

A CLI for common post-development tasks: devbox project setup and Flutter Android/iOS/web workflows.

## Requirements

- Go 1.27+
- Flutter SDK (for `martini flutter` build commands)

## Install

Install for your user (recommended; uses `go install` → `~/go/bin` by default):

```bash
make install-user
# or
./build.sh
```

System-wide install (copies the built binary to `/usr/local/bin`):

```bash
sudo make install
```

Make sure `~/go/bin` (or your `GOBIN`) is on your `PATH` for user installs.

You can also run without installing:

```bash
go run .
```

## Commands

### `martini devbox`

Interactive wizard that creates a `devbox.json` in the current directory with project name, custom environment variables, and scripts.

```bash
martini devbox
```

### `martini android`

Commands for native Android Studio projects (not Flutter).

```bash
martini android --generate-upload-keystore
cd /path/to/android/project && martini android --generate-upload-keystore
```

Creates an upload JKS in the project root and in your home directory, writes `key.properties`, and saves credentials to `~/{keystore-name}-upload-key.txt`. It does **not** modify `build.gradle.kts` (use this when Gradle signing is already configured).

Flags: `--generate-upload-keystore`, `--path`, `--keystore-name`, `--store-password`, `--key-password`, `--key-alias`, `--common-name`. Omitted keystore fields are prompted interactively.

### `martini flutter`

Interactive task picker, or run tasks directly with flags.

```bash
martini flutter
cd /path/to/flutter/project && martini flutter --build-release
martini flutter --path /path/to/flutter/project --build-release
```

| Task | Flag |
| --- | --- |
| Update Gradle package name | `--update-package-name` |
| Generate upload keystore | `--generate-upload-keystore` |
| Update signing config | `--update-signing-config` |
| Build release APKs and AAB | `--build-release` |
| Update Android, iOS, and web icons | `--update-icons` |

Common flags:

- `--path`, `-p` — Flutter project directory (defaults to the current working directory)
- `--icons-path` — IconKitchen output directory (required with `--update-icons`)

Examples:

```bash
martini flutter --path ~/Projects/myapp --update-package-name
martini flutter --path ~/Projects/myapp --generate-upload-keystore --update-signing-config
martini flutter --path ~/Projects/myapp --build-release
martini flutter --path ~/Projects/myapp --update-icons --icons-path ~/Downloads/IconKitchen-Output
```

Release builds produce a zip in your home directory with APKs, AAB, and SHA1 hashes.

## Development

```bash
make check    # fmt, vet, test
make build    # bin/martini
make run      # go run .
```

Cross-compile for another platform:

```bash
GOOS=linux GOARCH=amd64 ./build.sh
```

## License

MIT
