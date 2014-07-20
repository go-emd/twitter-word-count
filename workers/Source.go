package workers

import (
	"emd/log"
	"emd/worker"

	"github.com/darkhelmet/twitterstream"
)

type Source struct {
	worker.Work
}

var client *twitterstream.Client

func (w Source) Init() {
	for _, p := range w.Ports() {
		p.Open()
	}

	client = twitterstream.NewClient("9tXWK5ZjaQwjp4AoFUbsReNdY",
									 "qyNDhvcnqbbja3pBNcbQhOQzl0aXJdlch9WxTnwgR8QmPMuIbe",
									 "215144158-zzVvdHodezegE7orVZiYORaCIoJAA5NMFdQP9EOt",
									 "9ewMfIWuOjC2ODbJvlPYlzcfXZabjZxhIlOvl7QaSQAtC")

	log.INFO.Println("Worker " + w.Name() + " inited.")
}

func (w Source) Run() {
	log.INFO.Println("Source is running.")

	// Catch any errors that could happen
	defer func() {
		if r := recover(); r != nil {
			log.ERROR.Println("Uncaught error occurred, "+w.Name_+" is stopping.")

			w.Stop()
		}
	}()

	conn, err := client.Track("Stockmarket, Stocks")
	if err != nil {
		log.ERROR.Println(err)
		w.Stop()
	}

	// Main processing
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
		default:
			if tweet, err := conn.Next(); err == nil {
				w.Ports()["Source_and_Count"].Channel() <- tweet.Text
			} else {
				log.ERROR.Println("Unable to get next twwet.")
			}
		}
	}
}

func (w Source) Stop() {
	w.Ports()["MGMT_Source"].Close()
	w.Ports()["Source_and_Count"].Close()

	log.INFO.Println("Worker " + w.Name() + " stopped.")
}
