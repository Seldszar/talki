# Talki

> Customizable, self-hosted Discord voice widget overlay

This application is aimed towards content creators using the Discord Streamkit voice widget and would like more control over design and functionality.
At the moment, it will follow the client's user and display the current voice channel state without the need to update the browser source, including group calls.

A video demonstration is available here: https://twitter.com/0xSeldszar/status/1775597534922027058

## Requirements

- [Go](https://go.dev)
- [Node.js](https://nodejs.org)

## Build

```sh
$ make all
```

## Usage

After buildng the program, run the binary corresponding on your platform, for example on Windows:

```sh
$ talki-windows-amd64.exe
```

An icon should be added to your system tray, from which you can open its menu and widget page and more.

If you want to use a custom public path, add a `--public` flag with the path to your folder:

```sh
$ talki-windows-amd64.exe --public="my-widget/dist"
```

You can find all available flags by adding the `--help` flag:

```sh
$ talki-windows-amd64.exe --help
```

## License

Copyright (c) 2024-present Alexandre Breteau

This software is released under the terms of the MIT License.
See the [LICENSE](LICENSE) file for further information.
