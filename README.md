# EventService

Purpose of this application is a case stady of implementing event logging service in go. <br />

EventService listens for POST http requests on localhost:8080/log

### Setup

#### Application parameters

1. max_workers (defuault 20) - defines number of concurrent workers for events processing.
2. max_queue_size (defuault 500) - muximum number of events that can be buffered in queue before processing.
3. port (defuault 8080) - server port.

#### DB configuration

EventService stores events to PostgeSQL database. <br />
Following parameters cab be changed in Worker.go file

1. host (defuault "localhost")
2. port (defuault 5432)
3. user (defuault "postgres")
4. password = (defuault "postgres")
5. dbname   = (defuault "anoiio")


Requared table can be created using /db/init.sql script

### Run

EventService supports 3 events types.
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
