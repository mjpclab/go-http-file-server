# Go HTTP File Server
Simple command line based HTTP file server to share local file system.

![Go HTTP File Server pages](doc/ghfs.gif)

## Features
- More friendly UI than Apache/Nginx directory index page
- Adapt for mobile display
- Can download the whole contents of current directory as archive file if enabled
- Can upload files to current directory if enabled
- Can specify a custom template for page rendering
- Support location alias(mount another directory to url location)

## Compile
Minimal required Go version is 1.9.
```bash
cd src
go build main.go
```

If default html template file under `src/tpl` changed, need to re-embed templates into go files:
```bash
cd src
make tpls
```
Then compile the project again.

## Usage
```
server [options]

-l|--listen <ip|port|:port|ip:port|socket> ...
    IP and port the server listens on, e.g. ":80" or "127.0.0.1:80".
    If --cert and --key are specified, port listens for TLS connection.
    If port is not specified, use "80" for pure HTTP mode, or "443" for TLS mode.
    If value contains "/" then treat it as a unix socket file.
    Flag "-l" or "--listen" can be ommitted.
--listen-plain <ip|port|:port|ip:port|socket> ...
    Similar to --listen, but force to use non-TLS mode
--listen-tls <ip|port|:port|ip:port|socket> ...
    Similar to --listen, but force to use TLS mode, will failed if cert or key is not specified.

-r|--root <directory>
    Root directory of the server.
    Defaults to current working directory.

-a|--alias <separator><url-path><separator><file-system-path> ...
    Set path alias. e.g. ":/doc:/usr/share/doc"

-U|--global-upload
    Allow upload files for all url paths.
    Use it with care.

-u|--upload <url-path> ...
    Set url paths that allows to upload files.
    If filename exists, will try to add or increase numeric prefix.
    Use it with care.

-A|--global-archive
    Allow user to download the whole contents of current directory for all url paths.
    A download link will appear on top part of the page.
    Make sure there is no circular symbol links.

--archive <url-path> ...
    Allow user to download the whole contents of current directory for specific url paths.

-c|--cert <file>
    Specify TLS certificate file.

-k|--key <file>
    Specify key file of TLS certificate.

-t|--template <file>
    Use a custom template file for rendering pages, instead of builtin template.

-S|--show <wildcard> ...
-SD|--show-dir <wildcard> ...
-SF|--show-file <wildcard> ...
    If specified, files or directories match wildcards(except hidden by hide option) will be shown. 

-H|--hide <wildcard> ...
-HD|--hide-dir <wildcard> ...
-HF|--hide-file <wildcard> ...
    If specified, files or directories match wildcards will not be shown.

-L|--access-log <file>
    Access log file.
    Set "-" to use stdout.
    Set to empty to disable access log.

-E|--error-log <file>
    Error log file.
    Set "-" to use stderr.
    Set to empty to disable error log.
    Defaults to "-".

--config <file>
    External config file to load.

    Its content is option list of any other options,
    same as the form specified on command line,
    separated by whitespace characters.

    The external config's priority is lower than arguments specified on command line.
    If one option is specified on command line, then external config is ignored.
```

## Examples
Start server on port 8080, root directory is current working  directory:
```sh
server -l 8080
``` 

Start server on port 8080, root directory is /usr/share/doc:
```sh
server -l 8080 -r /usr/share/doc
```

Start server on default port, root directory is /tmp, and allow upload files to file system directory /tmp/data:
```sh
server -r /tmp -u /data
```

Share files from /etc, but also mount /usr/share/doc to url path /doc
```sh
server -r /etc -a :/doc:/usr/share/doc
```

Start server on port 8080, serve for HTTPS protocol
```sh
server -k /path/to/certificate/key -c /path/to/certificate/file -l 8080
```

Do not show hidden unix directories and files that starts with `.`.
Tips: wrap wildcard by quotes to prevent expanding by shell.
```sh
server -H '.*'
```

Show access log on console:
```sh
server -L -
```
