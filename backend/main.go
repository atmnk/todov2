package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qalens/todov2/db"
	"github.com/qalens/todov2/service"
)

func Authorize(ctx *gin.Context) {
	var token string
	cookie, err := ctx.Cookie("token")

	authorizationHeader := ctx.Request.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) != 0 && fields[0] == "Bearer" {
		token = fields[1]
	} else if err == nil {
		token = cookie
	}

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Unauthorized: " + err.Error()})
		return
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.Set("currentUser", db.User{
		Id:       uint(claims["id"].(float64)),
		Username: claims["username"].(string),
	})
	ctx.Set("claims", claims)
	ctx.Next()
}
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

type CreateTodo struct {
	Title string `json:"title"`
}
type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UpdateTodo struct {
	Title  *string        `json:"title"`
	Status *db.TodoStatus `json:"status"`
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(CORSMiddleware())
	// // Ping test
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "pong")
	// })
	r.POST("/user", func(ctx *gin.Context) {
		var createUserBody CreateUser
		ctx.ShouldBindBodyWithJSON(&createUserBody)
		user := &db.User{
			Username: createUserBody.Username,
			Password: createUserBody.Password,
		}
		if e := user.Create(db.DB()); e == nil {
			if token, err := service.GenerateToken(user); err == nil {
				ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": token, "message": "User Created"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "data": err.Error(), "message": "Internal server error"})
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "data": e.Error(), "message": "Bad request"})
		}
	})
	r.POST("/user/login", func(ctx *gin.Context) {
		var createUserBody CreateUser
		ctx.ShouldBindBodyWithJSON(&createUserBody)
		user := &db.User{
			Username: createUserBody.Username,
			Password: createUserBody.Password,
		}
		if e := user.Login(db.DB()); e == nil {
			if token, err := service.GenerateToken(user); err == nil {
				ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": token, "message": "User logged in"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "data": err.Error(), "message": "Internal server error"})
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "data": e.Error(), "message": "Bad request"})
		}
	})
	r.GET("/todo", Authorize, func(ctx *gin.Context) {
		user := ctx.MustGet("currentUser").(db.User)
		if todos, e := user.GetTodos(db.DB()); e == nil {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": todos, "message": "success"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "data": e.Error(), "message": "Bad request"})
		}
	})
	r.POST("/todo", Authorize, func(ctx *gin.Context) {
		var todoBody CreateTodo
		ctx.ShouldBindBodyWithJSON(&todoBody)
		user := ctx.MustGet("currentUser").(db.User)
		if todo, e := user.CreateTodo(db.DB(), todoBody.Title); e == nil {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": todo, "message": "Todo created"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "data": e.Error(), "message": "Bad request"})
		}
	})
	// 	resp := []Todo{}
	// 	mu.Lock()
	// 	for _, todo := range db {
	// 		resp = append(resp, todo)
	// 	}
	// 	mu.Unlock()
	// 	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": resp, "message": "success"})
	// })
	// r.POST("/todo", func(ctx *gin.Context) {
	// 	var todoBody CreateTodo
	// 	ctx.ShouldBindBodyWithJSON(&todoBody)
	// 	mu.Lock()
	// 	var maxKey uint = 0
	// 	for key := range db {
	// 		if maxKey < key {
	// 			maxKey = key
	// 		}
	// 	}
	// 	maxKey = maxKey + 1
	// 	todo := Todo{
	// 		Id:     maxKey,
	// 		Title:  todoBody.Title,
	// 		Status: StatusActive,
	// 	}
	// 	db[todo.Id] = todo
	// 	mu.Unlock()
	// 	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": todo, "message": "Todo Created"})
	// })
	r.PATCH("/todo/:id", func(ctx *gin.Context) {
		// if Id, e := GetId(ctx); e == nil {
		// 	var todoBody UpdateTodo
		// 	if e := ctx.ShouldBindBodyWithJSON(&todoBody); e == nil {
		// 		mu.Lock()
		// 		original := db[Id]
		// 		title := original.Title
		// 		status := original.Status
		// 		if todoBody.Title != nil {
		// 			title = *todoBody.Title
		// 		}
		// 		if todoBody.Status != nil {
		// 			status = *todoBody.Status
		// 		}
		// 		newTodo := Todo{
		// 			Id:     Id,
		// 			Title:  title,
		// 			Status: status,
		// 		}
		// 		db[Id] = newTodo
		// 		mu.Unlock()
		// 		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": newTodo, "message": "Todo Updated"})
		// 	} else {
		// 		ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "message": e.Error()})
		// 	}
		// } else {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "message": e.Error()})
		// }
	})
	r.DELETE("/todo/:id", func(ctx *gin.Context) {
		// if Id, e := GetId(ctx); e == nil {

		// 	mu.Lock()
		// 	delete(db, Id)
		// 	mu.Unlock()
		// 	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Todo Deleted"})
		// } else {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "message": e.Error()})
		// }

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
	db.Migrate(db.DB())
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
