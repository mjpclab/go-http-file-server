# Use apache as a reverse proxy to ghfs

Ensure to enable module `proxy` and `proxy_http`.

```sh
a2enmod proxy proxy_http
```

For example, proxying `/files` from apache to ghfs root path:

### Run ghfs and listen

```sh
ghfs -l 8080 -r /tmp/
```

### Config apache reverse proxy

```conf
ProxyPass /files/ http://localhost:8080/
ProxyPassReverse /files/ http://localhost:8080/
```
