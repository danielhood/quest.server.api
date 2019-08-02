quest.server.api
================
Server for quest API

Setup
-----
Install golang
`sudo apt-get install golang-go`

Install redis via docker
`docker run --name=redis -p 6379:6379 -d redis`

Or via linux subsystem on windows
`redis-server`

Code currently expects redis on default port on localhost.

Generate `server.key` and `server.pem` files in root of app for https on port 8443. If these don't exists, http is used on port 8080.

Key/Pem file generation:
```
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
```

Build server
`go build`

Run server
`quest.server.api`

