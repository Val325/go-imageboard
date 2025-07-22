package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Post  string `json:"post"`
	Time  string `json:"time"`
}

// Save to DB
var posts = []Post{}

func getPosts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, posts)
}

func postPosts(c *gin.Context) {
	var newPost Post

	if err := c.ShouldBind(&newPost); err != nil {
		return
	}

	//add to DB
	db, err := sql.Open("sqlite3", "./posts.db")
	dt := time.Now()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO posts(title, post, time) VALUES(?, ?, ?)", c.Request.PostForm["title"][0], c.Request.PostForm["post"][0], dt.Format("01-02-2006 15:04:05"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("New user inserted successfully")

	c.IndentedJSON(http.StatusCreated, posts)
}

func getPostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for _, a := range posts {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "post not found"})
}

func app(c *gin.Context) {
	posts = nil
	db, err := sql.Open("sqlite3", "./posts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, title, post, time FROM posts;")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		var post string
		var time string
		err = rows.Scan(&id, &title, &post, &time)
		if err != nil {
			log.Fatal(err)
		}
		var newPost Post
		newPost = Post{ID: id, Title: title, Post: post, Time: time}
		posts = append(posts, newPost)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	c.HTML(http.StatusOK, "all-posts.tmpl", posts)
}

func main() {
	db, err := sql.Open("sqlite3", "./posts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        title TEXT,
		post TEXT,
		time TEXT
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table 'posts' created successfully")

	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("template/*")

	router.GET("/", app)
	router.GET("/api/posts", getPosts)
	router.POST("/api/posts", postPosts)
	router.GET("/api/posts/:id", getPostByID)

	router.Run("localhost:8080")
}
