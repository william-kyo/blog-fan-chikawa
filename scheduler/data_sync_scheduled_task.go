package scheduler

import (
	"log"
	"time"
)

func (scheduler *Scheduler) DataSync() {
	scheduler.ScheduleAtFixedRate("dataSync", func() {
		log.Println("Data sync...")
	}, 5*time.Second)
}
