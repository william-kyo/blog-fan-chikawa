package scheduler

import (
	"log"
	"time"
)

func (scheduler *Scheduler) ImageLabelDetect() {
	scheduler.ScheduleAtFixedRate("imageLabelDetect", func() {
		log.Println("Image label detect starting...")
		scheduler.mediaService.DetectAndSaveImageLabels()
		log.Println("Image label detect finished...")
	}, 20*time.Second)
}
