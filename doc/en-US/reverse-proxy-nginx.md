# Use nginx as a reverse proxy to ghfs

Ensure nginx module `ngx_http_proxy_module` was installed.

For example, proxying `/files` from nginx to ghfs root path:

## Method 1: Proxy stripped path to ghfs(recommended)

### Run ghfs and listen

```sh
ghfs -l 8080 -r /tmp/
```

### Config nginx reverse proxy

Note to **preserve** tailing `/` in `proxy_pass`:

```conf
location /files {
 proxy_pass http://localhost:8080/;
}
```

### Result

When a request which its path is `/files/dirs` to nginx, the ghfs actually got request to `/dirs`.

## Method 2: Proxy original path to ghfs

### Run ghfs and listen, with prefix `/files`

```sh
ghfs -l 8080 -r /tmp/ --prefix /files
```

### Config nginx reverse proxy

Note to **omit** tailing `/` in `proxy_pass`:

```conf
location /files {
 proxy_pass http://localhost:8080;
}
```

### Result

When a request which its path is `/files/dirs` to nginx, the ghfs got the same request path as from nginx.

ghfs will strip prefix `/files` internally.
