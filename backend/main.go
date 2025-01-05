package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type TodoStatus string

const (
	StatusActive TodoStatus = "Active"
	StatusDone   TodoStatus = "Done"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// var db = make(map[string]string)
type Todo struct {
	Id     uint       `json:"id"`
	Title  string     `json:"title"`
	Status TodoStatus `json:"status"`
}
type CreateTodo struct {
	Title string `json:"title"`
}
type UpdateTodo struct {
	Title  *string     `json:"title"`
	Status *TodoStatus `json:"status"`
}

var db = make(map[uint]Todo)
var mu sync.Mutex

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(CORSMiddleware())
	// // Ping test
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "pong")
	// })
	r.GET("/todo", func(ctx *gin.Context) {
		resp := []Todo{}
		mu.Lock()
		for _, todo := range db {
			resp = append(resp, todo)
		}
		mu.Unlock()
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": resp, "message": "success"})
	})
	r.POST("/todo", func(ctx *gin.Context) {
		var todoBody CreateTodo
		ctx.ShouldBindBodyWithJSON(&todoBody)
		mu.Lock()
		var maxKey uint = 0
		for key := range db {
			if maxKey < key {
				maxKey = key
			}
		}
		maxKey = maxKey + 1
		todo := Todo{
			Id:     maxKey,
			Title:  todoBody.Title,
			Status: StatusActive,
		}
		db[todo.Id] = todo
		mu.Unlock()
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": todo, "message": "Todo Created"})
	})
	r.PATCH("/todo/:id", func(ctx *gin.Context) {
		if Id, e := GetId(ctx); e == nil {
			var todoBody UpdateTodo
			if e := ctx.ShouldBindBodyWithJSON(&todoBody); e == nil {
				mu.Lock()
				original := db[Id]
				title := original.Title
				status := original.Status
				if todoBody.Title != nil {
					title = *todoBody.Title
				}
				if todoBody.Status != nil {
					status = *todoBody.Status
				}
				newTodo := Todo{
					Id:     Id,
					Title:  title,
					Status: status,
				}
				db[Id] = newTodo
				mu.Unlock()
				ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": newTodo, "message": "Todo Updated"})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "message": e.Error()})
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "message": e.Error()})
		}
	})
	r.DELETE("/todo/:id", func(ctx *gin.Context) {
		if Id, e := GetId(ctx); e == nil {

			mu.Lock()
			delete(db, Id)
			mu.Unlock()
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Todo Deleted"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "message": e.Error()})
		}

	})

	// // Get user value
	// r.GET("/user/:name", func(c *gin.Context) {
	// 	user := c.Params.ByName("name")
	// 	value, ok := db[user]
	// 	if ok {
	// 		c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
	// 	} else {
	// 		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
	// 	}
	// })

	// // Authorized group (uses gin.BasicAuth() middleware)
	// // Same than:
	// // authorized := r.Group("/")
	// // authorized.Use(gin.BasicAuth(gin.Credentials{
	// //	  "foo":  "bar",
	// //	  "manu": "123",
	// //}))
	// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	// 	"foo":  "bar", // user:foo password:bar
	// 	"manu": "123", // user:manu password:123
	// }))

	// /* example curl for /admin with basicauth header
	//    Zm9vOmJhcg== is base64("foo:bar")

	// 	curl -X POST \
	//   	http://localhost:8080/admin \
	//   	-H 'authorization: Basic Zm9vOmJhcg==' \
	//   	-H 'content-type: application/json' \
	//   	-d '{"value":"bar"}'
	// */
	// authorized.POST("admin", func(c *gin.Context) {
	// 	user := c.MustGet(gin.AuthUserKey).(string)

	// 	// Parse JSON
	// 	var json struct {
	// 		Value string `json:"value" binding:"required"`
	// 	}

	// 	if c.Bind(&json) == nil {
	// 		db[user] = json.Value
	// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	// 	}
	// })

	return r
}
func GetId(ctx *gin.Context) (uint, error) {
	idString := ctx.Param("id")
	if id, e := strconv.ParseUint(idString, 10, 64); e == nil {
		return uint(id), nil
	} else {
		return 0, fmt.Errorf("invalid id")
	}
}
func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
