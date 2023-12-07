package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"net/http"
)

var db *sqlx.DB

func main() {
	// init db
	if err := initDB(); err != nil {
		fmt.Println("connect mysqlq failed", err)
		panic(err)
	}

	server := gin.Default()

	server.POST("/login", loginHandler)

	server.Run()
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func loginHandler(context *gin.Context) {
	// 1.从请求中获取用户的请求数据
	// 要么是form表单提交，要么是json格式提交。

	var reqDate Login
	if err := context.ShouldBind(&reqDate); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "request option error",
		})
	}

	fmt.Printf("reqData: %#v \n", reqDate)
	//Scontext.JSON(http.StatusOK, reqDate)
	//fmt.Println("////")
	// 1.对数据进行校验

	if u, err := QueryMysqlUser(reqDate.Username, reqDate.Password); err == nil {
		fmt.Println(u)
		context.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "hello " + reqDate.Username,
			"data": u,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "username or password error",
		})
	}
	// 1.返回响应
}

func initDB() (err error) {
	dsn := "root:123456@tcp(192.168.87.135:3306)/myblog?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

type User struct {
	// - 是忽略id字段
	Id       int    `db:"id" json:"-"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	// Desc 使用omitempty，是指该字段有值的时候显示，无值的时候忽略
	Desc string `json:"desc,omitempty"`
}

func QueryMysqlUser(username, password string) (*User, error) {
	// 查库
	sqlstr := "select id, username, password from user where username=? and password=?"
	var u User
	if err := db.Get(&u, sqlstr, username, password); err != nil {
		fmt.Printf("[get failed , err : %v]\n", err)
		return nil, err
	}

	return &u, nil

}
