package main

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	newSqliteHandler("./test.db")
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, s.allDB())
	})
	r.POST("/user", func(c *gin.Context) {
		num, _ := strconv.Atoi(c.PostForm("discount"))
		s.addDB(c.PostForm("name"), c.PostForm("app"), num, c.PostForm("content"))
	})
	r.Run()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Header", "Content-Type, Authorization, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, DELETE, POST")
		c.Next()
	}
}

var s sqliteHandler

type sqliteHandler struct {
	db *sql.DB
}

func newSqliteHandler(filepath string) {
	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS stores (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			app TEXT,
			discount INTEGER,
			content TEXT,
			createdAt DATETIME
		);`)
	statement.Exec()
	s.db = database
}

func (s *sqliteHandler) addDB(name string, app string, discount int, content string) {
	stmt, err := s.db.Prepare("INSERT INTO stores (name,app,discount,content,createdAt) VALUES(?,?,?,?,datetime('now'))")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(name, app, discount, content)
	if err != nil {
		panic(err)
	}
}

type Store struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	App       string    `json:"app"`
	Discount  int       `json:"discount"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *sqliteHandler) allDB() []*Store {
	stores := []*Store{}
	rows, err := s.db.Query("SELECT id, name, app, discount, content, createdAt FROM stores")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var store Store
		rows.Scan(&store.ID, &store.Name, &store.App, &store.Discount, &store.Content, &store.CreatedAt)
		stores = append(stores, &store)
	}
	return stores
}
