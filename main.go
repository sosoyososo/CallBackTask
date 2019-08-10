package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

type Config struct {
	MysqlHost     string
	MysqlUserName string
	MysqlUserPswd string

	ServerPort string
}

const (
	MysqlDBName = "NetTaskService"
)

var (
	db   *gorm.DB
	conf Config
)

func init() {
	fileName := "./conf-dev.json"
	if gin.ReleaseMode == gin.Mode() {
		fileName = "./conf.json"
	}
	f, err := os.Open(fileName)
	if nil != err {
		panic(err)
	}
	buf, err := ioutil.ReadAll(f)
	if nil != err {
		panic(err)
	}
	err = json.Unmarshal(buf, &conf)
	if nil != err {
		panic(err)
	}

	urlstr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.MysqlUserName,
		conf.MysqlUserPswd,
		conf.MysqlHost,
		MysqlDBName)
	mysql, err := gorm.Open("mysql", urlstr)
	if nil != err {
		panic(err)
	}
	db = mysql
}

func main() {
	port := conf.ServerPort
	if len(port) == 0 {
		port = ":8080"
	}

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.Any("listTask", ListTask)
	r.Any("addTask", AddTask)
	r.Any("cancelTask", CancelTask)

	r.Run(port)
}
