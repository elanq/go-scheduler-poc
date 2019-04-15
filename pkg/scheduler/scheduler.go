package scheduler

import (
	"log"
	"sync"
	"time"
)

//Task is task
type Task func(id string) error

//Schedule is schedule
type Schedule struct {
	mu    sync.Mutex
	done  bool
	task  Task
	param string
	timer time.Duration
}

//NewSchedule creates new schedule
func NewSchedule(t Task, dur time.Duration) *Schedule {
	return &Schedule{
		task:  t,
		done:  false,
		timer: dur,
		param: "",
	}
}

//SetParam sets param
func (s *Schedule) SetParam(param string) {
	s.param = param
}

//SetTask sets task
func (s *Schedule) SetTask(t Task) {
	s.task = t
}

//SetDone status
func (s *Schedule) SetDone() {
	s.mu.Lock()
	s.done = true
	s.mu.Unlock()
}

//Start starts Schedule
func (s *Schedule) Start() {
	errChan := make(chan error, 1)
	done := make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(s.timer)
		for {
			select {
			case e := <-errChan:
				log.Println(e)
			case <-ticker.C:
				//handle double done
				if s.done {
					return
				}

				err := s.task(s.param)
				if err != nil {
					errChan <- err
					s.SetDone()
					done <- true
				}
			}
		}
	}()
}
