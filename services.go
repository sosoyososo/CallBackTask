package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddTask(c *gin.Context) {
	var t Task
	if nil != c.BindJSON(&t) {
		c.JSON(http.StatusOK, ParameterErrResult())
	} else {
		err := t.Schedule()
		if err == nil {
			c.JSON(http.StatusOK, SucceedResult(t))
		} else {
			c.JSON(http.StatusOK, DBErrResult())
		}
	}
}

func CancelTask(c *gin.Context) {
	id, exist := c.GetPostForm("id")
	if !exist || len(id) == 0 {
		c.JSON(http.StatusOK, ParameterErrResult())
		return
	}

	t := TaskWithID(id)
	if t == nil {
		c.JSON(http.StatusOK, SucceedResult(nil))
	} else {
		err := t.Cancel()
		if nil == err {
			c.JSON(http.StatusOK, SucceedResult(t))
		} else {
			c.JSON(http.StatusOK, DBErrResult())
		}
	}
}

func ListTask(c *gin.Context) {
	groupKey, exist := c.GetPostForm("group")
	if !exist || len(groupKey) == 0 {
		c.JSON(http.StatusOK, ParameterErrResult())
		return
	}
	c.JSON(http.StatusOK, SucceedResult(TaskWithGroupKey(groupKey)))
}
