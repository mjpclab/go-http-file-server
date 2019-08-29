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
-l|--listen <ip|[:]port|ip:port>
    Optional IP and port the server listens on, e.g. ":80" or "127.0.0.1:80".
    If port is not specified, use "80" for pure HTTP mode, or "443" for TLS mode.
    flag "-l" or "--listen" can be ommitted.

-r|--root <directory>
    Root directory of the server.
    Defaults to current working directory.

-a|--alias <separator><url-path><separator><file-system-path>, ...
    Set path alias. e.g. ":/doc:/usr/share/doc"

-u|--upload <url-path>, ...
    Set url path that allows to upload files.
    If filename exists, will try to add or increase numeric prefix.
    Use it with care.

-c|--cert <file>
    Specify TLS certificate file.
    If both "cert" and "key" are specified, the server serves in TLS mode for HTTPS protocol.

-k|--key <file>
    Specify key file of TLS certificate.
    If both "cert" and "key" are specified, the server serves in TLS mode for HTTPS protocol.

-t|--template <file>
    Use a custom template file for rendering pages, instead of builtin template.

-S|--show <wildcard>, ...
-SD|--show-dir <wildcard>, ...
-SF|--show-file <wildcard>, ...
    If specified, files or directories match wildcards(except hidden by hide option) will be shown. 

-H|--hide <wildcard>, ...
-HD|--hide-dir <wildcard>, ...
-HF|--hide-file <wildcard>, ...
    If specified, files or directories match wildcards will not be shown.

-L|--access-log <file>
    Access log file.
    Set "-" to use stdout.
    Set to empty to disable access log.

-E|--error-log <file>
    Error log file.
    Set "-" to use stderr.
    Set to empty to disable error log.
```
