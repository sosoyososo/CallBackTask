package main

import (
	"time"
)

type Timer struct {
	task   *Task
	timer  *time.Timer
	ticker *time.Ticker
}

var (
	timers    = []Timer{}
	timerLock = make(chan int, 1)
)

func init() {
	timerLock <- 1
}

func (t *Task) ScheduledItem() *Timer {
	for _, tmp := range timers {
		if tmp.task.ID == t.ID {
			return &tmp
		}
	}
	return nil
}

func scheduleTasksTimer(list *[]Task) {
	for _, t := range *list {
		scheduleTaskTimerIfNeeded(&t)
	}
}

/**
 * single fire task : fired on delay + duration
 * repeat fire task : first fired on delay, then repeatly fired on new duration end
 */
func scheduleTaskTimerIfNeeded(t *Task) {
	if t.ScheduledItem() == nil {
		var timer *time.Timer = nil
		item := Timer{task: t}
		if t.Repeat {
			timer = time.AfterFunc(t.Delay, func() {
				ticker := time.NewTicker(t.Duration)
				t.ScheduledItem().ticker = ticker
				go func() {
					for range ticker.C {
						t.Fire()
					}
				}()
			})
			item.timer = timer
		} else {
			timer = time.AfterFunc(t.Delay+t.Duration, func() {
				t.Fire()
			})
			item.timer = timer
		}

		<-timerLock
		timers = append(timers, item)
		timerLock <- 1
	}
}

func cancelTaskTimerIfNeeded(t *Task) {
	item := t.ScheduledItem()
	item.timer.Stop()

	<-timerLock
	list := []Timer{}
	for _, tmp := range timers {
		if tmp.task.ID != t.ID {
			list = append(list, tmp)
		} else {
			tmp.timer.Stop()
			if tmp.ticker != nil {
				tmp.ticker.Stop()
			}
		}
	}
	timers = list
	timerLock <- 1
}