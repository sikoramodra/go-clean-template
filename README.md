# Go Clean template

## TODO:

- improve `internal/controller/restapi`
- fix unit and integration tests

## After git clone:

```sh
make bin-deps
mkcert -cert-file nginx/certs/fullchain.pem \
       -key-file  nginx/certs/privkey.pem \
       "*.lvh.me" lvh.me localhost 127.0.0.1 ::1
```
