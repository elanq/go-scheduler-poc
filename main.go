package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"./pkg/scheduler"
)

func main() {
	var wg sync.WaitGroup

	done := make(chan bool)
	task1 := func(id string) error {
		body := fmt.Sprintf("Task %v running", id)
		log.Println(body)
		return errors.New("something happened with task " + id)
	}

	ss := make([]*scheduler.Schedule, 0)

	//lets try 1000 task
	for i := int64(1); i <= 1000; i++ {
		s := scheduler.NewSchedule(task1, time.Duration(i)*time.Millisecond)
		s.SetParam(strconv.Itoa(int(i)))
		ss = append(ss, s)
	}
	for i := range ss {
		ss[i].Start()
		wg.Add(1)
	}

	select {
	case <-done:
		wg.Done()
	}

	wg.Wait()
}
