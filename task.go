package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/satori/go.uuid"
)

type ErrCode int

const (
	ErrCodeNone ErrCode = iota
	ErrCodeRedisErr
	ErrCodeWrongParameter
	ErrCodeMysqlErr
)

var (
	tasks *[]Task = nil
	tLock         = make(chan int, 1)
)

type Task struct {
	ID          string    `json:"id"`
	GroupKey    string    `json:"groupKey"`
	CallBackURL string    `json:"callBackURL"`
	Index       int       `json:"index"`
	Create      time.Time `json:"create"`
	Update      time.Time `json:"update"`
	Closed      bool      `json:"closed"`

	FirstFire time.Time     `json:"firstFire"` //default will be now
	Delay     time.Duration `json:"delay"`     //seconds
	Duration  time.Duration `json:"duration"`  //seconds
	Repeat    bool          `json:"repeat"`
}

func init() {
	tLock <- 1
	/**
	 * create table if not exist
	 */
	if false == db.HasTable(&Task{}) {
		db.CreateTable(&Task{})
	}

	/**
	 * load task from db to mem, and reschedule with timer
	 */
	go func() {
		LoadTasks()
		scheduleTasksTimer(tasks)
	}()
}

func (t *Task) InitBase() {
	id, err := uuid.NewV4()
	if nil == err {
		t.ID = id.String()
	}
	t.Create = time.Now()
	t.Update = t.Create
	if t.FirstFire.Unix() == 0 {
		t.FirstFire = time.Now()
	}
}

func LoadTasks() error {
	var list []Task
	err := db.Model(&Task{}).Where("closed = ?", false).Find(&list).Error
	if nil != err {
		return err
	}
	tasks = &list
	return nil
}

func TaskWithGroupKey(groupKey string) []*Task {
	var list []*Task
	for _, t := range *tasks {
		if t.GroupKey == groupKey {
			list = append(list, &t)
		}
	}
	return list
}

func TaskWithID(id string) *Task {
	for _, t := range *tasks {
		if t.ID == id {
			return &t
		}
	}
	return nil
}

func (t *Task) Schedule() error {
	if len(t.ID) == 0 {
		t.InitBase()
	}

	err := db.Create(t).Error
	if nil != err {
		return err
	}

	<-tLock
	list := append(*tasks, *t)
	tasks = &list
	tLock <- 1
	scheduleTaskTimerIfNeeded(t)

	return nil
}

func (t *Task) Cancel() error {
	/**
	 * update mysql, remove from mem, and remove related timer
	 * */
	t.Closed = true
	t.Update = time.Now()
	err := db.Save(t).Error
	if nil != err {
		return err
	}

	<-tLock
	list := []Task{}
	for _, tmp := range *tasks {
		if t.ID != tmp.ID {
			list = append(list, tmp)
		}
	}
	tasks = &list
	tLock <- 1

	cancelTaskTimerIfNeeded(t)

	return nil
}

func (t *Task) Fire() {
	form := url.Values{"index": {strconv.Itoa(t.Index)}}
	fmt.Fprintf(gin.DefaultWriter, "task fired with %v \n %v \n", t, form)
	resp, err := http.PostForm(t.CallBackURL, form)
	if err != nil {
		fmt.Fprintln(gin.DefaultWriter, err)
		rescheduleSingleFiredTaskTimer(t)
	} else {
		t.Cancel()
		fmt.Fprintln(gin.DefaultWriter, resp)
	}
}
