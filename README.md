# Blockbuffer Automated Video Converter

## Overview

Inspired by HandBrake, this is a web-enabled video transcoder that uses the [FFmpeg](https://ffmpeg.org/) library to convert videos to a variety of formats with options for manual and automatic processing, saved transcoding profiles, multi-output conversion, and much more. The application is built with a [Go](https://golang.org/) server and a [Nuxt 3](https://nuxtjs.org/) web client.

## Planned Features

The core functionality of the application is to convert videos to a variety of formats using FFmpeg. The current version handles individual files being automatically converted to a default format and placed in the `output` directory. The following features are planned for the application:

- [x] Automated video transcoding (single file configuration)
- [ ] Manual video transcoding (directory configuration)
- [ ] Transcoding profiles
  - [ ] Video codec
  - [ ] Audio codec
- [ ] Multi-output rules
  - [ ] Transcoding profiles
  - [ ] Output directories
- [ ] Web interface
  - [ ] Global settings
  - [x] Video upload
  - [ ] Video download
  - [ ] Video conversion status
  - [ ] Configuration options
  - [ ] Transcoding profiles
  - [ ] Multi-output rules
- [ ] Docker containerization

This is an evolving project and the list of features will be updated as development progresses.

## Installation

Clone the repository, cd into project directory, and install dependencies:
```bash
git clone https://github.io/nendocreative/blockbuffer.git --depth 1

cd blockbuffer

go mod download

yarn --cwd ./blockbuffer-fe install
```

## Development Server

Start the development server on `http://localhost:8080`:

```bash
go run .
```

This will start the Go server and run Nuxt3 in a development mode, allowing for hot-reloading and other development features. Be sure to access it from the port Go is running (default: 8080) instead of the Nuxt port (default: 3000). Otherwise, the API calls will not work.

## Production Build

Build the frontend for production as a static site:
```bash
yarn --cwd ./blockbuffer-fe generate 
```

Build the Go server for production:
```bash
go build .
```

you should now have a `blockbuffer` binary in the project directory. Run it to start the server.
