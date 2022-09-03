# 启动docker容器来共享本地文件系统

## 镜像名称
`mjpclab/ghfs`

## 容器默认配置

- 根目录为 `/var/ghfs/`
- 在 `8080` 私有端口提供纯HTTP服务
- 在 `8443` 私有端口提供HTTPS服务
- TLS证书位于 `/etc/server.crt`
- TLS证书私钥位于 `/etc/server.key`

## 提供纯HTTP服务

```sh
docker run \
-v /PATH/TO/ROOT:/var/ghfs \
-p HOST_PORT:8080 \
mjpclab/ghfs
```

## 提供HTTPS服务

```sh
docker run \
-v /PATH/TO/CERT:/etc/server.crt \
-v /PATH/TO/KEY:/etc/server.key \
-v /PATH/TO/ROOT:/var/ghfs \
-p HOST_PORT:8443 \
mjpclab/ghfs
```
