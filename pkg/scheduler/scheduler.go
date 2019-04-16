package scheduler

import (
	"sync"
	"time"
)

//Task is task
type Task func(id string) error

//Schedule is schedule
type Schedule struct {
	mu                sync.Mutex
	id                string
	done              bool
	task              Task
	param             string
	execTime          time.Time
	nextExecTime      time.Time
	recurringInterval time.Duration
}

//NewSchedule creates new schedule. By default scheduled task is not recurring
func NewSchedule(taskID string, t Task, execTime time.Time) *Schedule {
	return &Schedule{
		id:                taskID,
		task:              t,
		done:              false,
		recurringInterval: 0 * time.Second,
		param:             "",
		execTime:          execTime,
	}
}

//SetTaskInterval sets interval duration for task
//this method also change nextExecTime value by current execTime
func (s *Schedule) SetTaskInterval(dur time.Duration) *Schedule {
	s.recurringInterval = dur
	s.nextExecTime = s.execTime.Add(dur)
	return s
}

//SetParam sets param
func (s *Schedule) SetParam(param string) *Schedule {
	s.param = param
	return s
}

//Stop stops schedule
func (s *Schedule) Stop() {
	s.mu.Lock()
	s.done = true
	s.mu.Unlock()
}

//HashID returns hashed ID of schedule job. use this method to store in repository
func (s *Schedule) HashID() string {
	return hash(s.id)
}

//Dispatch executes task from scheduler
func (s *Schedule) Dispatch() error {
	if s.done {
		return nil
	}

	err := s.task(s.param)
	if err != nil {
		return err
	}
	//set status to dispatched if has no recurring status
	if !s.nextExecTime.IsZero() {
		s.execTime = s.nextExecTime
		s.nextExecTime = s.execTime.Add(s.recurringInterval)

		UpdateTime(s)
	}
	return nil
}

//Start starts Schedule
func (s *Schedule) Start() {
	//	errChan := make(chan error, 1)
	//	done := make(chan bool, 1)
	//	go func() {
	//		ticker := time.NewTicker(s.timer)
	//		for {
	//			select {
	//			case e := <-errChan:
	//				log.Println(e)
	//			case <-ticker.C:
	//				//handle double done
	//				if s.done {
	//					return
	//				}
	//
	//				err := s.task(s.param)
	//				if err != nil {
	//					errChan <- err
	//					s.SetDone()
	//					done <- true
	//				}
	//			}
	//		}
	//	}()
}
