# Core project code is here:
This is the directory where actual api data processing happens.

## Directory structure:
```
pkg
│   ├── dal
│   │   ├── dalintf.go
│   │   ├── entities.go
│   │   └── mysql.go
│   ├── handlers
│   │   └── app.go
│   ├── README.md
│   └── util
│       ├── httpErrorUtil.go
│       └── utils.go

```
## dal (i.e. data access layer):
This directory contains `Message` data model (`entities.go`), interface to the database driver (`dalintf.go`) and database communication code(`mysql.go`).

`entities.go` defines the Message data model:
```
package dal

//Message ... message object
type Message struct {
	ID      int    `json:"id,omitempty" example:"1"`
	Message string `json:"message,omitempty" example:"it is what is it"`
}

```

`dalintf.go` contains CRUD methods for Message
```
package dal

// DalInterface ... is an interface to data access layer methods
type DalInterface interface {
	AddMessage(msg *Message) error
	GetMessage(msg *Message) error
	UpdateMessage(msg *Message) error
	DeleteMessage(msg *Message) error
	GetAll() ([]Message, error)
}

```
## handler:
This is the REST API requests handler for rest-go-microservice. The `app.go` contains the application structure, swagger specs and handler methods for the REST endpoints.

Currently, following API endpoints are available for rest-go-microservice: 
```
1) A rest-go-microservice hello message GET api at /hello

2) Fetch a list of messages in response to a valid GET request at /messages, 

3) Create a new message in response to a valid POST request at /messages,

4) Check if a given message is a palindrome in response to a valid GET request at /messages/{id},

5) Fetch a message in response to a valid GET request at /messages/{id}

6) Update a message in response to a valid PUT request at /messages/{id},

7) Delete a message in response to a valid DELETE request at /messages/{id},

The {id} will determine which message the request will work with.
```
The `WithStats(h http.HandlerFunc)` function in `app.go` handlers with stats reporting. It tracks metrics such as request recieve time, latency.
```
func WithStats(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		tags := util.GetStatsTags(r)
		util.RequestFrom(tags, start)

		h(w, r)

		duration := time.Since(start)
		util.RecordLatency(tags, duration)
	}
}

```
## util:
This directory contains various util functions required by the rest-go-microservice handlers.



