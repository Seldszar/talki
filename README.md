# Talki

> Customizable, self-hosted Discord voice widget overlay

This application is aimed towards content creators using the Discord Streamkit voice widget and would like more control over design and functionality.
At the moment, it will follow the client's user and display the current voice channel state without the need to update the browser source, including group calls.

A video demonstration is available here: https://twitter.com/0xSeldszar/status/1775597534922027058

## Requirements

- [Go](https://go.dev)
- [Node.js](https://nodejs.org)

## Usage

```sh
$ go run main.go --help
```

In case you want to use the provided theme, go to the `public` folder and build the web project with the following commands:

```sh
$ npm install
$ npm run build
```

Then run the application with the following `public` flag:

```sh
$ go run main.go --public="public/dist"
```

## License

Copyright (c) 2024-present Alexandre Breteau

This software is released under the terms of the MIT License.
See the [LICENSE](LICENSE) file for further information.
