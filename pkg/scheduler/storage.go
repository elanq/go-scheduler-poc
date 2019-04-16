package scheduler

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
)

//TaskRepository is simple key value memory to store Schedule
type TaskRepository map[string]*Schedule

//ScheduleRepository contains schedules by its execution time
type ScheduleRepository map[string][]*Schedule

var (
	repository         TaskRepository
	scheduleRepository ScheduleRepository
	mr                 sync.Mutex
)

func init() {
	repository = make(TaskRepository, 0)
	scheduleRepository = make(ScheduleRepository, 0)
}

func (r ScheduleRepository) addSchedule(s *Schedule) {
	k := hashTime(s.execTime.Unix())
	r[k] = append(r[k], s)
}

//ExpireSchedules expires schedules by execution time
func ExpireSchedules(t time.Time) {
	k := hashTime(t.Unix())
	delete(scheduleRepository, k)
}

//UpdateTime updates schedule execution time
func UpdateTime(s *Schedule) {
	scheduleRepository.addSchedule(s)
}

//Register registers task schedule to repository
func Register(s *Schedule) error {
	mr.Lock()
	if s.id == "" {
		return errors.New("Invalid id")
	}

	if s.execTime.IsZero() {
		return fmt.Errorf("task %v has no execution time", s.id)
	}

	if _, found := repository[s.id]; found {
		return errors.New("duplicate id")
	}

	repository[hash(s.id)] = s
	UpdateTime(s)

	mr.Unlock()

	return nil
}

//GetSchedules populate schedules by its schedule time
func GetSchedules(execTime int64) []*Schedule {
	k := hashTime(execTime)
	return scheduleRepository[k]
}

func hashTime(execTime int64) string {
	hash := sha1.New()
	io.WriteString(hash, strconv.FormatInt(execTime, 10))
	return fmt.Sprintf("%x", string(hash.Sum(nil)))
}

func hash(id string) string {
	hash := sha1.New()
	io.WriteString(hash, id)
	return fmt.Sprintf("%x", string(hash.Sum(nil)))
}
