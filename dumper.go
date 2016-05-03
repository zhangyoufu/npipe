package npipe

import (
	"fmt"
	"time"
)

type logEvent struct  {
		Name string
		Timestamp time.Time
		}
		
var eventBuffer []logEvent
const bufSize = 500
var eventIndex int		
var logChan chan logEvent
var finished chan int


func init(){
	logChan = make(chan logEvent, 100)
	finished = make(chan int)
	go func() {
		for {
			select {
				case e := <-logChan:
					if len(eventBuffer) < bufSize {
						eventBuffer = append(eventBuffer, e)
					} else {
						eventBuffer[eventIndex % bufSize] = e
					}
					eventIndex = eventIndex + 1
				case <-finished:
					for i := 0; i < len(eventBuffer); i++ {
						currentIndex := (eventIndex + i) % len(eventBuffer)
						fmt.Println(eventBuffer[currentIndex].Timestamp.Format("Mon Jan 2 15:04:05.0000 -0700 MST 2006"), eventBuffer[currentIndex].Name) 
					}
					finished <- 0
					return
			}
		}
	}()
}

func AddEvent(desc string) {
	logChan <- logEvent{desc, time.Now()}
}

func Dump() {
	finished <- 1
	<- finished
}
		 