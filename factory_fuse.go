package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

type DBCred struct {
  User string
  Pass string
  Host string
  Port string
  DBName string
}

func DBInit() (*sql.DB, error) {
  /*cfg := mysql.Config{
    User:   os.Getenv("DBUSER"),
    Passwd: os.Getenv("DBPASS"),
    Net:    "tcp",
    Addr:   "127.0.0.1:3306",
    DBName: "mine",
  }*/

  cred := DBCred {
    User: os.Getenv("DBUSER"),
    Pass: os.Getenv("DBPASS"),
    Host: "localhost",
    Port: "5432",
    DBName: "ff",
  };

  connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cred.User, cred.Pass, cred.Host, cred.Port, cred.DBName);

  return sql.Open("postgres", connStr) 
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
      "error": if_server_err.Error(),
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
  row := db.QueryRow("SELECT id, name, email FROM users WHERE email = $1 AND password = $2", user.Email, user.Password)
  err = row.Scan(&rec_user.Id, &rec_user.Name, &rec_user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error 1": err.Error()})
		return
	}

  _, err = db.Exec("UPDATE users SET last_login = $1, is_logged_in = $2 WHERE id = $3", cur_date, "1", rec_user.Id)
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

  _, err = db.Exec("UPDATE users SET is_logged_in = $1 WHERE id = $2", "0", id.Id)
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
  

  result, err := db.Exec("INSERT INTO users (name, email, password, last_login) VALUES ($1, $2, $3, $4)", user.Name, user.Email, user.Password, cur_date)
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