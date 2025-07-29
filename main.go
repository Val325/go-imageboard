package main

import (
	"database/sql"
	"fmt"
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

var postsPerPage = 5
var pagesPagination = 5
var maxPage = 10

// Save to DB
var posts = []Post{}

func getPosts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, posts)
	c.Redirect(http.StatusMovedPermanently, "/")

}

/*
	    CREATE TABLE IF NOT EXISTS posts (
	        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			mainpost INTEGER,
			isMain INTEGER,
	        title TEXT,
			post TEXT,
			time TEXT
	    );
*/
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

	_, err = db.Exec("INSERT INTO posts(mainpost, isMain, title, post, time) VALUES(?, ?, ?, ?, ?)", 0, 1, c.Request.PostForm["title"][0], c.Request.PostForm["post"][0], dt.Format("01-02-2006 15:04:05"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("New user inserted successfully")

	var idMax int
	err_max := db.QueryRow("SELECT MAX(id) FROM posts;").Scan(&idMax)
	if err_max != nil {
		log.Fatal(err_max)
	}
	if idMax >= maxPage*postsPerPage {
		//delete min posts
		var idMin int
		err_min := db.QueryRow("SELECT MIN(id) FROM posts;").Scan(&idMin)
		if err_min != nil {
			log.Fatal(err_min)
		}
		_, err = db.Exec("DELETE FROM posts WHERE id = ?;", idMin)
		if err != nil {
			log.Fatal(err)
		}
	}

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

func startpage(c *gin.Context) {
	var page = 1
	posts = nil
	db, err := sql.Open("sqlite3", "./posts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, mainpost, isMain, title, post, time FROM posts;")
	if err != nil {
		log.Fatal(err)
	}
	//SELECT MAX(id) FROM posts;
	defer rows.Close()
	for rows.Next() {
		var id int
		var mainpost int
		var isMain int
		var title string
		var post string
		var time string

		err = rows.Scan(&id, &mainpost, &isMain, &title, &post, &time)
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

	var idMax int
	var idMin int
	//var remainderPosts int
	var isIsExist bool
	err_exist := db.QueryRow("SELECT EXISTS(SELECT MIN(id) FROM posts);").Scan(&isIsExist)
	if err_exist != nil {
		log.Fatal(err_exist)
	}
	fmt.Print("isExist: ", isIsExist, "\n")
	if isIsExist == false {
		page = 0
		postsPerPage = 0
		pagesPagination = 0

		c.HTML(http.StatusOK, "all-posts.tmpl", gin.H{
			"posts":      posts[(page*postsPerPage)-postsPerPage : page*postsPerPage],
			"nums_pages": CalculateRangeArray(1, pagesPagination+1),
			"page":       page,
		})
	} else {
		//err_max := db.QueryRow("SELECT MAX(id) FROM posts;").Scan(&idMax)
		db.QueryRow("SELECT MAX(id) FROM posts;").Scan(&idMax)
		//if err_max != nil {
		//	log.Fatal(err_max)
		//}
		fmt.Print("MAX(id): ", idMax, "\n")
		//err_min := db.QueryRow("SELECT MIN(id) FROM posts;").Scan(&idMin)
		db.QueryRow("SELECT MIN(id) FROM posts;").Scan(&idMin)
		//if err_min != nil {
		//	log.Fatal(err_min)
		//}
		fmt.Print("MIN(id): ", idMin, "\n")
		//remainderPosts = idMax % 5
		//fmt.Println("id max: ", idMax, "\n")
		//fmt.Println("id min: ", idMin, "\n")
		//fmt.Println("page: ", page, "\n")
		//fmt.Println("postsPerPage: ", postsPerPage, "\n")
		//fmt.Println("remainder posts: ", 5-remainderPosts, "\n")
		pagesPagination := CalculatePages(len(posts), 5)
		if page > pagesPagination {
			page = pagesPagination
		}
		if page <= 0 {
			page = 1
		}

		var pagesMax int
		pagesMax = (page * postsPerPage)
		if (page * postsPerPage) > len(posts) {
			pagesMax = len(posts)
		}
		//(page * postsPerPage)
		c.HTML(http.StatusOK, "all-posts.tmpl", gin.H{
			"posts":      posts[(page*5)-postsPerPage : pagesMax],
			"nums_pages": CalculateRangeArray(1, pagesPagination+1),
			"page":       page,
		})

	}

}

func app(c *gin.Context) {
	var page int

	if isNumeric(c.Param("page")) {
		var pageNum int
		pageNum, err_page := strconv.Atoi(c.Param("page"))
		if err_page != nil {
			log.Fatal(err_page)
		}
		page = pageNum
	} else {
		page = 1
	}

	posts = nil
	db, err := sql.Open("sqlite3", "./posts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, mainpost, isMain, title, post, time FROM posts;")
	if err != nil {
		log.Fatal(err)
	}
	//SELECT MAX(id) FROM posts;
	defer rows.Close()
	for rows.Next() {
		var id int
		var mainpost int
		var isMain int
		var title string
		var post string
		var time string

		err = rows.Scan(&id, &mainpost, &isMain, &title, &post, &time)
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

	var idMax int
	var idMin int
	//var remainderPosts int
	var isIsExist bool
	err_exist := db.QueryRow("SELECT EXISTS(SELECT MIN(id) FROM posts);").Scan(&isIsExist)
	if err_exist != nil {
		log.Fatal(err_exist)
	}

	//fmt.Print("len: ", len(posts), "\n")
	//nums_page := CalculatePages(len(posts), 5)
	//fmt.Print("amount pages: ", nums_page, "\n")
	//c.HTML(http.StatusOK, "all-posts.tmpl", posts[(page*5)-5:page*5])
	fmt.Print("isExist: ", isIsExist, "\n")
	if isIsExist == false {
		page = 0
		postsPerPage = 0
		pagesPagination = 0

		c.HTML(http.StatusOK, "all-posts.tmpl", gin.H{
			"posts":      posts[(page*postsPerPage)-postsPerPage : page*postsPerPage],
			"nums_pages": CalculateRangeArray(1, pagesPagination+1),
			"page":       page,
		})
	} else {
		db.QueryRow("SELECT MAX(id) FROM posts;").Scan(&idMax)
		db.QueryRow("SELECT MIN(id) FROM posts;").Scan(&idMin)
		//err_max := db.QueryRow("SELECT MAX(id) FROM posts;").Scan(&idMax)
		//if err_max != nil {
		//	log.Fatal(err_max)
		//}
		//err_min := db.QueryRow("SELECT MIN(id) FROM posts;").Scan(&idMin)
		//if err_min != nil {
		//	log.Fatal(err_min)
		//}
		//remainderPosts = idMax % 5
		//fmt.Println("id max: ", idMax, "\n")
		//fmt.Println("id min: ", idMin, "\n")
		//fmt.Println("page: ", page, "\n")
		//fmt.Println("postsPerPage: ", postsPerPage, "\n")
		//fmt.Println("remainder posts: ", 5-remainderPosts, "\n")
		pagesPagination := CalculatePages(len(posts), 5)
		if page > pagesPagination {
			page = pagesPagination
		}
		if page <= 0 {
			page = 1
		}

		var pagesMax int
		pagesMax = (page * postsPerPage)
		if (page * postsPerPage) > len(posts) {
			pagesMax = len(posts)
		}
		//(page * postsPerPage)
		c.HTML(http.StatusOK, "all-posts.tmpl", gin.H{
			"posts":      posts[(page*5)-postsPerPage : pagesMax],
			"nums_pages": CalculateRangeArray(1, pagesPagination+1),
			"page":       page,
		})
	}

}

func postSubPosts(c *gin.Context) {
	var postsNum int
	postsNum, err_post := strconv.Atoi(c.Param("id"))
	if err_post != nil {
		log.Fatal(err_post)
	}

	//add to DB
	db, err := sql.Open("sqlite3", "./posts.db")
	dt := time.Now()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO posts(mainpost, isMain, title, post, time) VALUES(?, ?, ?, ?, ?)", postsNum, 0, c.Request.PostForm["title"][0], c.Request.PostForm["post"][0], dt.Format("01-02-2006 15:04:05"))

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
		mainpost INTEGER,
		isMain INTEGER,	
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

	router.GET("/", startpage)
	router.GET("/:page", app)
	router.GET("/api/posts", getPosts)
	router.POST("/api/posts", postPosts)
	router.POST("/api/posts/:id", postSubPosts)
	router.GET("/api/posts/:id", getPostByID)

	router.Run("localhost:8080")
}
