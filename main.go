package main

import (
	"blog-fanchiikawa-service/greetings"
	"blog-fanchiikawa-service/user"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetPrefix("blog-fanchiikawa-service: ")
	log.SetFlags(0)

	router := gin.Default()
	router.GET("/users", getUserList)
	router.POST("/login", login)

	router.Run("localhost:8080")
}

func getUserList(c *gin.Context) {
	users := user.GetUserList()
	c.IndentedJSON(http.StatusOK, users)
}

func login(c *gin.Context) {
	var newUser user.User

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	newUser.ID = len(user.UserList) + 1
	newUser.CreatedAt = time.Now()

	user.SaveUser(newUser)

	message, err := greetings.Hello(newUser.Nickname)
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, message)
}
