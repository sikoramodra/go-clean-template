# Go Clean template

## TODO:

- better supertokens configuration
- check tracing, metrics, swagger+supertokens
- improve `internal/controller/restapi`
- `ervone` -> `sikoramodra`
- fix unit and integration tests

## After git clone:

1.

```sh
mkcert -cert-file nginx/certs/fullchain.pem \
       -key-file  nginx/certs/privkey.pem \
       "*.lvh.me" lvh.me app.lvh.me jaeger.lvh.me localhost 127.0.0.1 ::1
```
