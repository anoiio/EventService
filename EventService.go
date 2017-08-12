package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/anoiio/eventservice/worker"
)

func createEvent(body []byte) (worker.LogEvent, error) {
	var reqMap map[string]interface{}

	json.Unmarshal(body, &reqMap)
	fmt.Printf("Got request with body: %s\n", reqMap)

	eventType, err := worker.ValidateEventType(reqMap[worker.TypeKey].(string))

	if err != nil {
		return worker.LogEvent{}, err
	}

	event := worker.LogEvent{Type: eventType, Payload: reqMap[worker.PayloadKey].(map[string]interface{})}
	return event, nil
}

func requestHandler(w http.ResponseWriter, r *http.Request, eventQueue chan worker.LogEvent) {

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := createEvent(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	eventQueue <- event

	// Render success.
	w.WriteHeader(http.StatusCreated)
}

func main() {
	var (
		maxWorkers   = flag.Int("max_workers", 20, "The number of workers to start")
		maxQueueSize = flag.Int("max_queue_size", 500, "The size of event queue")
		port         = flag.String("port", "8080", "The server port")
	)
	flag.Parse()

	eventQueue := make(chan worker.LogEvent, *maxQueueSize)

	// Start the dispatcher.
	dispatcher := worker.NewDispatcher(eventQueue, *maxWorkers)
	dispatcher.Run()

	// Start the HTTP handler.
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r, eventQueue)
	})
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
