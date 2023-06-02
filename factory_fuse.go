package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type ArcheId struct {
	Id string `json:"id"`
}

type ArcheUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id         string
	Name       string
	Email      string
	LastLogin  string
	IsLoggedIn bool
	Password   string
}

type DBCred struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	DBName string `yaml:"dbname"`
}

func ReadConfig(yaml_f string) DBCred {
	var cred DBCred

	file_b, err := os.ReadFile(yaml_f)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file_b, &cred)
	if err != nil {
		panic(err)
	}

	return cred
}

func DBInit() (*sql.DB, error) {
	/*cfg := mysql.Config{
	  User:   os.Getenv("DBUSER"),
	  Passwd: os.Getenv("DBPASS"),
	  Net:    "tcp",
	  Addr:   "127.0.0.1:3306",
	  DBName: "mine",
	}*/

	/*
	  cred := DBCred {
	    User: os.Getenv("FF_DBUSER"),
	    Pass: os.Getenv("FF_DBPASS"),
	    Host: os.Getenv("FF_DBHOST"),
	    Port: "5432",
	    DBName: "ff",
	  };*/

	cred := ReadConfig("db_cred.yaml")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cred.User, cred.Pass, cred.Host, cred.Port, cred.DBName)

	return sql.Open("postgres", connStr)
}

func FFCreateAlphaNumericString(length int) string {
	rand_gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand_gen.Intn(len(letterBytes))]
	}

	return string(b)
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
			"message":      "success",
			"serverStatus": "error",
			"error":        if_server_err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "success",
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
		c.JSON(http.StatusAccepted, gin.H{
			"error 1": err.Error(),
			"message": "usernotfound",
		})
		return
	}

	session_id := FFCreateAlphaNumericString(32)

	_, err = db.Exec("UPDATE users SET last_login = $1, is_logged_in = $2 WHERE id = $3", cur_date, "TRUE", rec_user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error 2": err.Error()})
		return
	}
	_, err = db.Exec("INSERT INTO active_sessions (session_id, user_name, user_email) VALUES ($1, $2, $3)", session_id, rec_user.Name, rec_user.Email)
	// _, err = db.Exec("INSERT INTO users (last_login, is_logged_in) VALUES (?, ?)", cur_date, "1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error 3": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "id": session_id})
}

func LogoutHandler(c *gin.Context) {
	var id ArcheId
	var err error

	if err = c.BindJSON(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user_email string

	err = db.QueryRow("SELECT user_email FROM active_sessions WHERE session_id = $1", id.Id).Scan(&user_email)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	_, err = db.Exec("UPDATE users SET is_logged_in = $1 WHERE email = $2", "FALSE", user_email)
	// _, err = db.Exec("INSERT INTO users (last_login, is_logged_in) VALUES (?, ?)", cur_date, "1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error 2": err.Error()})
		return
	}

	_, err = db.Exec("DELETE FROM active_sessions WHERE session_id = $1", id.Id)
	// _, err = db.Exec("INSERT INTO users (last_login, is_logged_in) VALUES (?, ?)", cur_date, "1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error 3": err.Error()})
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

	_, err := db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session_id := FFCreateAlphaNumericString(32)
	cur_date := time.Now().Format("2006-01-02 15:04:05")

	_, err = db.Exec("UPDATE users SET last_login = $1, is_logged_in = $2 WHERE email = $3", cur_date, "TRUE", user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error 2": err.Error()})
		return
	}
	_, err = db.Exec("INSERT INTO active_sessions (session_id, user_name, user_email) VALUES ($1, $2, $3)", session_id, user.Name, user.Email)
	// _, err = db.Exec("INSERT INTO users (last_login, is_logged_in) VALUES (?, ?)", cur_date, "1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error 3": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "id": session_id})
}

func CheckSessionHandler(c *gin.Context) {
	var reader map[string]interface{}

	var err error
	var user_name string
	var user_email string

	if err = c.BindJSON(&reader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.QueryRow("SELECT user_name, user_email FROM active_sessions WHERE session_id = $1", reader["id"]).Scan(&user_name, &user_email)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "success",
			"status":  "invalid",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"status":  "open",
		"user": gin.H{
			"name":  user_name,
			"email": user_email,
		},
	})
}

var db *sql.DB

func main() {

	var err error
	db, err = DBInit()
	if err != nil {
		panic(err.Error())
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"*"},
		// ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", MainHandler)
	r.GET("/ok", OkHandler)
	r.POST("/login", LoginHandler)
	r.POST("/logout", LogoutHandler)
	r.POST("/signup", SignUpHandler)
	r.POST("/checksession", CheckSessionHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
