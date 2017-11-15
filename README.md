# Go HTTP Server
Super-simple static file servin' HTTP 1.1 server built on top of raw TCP sockets.

## Features
- Supports GET requests only
- Support for TLS
- Support for Gzipped responses
- Support for 304 cache requests
- Correct HTTP Status codes (mostly)
- Correct mimetype setting (mostly)
- One request maps to one goroutine

## TODO
- Different versions utilizing thread pools, async/event-driven patterns...
