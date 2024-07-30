<h1 align=center>Command and Control HTTP Server</h1> 
<p align=center>written in Gin+Go</p>

## Compiling

compile code - `go build`
run code - `go run main.go`
add dependencies - `go get .`

## Endpoints
- `/reg` - register new agent (hostname, ip, ) - POST
- `/task/<name>` - fetch all tasks for agents (name) - GET
- `/result/<name>` - fetch/upload results from agents (name) - POST

## TODO

- [ ] Add Landing Page
- [ ] Add Login Page
- [ ] Add functionality to Dashboard Page (JS or HTMX, still debating...)
- [ ] Add functionality to remove finished tasks

## Usage

Either log in to the dashboard using the login page (WIP) or just use the CLI client provided in the `/Seraph/client` directory.

## ...
More to come in the future...
