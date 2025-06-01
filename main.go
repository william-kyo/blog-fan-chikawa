package main

import (
	"blog-fanchiikawa-service/db"
	"blog-fanchiikawa-service/greetings"
	"blog-fanchiikawa-service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func main() {
// 	log.SetPrefix("blog-fanchiikawa-service: ")
// 	log.SetFlags(0)

// 	db.InitDB()

// 	router := gin.Default()
// 	// 添加全局异常处理中间件
// 	router.Use(middleware.Recovery())

// 	router.GET("/users", getUserList)
// 	router.POST("/login", login)

// 	router.Run("localhost:8080")
// }

func getUserList(c *gin.Context) {
	rows, err := db.DB.Query("SELECT * FROM user limit 10")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	users := []user.User{}
	for rows.Next() {
		var u user.User
		err := rows.Scan(&u.ID, &u.Nickname, &u.Email, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, u)
	}
	c.JSON(http.StatusOK, users)
}

func login(c *gin.Context) {
	var newUser user.User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.DB.Exec("INSERT INTO user (nickname, email) VALUES (?, ?)", newUser.Nickname, newUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUser.ID = lastInsertId
	message, err := greetings.Hello(newUser.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}
