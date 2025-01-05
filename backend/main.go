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
				ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": token, "message": "User Created"})
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
			ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": todo, "message": "Todo created"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "data": e.Error(), "message": "Bad request"})
		}
	})
	r.PATCH("/todo/:id", Authorize, func(ctx *gin.Context) {
		if Id, e := GetId(ctx); e == nil {
			var todoBody db.UpdateTodo
			ctx.ShouldBindBodyWithJSON(&todoBody)
			if e := todoBody.Validate(); e == nil {
				user := ctx.MustGet("currentUser").(db.User)
				if todo, e := user.UpdateTodo(db.DB(), Id, todoBody); e == nil {
					ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": todo, "message": "Todo updated"})
				} else {
					ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "data": e.Error(), "message": "Not found"})
				}
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "data": e.Error(), "message": "Bad request"})
			}
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "data": e.Error(), "message": "Not found"})
		}
	})
	r.DELETE("/todo/:id", Authorize, func(ctx *gin.Context) {
		if Id, e := GetId(ctx); e == nil {
			user := ctx.MustGet("currentUser").(db.User)
			if e := user.DeleteTodo(db.DB(), Id); e == nil {
				ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": Id, "message": "Todo deleted"})
			} else {
				ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "data": e.Error(), "message": "Bad request"})
			}
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "data": e.Error(), "message": "Not found"})
		}
	})
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
	r.Run(":8080")
}
