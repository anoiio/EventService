# EventService

The purpose of this application is a case study of implementing event logging service in go. <br />

Design and code based on this post: [Handling 1 Million Requests per Minute with Go](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/) <br />

EventService listens for HTTP POST requests on localhost:8080/log

### Setup

#### Application parameters

1. max_workers (default 20) - defines a number of concurrent workers for events processing.
2. max_queue_size (default 500) - maximum number of events that can be buffered in the queue before processing.
3. port (default 8080) - server port.

#### DB configuration

EventService stores events to PostgreSQL database. <br />
Following parameters, cab be changed in Worker.go file

1. host (default "localhost")
2. port (default 5432)
3. user (default "postgres")
4. password (default "postgres")
5. dbname (default "anoiio")


Required table can be created using /db/init.sql script

### Run

To start the application: $GOPATH/bin/eventservice  <br />

EventService supports 3 events types. <br />
They can be submitted to the service by issuing HTTP POST request with the following body:

1. Click

```javascript
{
  "type": "Click",
  "payload": {
    "date_time": 1502542550914,
    "transaction_id": 11111,
    "ad_type": "video",
    "time_to_click": 5000,
    "user_id": "fydsgf-t73r8nf"
  }
}
```


2. Impression

```javascript
{
  "type": "Impression",
  "payload": {
    "date_time": 1502542550914,
    "transaction_id": 22222,
    "ad_type": "video",
    "user_id": "238rjhwefn-vmedo"
  }
}
```


3. Completion

```javascript
{
  "type": "Completion",
  "payload": {
    "date_time": 1502542550914,
    "transaction_id": 33333,
    "ad_type": "video",
    "user_id": "7fh437fh3ug-rnbjr"
  }
}
```
