package workers

import (
	"github.com/go-emd/emd/log"
	"github.com/go-emd/emd/worker"
	"strings"
)

type Count struct {
	worker.Work
}

var counters map[string]int
func highestCount() string{
	word := ""
	highest := 0

	for w, c := range counters {
		if c > highest {
			highest = c
			word = w
		}
	}

	return word
}

func (w Count) Init() {
	for _, p := range w.Ports() {
		p.Open()
	}

	counters = make(map[string]int)

	log.INFO.Println("Worker " + w.Name() + " inited.")
}

func (w Count) Run() {
	log.INFO.Println("Count is running.")

	// Catch any errors that could happen
	defer func() {
		if r := recover(); r != nil {
			log.ERROR.Println("Uncaught error occurred, notifying leader and exiting.")

			w.Stop()
		}
	}()

	for {
		select {
		case cmd := <-w.Ports()["MGMT_Count"].Channel():
			if cmd == "STOP" {
				w.Stop()
				return
			} else if cmd == "STATUS" {
				w.Ports()["MGMT_Count"].Channel() <- "Healthy"
			} else if cmd == "METRICS" {
				w.Ports()["MGMT_Count"].Channel() <- Metric{"health": "TODO metrics."}
			}
		case tweet := <-w.Ports()["Source_and_Count"].Channel():
			words := strings.Split(tweet.(string), " ")
			for _, word := range words {
				counters[word] += 1
			}

			w.Ports()["Sink_and_Count"].Channel() <- highestCount()
		}
	}
}

func (w Count) Stop() {
	w.Ports()["MGMT_Count"].Close()
	w.Ports()["Source_and_Count"].Close()
	w.Ports()["Sink_and_Count"].Close()

	log.INFO.Println("Worker " + w.Name() + " stopped.")
}
