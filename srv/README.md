<h1 align=center>Command and Control HTTP Server</h1> 
<p align=center>written in Gin+Go</p>

## Compiling

compile code - `go build`
run code - `go run src/main.go`

## Endpoints
`/reg` - register new agent (hostname, ip, ) - POST
`/task/<name>` - fetch all tasks for agents (name) - GET
`/result/<name>` - fetch/upload results from agents (name) - POST

## ...

More to come in the future...
