package scheduler

import (
	"log"
	"time"
)

func (scheduler *Scheduler) ImageTextDetect() {
	scheduler.ScheduleAtFixedRate("imageLabelDetect", func() {
		log.Println("Image label detect starting...")
		scheduler.mediaService.DetectAndSaveImageText()
		log.Println("Image label detect finished...")
	}, 20*time.Second)
}
