package main

import (
	"fmt"
	"time"

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
	ID          string        `json:"id"`
	GroupKey    string        `json:"groupKey"`
	Delay       time.Duration `json:"delay"`
	Duration    time.Duration `json:"duration"`
	Repeat      bool          `json:"repeat"`
	CallBackURL string        `json:"callBackURL"`
	Index       int           `json:"index"`
	Create      time.Time     `json:"create"`
	Update      time.Time     `json:"update"`
	Closed      bool          `json:"closed"`
}

func init() {
	/**
	 * create table if not exist
	 */
	if false == db.HasTable(&Task{}) {
		db.CreateTable(&Task{})
	}

	/**
	 * load task from db to mem, and reschedule with timer
	 */
	LoadTasks()
	scheduleTasksTimer(tasks)
	tLock <- 1
}

func (t *Task) InitBase() {
	id, err := uuid.NewV4()
	if nil == err {
		t.ID = id.String()
	}
	t.Create = time.Now()
	t.Update = t.Create
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
	fmt.Printf("task fired with %v \n", t)
}
