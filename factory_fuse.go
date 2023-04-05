package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type ArcheId struct {
  Id string `json:"id"`
}

type ArcheUser struct {
  Name string `json:"name"`
  Email string `json:"email"`
  Password string `json:"password"`
}

type User struct {
  Id string
  Name string
  Email string
  LastLogin string
  IsLoggedIn bool
  Password string
}

func DBInit() (*sql.DB, error) {
  cfg := mysql.Config{
    User:   os.Getenv("DBUSER"),
    Passwd: os.Getenv("DBPASS"),
    Net:    "tcp",
    Addr:   "127.0.0.1:3306",
    DBName: "mine",
  }

  return sql.Open("mysql", cfg.FormatDSN()) 
}

func MainHandler(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{
    "message": "success",
  })
}

func OkHandler(c *gin.Context) {
  if_server_err := db.Ping()
  if if_server_err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": "success",
      "serverStatus": "error",
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "message": "success",
    "serverStatus": "ok",
  })
}

func LoginHandler(c *gin.Context) {
  var user ArcheUser
  var err error

  if err = c.BindJSON(&user); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  var rec_user User
  cur_date := time.Now().Format("2006-01-02 15:04:05")
  row := db.QueryRow("SELECT id, name, email FROM users WHERE email = ? AND password = ?", user.Email, user.Password)
  err = row.Scan(&rec_user.Id, &rec_user.Name, &rec_user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error 1": err.Error()})
		return
	}

  _, err = db.Exec("UPDATE users SET last_login = ?, is_logged_in = ? WHERE id = ?", cur_date, "1", rec_user.Id)
  // _, err = db.Exec("INSERT INTO users (last_login, is_logged_in) VALUES (?, ?)", cur_date, "1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error 2": err.Error()})
		return
	}

  c.JSON(http.StatusOK, gin.H{"message": "success", "id": rec_user.Id})
}

func LogoutHandler(c *gin.Context) {
  var id ArcheId
  var err error

  if err = c.BindJSON(&id); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  _, err = db.Exec("UPDATE users SET is_logged_in = ? WHERE id = ?", "0", id.Id)
  // _, err = db.Exec("INSERT INTO users (last_login, is_logged_in) VALUES (?, ?)", cur_date, "1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error 2": err.Error()})
		return
	}

  c.JSON(http.StatusOK, gin.H{"message": "success"})
}


func SignUpHandler(c *gin.Context) {
  var user ArcheUser
  if err := c.BindJSON(&user); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  cur_date := time.Now().Format("2006-01-02 15:04:05")
  

  result, err := db.Exec("INSERT INTO users (name, email, password, last_login) VALUES (?, ?, ?, ?)", user.Name, user.Email, user.Password, cur_date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

  id, _ := result.LastInsertId()
  c.JSON(http.StatusOK, gin.H{"message": "success", "id": id})
}

var db *sql.DB

func main() {

  var err error
  db, err = DBInit()
  if err != nil {
    panic(err.Error())
  }

  r := gin.Default()
  r.GET("/", MainHandler)
  r.GET("/ok", OkHandler)
  r.POST("/login", LoginHandler)
  r.POST("/logout", LogoutHandler)
  r.POST("/signup", SignUpHandler)
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}