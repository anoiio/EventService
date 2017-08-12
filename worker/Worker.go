package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	// will justify later
	_ "github.com/lib/pq"
)

// EventType - a type for event type representation
type EventType string

// evalable event constants, db config
const (
	TypeKey    string = "type"
	PayloadKey string = "payload"

	Impression EventType = "Impression"
	Click      EventType = "Click"
	Completion EventType = "Completion"

	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "anoiio"
)

// ValidateEventType checks inputs string event type against constants and return EventType representation if exists
func ValidateEventType(eventTypeStr string) (EventType, error) {

	inType := EventType(eventTypeStr)

	switch inType {
	case Impression:
		return Impression, nil
	case Click:
		return Click, nil
	case Completion:
		return Completion, nil

	default:
		return EventType(""), errors.New("Not supported event type: " + eventTypeStr)
	}
}

// LogEvent - event to be logged
type LogEvent struct {
	Type    EventType
	Payload map[string]interface{}
}

// Worker - worker goroutine struct
type Worker struct {
	eventQueue chan LogEvent
	workerPool chan chan LogEvent
	quitChan   chan bool
	db         *sql.DB
}

// NewWorker creates Worker
func NewWorker(workerPool chan chan LogEvent) Worker {

	return Worker{
		eventQueue: make(chan LogEvent),
		workerPool: workerPool,
		quitChan:   make(chan bool),
		db:         nil,
	}
}

// Start - opens db connection, starts Worker processing in new goroutine
func (w Worker) Start() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	w.db = db
	go process(w)
}

func process(w Worker) {
	fmt.Printf("Worker started\n")
	defer w.db.Close()
	for {
		w.workerPool <- w.eventQueue

		select {
		case event := <-w.eventQueue:
			w.handleEvent(event)
		case <-w.quitChan:
			fmt.Printf("Worker stopping\n")
			return
		}
	}
}

// Stop sends 'stop' signal to Worker
func (w Worker) Stop() {
	go func() {
		w.quitChan <- true
	}()
}

func (w Worker) handleEvent(event LogEvent) error {
	switch event.Type {
	case Impression:
		w.handleImpression(event.Payload)
	case Click:
		w.handleClick(event.Payload)
	case Completion:
		w.handleCompletion(event.Payload)
	default:
		return errors.New("Not supported event type")
	}
	return nil
}

func (w Worker) handleClick(data map[string]interface{}) {
	dateTime := int64(data["date_time"].(float64))
	transactionID := data["transaction_id"]
	adType := data["ad_type"]
	timeToClick := data["time_to_click"]
	userID := data["user_id"]

	_, err := w.db.Exec(fmt.Sprintf("INSERT INTO EVENTS (type, date_time, transaction_id, ad_type, time_to_click, user_id) VALUES ($1, $2, $3, $4, $5, $6)"),
		Click, time.Unix(0, dateTime*int64(time.Millisecond)), transactionID, adType, timeToClick, userID)

	if err != nil {
		fmt.Printf("DB error: %s\n", err)
	}
}

func (w Worker) handleCompletion(data map[string]interface{}) {
	dateTime := int64(data["date_time"].(float64))
	transactionID := data["transaction_id"]

	_, err := w.db.Exec(fmt.Sprintf("INSERT INTO EVENTS (type, date_time, transaction_id) VALUES ($1, $2, $3)"), Completion, time.Unix(0, dateTime*int64(time.Millisecond)), transactionID)
	if err != nil {
		fmt.Printf("DB error: %s\n", err)
	}
}

func (w Worker) handleImpression(data map[string]interface{}) {
	dateTime := int64(data["date_time"].(float64))
	transactionID := data["transaction_id"]
	adType := data["ad_type"]
	userID := data["user_id"]
	_, err := w.db.Exec(fmt.Sprintf("INSERT INTO EVENTS (type, date_time, transaction_id, ad_type, user_id) VALUES ($1, $2, $3, $4, $5)"), Impression, time.Unix(0, dateTime*int64(time.Millisecond)), transactionID, adType, userID)
	if err != nil {
		fmt.Printf("DB error: %s\n", err)
	}
}
