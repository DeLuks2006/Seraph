<h1 align=center>Command and Control HTTP Server</h1> 
<p align=center>written in Gin+Go</p>

## Compiling

compile code - `go build`
run code - `go run src/main.go`
add dependencies - `go get .`

## Endpoints
`/reg` - register new agent (hostname, ip, ) - POST
`/task/<name>` - fetch all tasks for agents (name) - GET
`/result/<name>` - fetch/upload results from agents (name) - POST

## TODO

- [ ] Add mini-config file to make logon less annoying
- [ ] Add functionality to remove finished tasks

## ...

More to come in the future...
