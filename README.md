# Martini

A CLI for common post-development tasks: devbox project setup and Flutter Android/iOS/web workflows.

## Requirements

- Go 1.27+
- Flutter SDK (for `martini flutter` build commands)

## Install

Build a static binary and install it to your Go bin directory (`~/go/bin` by default):

```bash
./build.sh
```

Make sure `~/go/bin` is on your `PATH`.

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

### `martini flutter`

Interactive task picker, or run tasks directly with flags.

```bash
martini flutter
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

- `--path`, `-p` — Flutter project directory
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
go test ./...
```

Cross-compile for another platform:

```bash
GOOS=linux GOARCH=amd64 ./build.sh
```

## License

MIT
