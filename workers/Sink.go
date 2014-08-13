package workers

import (
	"github.com/go-emd/emd/log"
	"github.com/go-emd/emd/worker"
)

type Sink struct {
	worker.Work
}

func (w Sink) Init() {
	for _, p := range w.Ports() {
		p.Open()
	}

	log.INFO.Println("Worker " + w.Name() + " inited.")
}

func (w Sink) Run() {
	log.INFO.Println("Sink is running.")

	// Catch any errors that could happen
	defer func() {
		if r := recover(); r != nil {
			log.ERROR.Println("Uncaught error occurred, notifying leader and exiting.")

			w.Stop()
		}
	}()

	for {
		select {
		case cmd := <-w.Ports()["MGMT_Sink"].Channel():
			if cmd == "STOP" {
				w.Stop()
				return
			} else if cmd == "STATUS" {
				w.Ports()["MGMT_Sink"].Channel() <- "Healthy"
			} else if cmd == "METRICS" {
				w.Ports()["MGMT_Sink"].Channel() <- Metric{"name": "value"}
			}
		case data := <-w.Ports()["Sink_and_Count"].Channel():
			log.INFO.Println(data)
		}
	}
}

func (w Sink) Stop() {
	w.Ports()["MGMT_Sink"].Close()
	w.Ports()["Sink_and_Count"].Close()

	log.INFO.Println("Worker " + w.Name() + " stopped.")
}
