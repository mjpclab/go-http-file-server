# 使用apache反向代理到ghfs

确保启用模块`proxy`和`proxy_http`。

```sh
a2enmod proxy proxy_http
```

例如，apache代理`/files`到ghfs根目录：

### 运行ghfs并侦听

```sh
ghfs -l 8080 -r /tmp/
```

### 配置apache反向代理
```
ProxyPass /files/ http://localhost:8080/
ProxyPassReverse /files/ http://localhost:8080/
```
