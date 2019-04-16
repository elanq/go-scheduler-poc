package main

import (
	"fmt"
	"log"
	"time"

	"../pkg/scheduler"
)

func main() {
	forever := make(chan bool, 1)

	task1 := func(id string) error {
		body := fmt.Sprintf("Task %v running", id)
		log.Println(body)
		//return errors.New("something happened with task " + id)
		return nil
	}
	t := time.Now().Add(5 * time.Second)

	s := scheduler.NewSchedule("123", task1, t).
		SetTaskInterval(2 * time.Second).
		SetParam("hehehe")

	scheduler.Register(s)

	//Scheduler thread
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				dispatch(time.Now())
			}
		}
	}()

	<-forever
}

func dispatch(t time.Time) {
	schedules := scheduler.GetSchedules(t.Unix())
	if len(schedules) == 0 {
		return
	}

	for i := range schedules {
		go func(s *scheduler.Schedule) {
			err := s.Dispatch()
			if err != nil {
				log.Println(err)
				s.Stop()
			}
		}(schedules[i])
	}

	scheduler.ExpireSchedules(t)
}
