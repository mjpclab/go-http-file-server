# Go HTTP File Server
Share local file system by HTTP.

![Go HTTP File Server pages](doc/ghfs.gif)

## Compile
Minimal required Go version is 1.10.
```bash
cd src
go build main.go
```

If default html template file under `src/tpl` changed, run
```bash
cd src
make tpls
```
to re-embed templates into go files. Then compile the project again.

## Usage
```
-l <ip|[:]port|ip:port>
--listen <ip|[:]port|ip:port>
    Optional IP and port the server listens on, e.g. ":80" or "127.0.0.1:80".
    If port is not specified, use "80" for pure HTTP mode, or "443" for TLS mode.
    flag "-l" or "--listen" can be ommitted.

-r <directory>
--root <directory>
    Root directory of the server.
    Defaults to current working directory.

-a <separator><url-path><separator><file-system-path>
--alias <separator><url-path><separator><file-system-path>
    Set path alias. e.g. ":/doc:/usr/share/doc"

-u <url-path>
--upload <url-path>
    Set url path that allows to upload files.
    If filename exists, will try to add or increase numeric prefix.
    Use it with care.

-c <file>
--cert <file>
    Specify TLS certificate file.
    If both "cert" and "key" are specified, the server serves in TLS mode for HTTPS protocol.

-k <file>
--key <file>
    Specify key file of TLS certificate.
    If both "cert" and "key" are specified, the server serves in TLS mode for HTTPS protocol.

-t <file>
--template <file>
    Use a custom template file for rendering pages, instead of builtin template.

-L <file>
--access-log <file>
    Access log file.
    Set "-" to use stdout.
    Set to empty to disable access log.

-E <file>
--error-log <file>
    Error log file.
    Set "-" to use stderr.
    Set to empty to disable error log.
```
