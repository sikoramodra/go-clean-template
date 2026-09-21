# Go Clean template

## TODO:

- check swagger+supertokens
- improve `internal/controller/restapi`
- fix unit and integration tests

## After git clone:

1.

```sh
make bin-deps
mkcert -cert-file nginx/certs/fullchain.pem \
       -key-file  nginx/certs/privkey.pem \
       "*.lvh.me" lvh.me localhost 127.0.0.1 ::1
```
