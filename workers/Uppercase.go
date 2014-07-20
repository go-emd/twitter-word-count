package workers

import (
	"emd/log"
	"emd/worker"
	"strings"
)

type Uppercase struct {
	worker.Work
}

func (w Uppercase) Init() {
	for _, p := range w.Ports() {
		p.Open()
	}

	log.INFO.Println("Worker " + w.Name() + " inited.")
}

func (w Uppercase) Run() {
	log.INFO.Println("Uppercase is running.")

	// Catch any errors that could happen
	defer func() {
		if r := recover(); r != nil {
			log.ERROR.Println("Uncaught error occurred, notifying leader and exiting.")

			w.Stop()
		}
	}()

	for {
		select {
		case cmd := <-w.Ports()["MGMT_Uppercase"].Channel():
			if cmd == "STOP" {
				w.Stop()
				return
			} else if cmd == "STATUS" {
				w.Ports()["MGMT_Uppercase"].Channel() <- "Healthy"
			} else if cmd == "METRICS" {
				w.Ports()["MGMT_Uppercase"].Channel() <- Metric{"health": "TODO metrics."}
			}
		case data := <-w.Ports()["Source_and_Uppercase"].Channel():
			w.Ports()["Sink_and_Uppercase"].Channel() <- strings.ToUpper(data.(string))
		}
	}
}

func (w Uppercase) Stop() {
	w.Ports()["MGMT_Uppercase"].Close()
	w.Ports()["Source_and_Uppercase"].Close()
	w.Ports()["Sink_and_Uppercase"].Close()

	log.INFO.Println("Worker " + w.Name() + " stopped.")
}
