package workers

import (
	"emd/log"
	"emd/worker"
)

type Source struct {
	worker.Work
}

func (w Source) Init() {
	for _, p := range w.Ports() {
		p.Open()
	}

	log.INFO.Println("Worker " + w.Name() + " inited.")
}

func (w Source) Run() {
	log.INFO.Println("Source is running.")

	// Catch any errors that could happen
	defer func() {
		if r := recover(); r != nil {
			log.ERROR.Println("Uncaught error occurred, notifying leader and exiting.")

			w.Stop()
		}
	}()

	w.Ports()["Source_and_Uppercase"].Channel() <- "uppercase this"

	for {
		select {
		case cmd := <-w.Ports()["MGMT_Source"].Channel():
			if cmd == "STOP" {
				w.Stop()
				return
			} else if cmd == "STATUS" {
				w.Ports()["MGMT_Source"].Channel() <- "Healthy"
			} else if cmd == "METRICS" {
				w.Ports()["MGMT_Source"].Channel() <- Metric{"metric": "TODO metrics."}
			}
		case data := <-w.Ports()["Source_and_Uppercase"].Channel():
			log.INFO.Println(data)
		}
	}
}

func (w Source) Stop() {
	w.Ports()["MGMT_Source"].Close()
	w.Ports()["Source_and_Uppercase"].Close()

	log.INFO.Println("Worker " + w.Name() + " stopped.")
}
