# 使用nginx反向代理到ghfs

确保nginx模块`ngx_http_proxy_module`已安装。

例如，nginx代理`/files`到ghfs根目录：

## 方法1：代理剥离的路径到ghfs（推荐）

### 运行ghfs并侦听

```sh
ghfs -l 8080 -r /tmp/
```

### 配置nginx反向代理

注意**保留**`proxy_pass`尾部的`/`。

```conf
location /files {
 proxy_pass http://localhost:8080/;
}
```

### 结果

当请求路径`/files/dirs`到达nginx，ghfs实际得到的路径为`/dirs`。

## 方法2：代理原始路径到ghfs

### 运行ghfs并侦听，前缀为`/files`

```sh
ghfs -l 8080 -r /tmp/ --prefix /files
```

### 配置nginx反向代理

注意**省略**`proxy_pass`尾部的`/`。

```conf
location /files {
 proxy_pass http://localhost:8080;
}
```

### 结果

当请求路径`/files/dirs`到达nginx，ghfs也获得同样的路径。

ghfs会在内部剥离`/files`前缀。
