<h1 align=center>Command and Control HTTP Server</h1> 
<p align=center>written in Gin+Go</p>

## Compiling

compile code - `go build`
run code - `go run main.go`
add dependencies - `go get .`

## Endpoints
- `/users` - register user or list all users - POST/GET
- `/users/<name>` - get data for specific user - GET
- `/task/<name>` - add or fetch all tasks for agents (name) - POST/GET
- `/result/<name>` - fetch/upload results from agents (name) - POST/GET

## TODO

- [ ] Add Landing Page
- [ ] Add Login Page
- [ ] Add functionality to Dashboard Page (JS or HTMX, still debating...)
- [ ] Add functionality to remove finished tasks

## Usage

Either log in to the dashboard using the login page (WIP) or just use the CLI client provided in the `/Seraph/client` directory.

## ...
More to come in the future...
