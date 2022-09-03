# Start a docker container to share local file system

## Image name
`mjpclab/ghfs`

## Container's default configuration

- root directory is `/var/ghfs/`
- serve plain HTTP on private port `8080`
- serve HTTPS on private port `8443`
- TLS cert is located at `/etc/server.crt`
- TLS key is located at `/etc/server.key`

## Serve by plain HTTP

```sh
docker run \
-v /PATH/TO/ROOT:/var/ghfs \
-p HOST_PORT:8080 \
mjpclab/ghfs
```

## Serve by HTTPS

```sh
docker run \
-v /PATH/TO/CERT:/etc/server.crt \
-v /PATH/TO/KEY:/etc/server.key \
-v /PATH/TO/ROOT:/var/ghfs \
-p HOST_PORT:8443 \
mjpclab/ghfs
```
