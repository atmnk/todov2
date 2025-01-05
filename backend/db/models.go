package db

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB = nil

type TodoStatus string

const (
	StatusActive TodoStatus = "Active"
	StatusDone   TodoStatus = "Done"
)

func DB() *gorm.DB {
	if db == nil {
		var err error
		if db, err = gorm.Open(sqlite.Open("file:todov2?mode=memory&cache=shared"), &gorm.Config{}); err != nil {
			panic(err)
		} else {
			return db
		}
	} else {
		return db
	}
}

type User struct {
	Id        uint      `json:"id" grom:"primaryKey"`
	Username  string    `json:"name"`
	Password  string    `json:"password"`
	Todos     []Todo    `json:"todos" gorm:"foreignKey:UserId"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Todo struct {
	Id        uint       `json:"id" grom:"primaryKey"`
	Title     string     `json:"title"`
	Status    TodoStatus `json:"status"`
	UserId    uint       `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Todo{})
}
func (user *User) Create(db *gorm.DB) error {
	var existing User
	if e := db.Model(user).First(&existing, "username = ?", user.Username).Error; e == nil {
		return fmt.Errorf("user already exists")
	} else {
		return db.Create(user).Error
	}
}
func (user *User) GetTodos(db *gorm.DB) ([]Todo, error) {
	if e := db.Preload("Todos").First(user).Error; e == nil {
		return user.Todos, nil
	} else {
		return nil, e
	}
}
func (user *User) CreateTodo(db *gorm.DB, title string) (*Todo, error) {
	todo := &Todo{
		Title:  title,
		Status: StatusActive,
		UserId: user.Id,
	}
	if e := todo.Create(db); e == nil {
		return todo, nil
	} else {
		return nil, e
	}
}
func (user *User) Login(db *gorm.DB) error {
	if e := db.Model(user).First(user, "username = ? AND password = ?", user.Username, user.Password).Error; e != nil {
		return fmt.Errorf("invalid credentials")
	} else {
		return nil
	}
}
func (todo *Todo) Create(db *gorm.DB) error {
	return db.Create(todo).Error
}
