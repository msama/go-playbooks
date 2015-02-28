package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// the sample data structure that we want to fetch
type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Id        string    `json:"id"`
}

type DataSource interface {
	QueryMessage(id string) (*Message, error)
}

func NewFakeDataSource() DataSource {
	return &fakeDataSource{}
}

// A fake DB data source.
// This simulates a database or any other data source used to fetch
// or to generate the message.
type fakeDataSource struct{}

func (s *fakeDataSource) QueryMessage(id string) (*Message, error) {
	if id == "" {
		return nil, fmt.Errorf("Id cannot be empty")
	}
	ts := time.Now()
	return &Message{
		Timestamp: ts,
		Id:        id,
	}, nil
}

type RequestCache interface {
	GetValue(key string) []byte
	SetValue(key string, val []byte)
}

func NewFakeCache() RequestCache {
	return &fakeCache{
		cache: map[string][]byte{},
	}
}

// A fake Cache.
// This simulates the component storing data in the meme cache.
type fakeCache struct {
	cache map[string][]byte
}

func (s *fakeCache) GetValue(key string) []byte {
	val, ok := s.cache[key]
	if !ok {
		return nil
	}
	return val
}

func (s *fakeCache) SetValue(key string, val []byte) {
	s.cache[key] = val
}

// The function used to generate keys.
// The request is used to create a unique key that unique key is used
// to store cahced values.
type hashFunction func(req *http.Request) string

// A dummy implementation that uses an id param
func MessageIdHash(req *http.Request) string {
	values := req.URL.Query()
	return values.Get("id")
}

// Factory method for jsonCacheHandler
func NewJsonCacheHandler(hashFn hashFunction) http.Handler {
	return &jsonCacheHandler{
		hashFn: hashFn,
		db:     NewFakeDataSource(),
		cache:  NewFakeCache(),
	}
}

// An http.Handler which can recover from panic events
type jsonCacheHandler struct {
	hashFn hashFunction
	db     DataSource
	cache  RequestCache
}

func (h *jsonCacheHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	key := h.hashFn(req)
	var messageJson []byte

	// Attempt to read from the cache
	val := h.cache.GetValue(key)
	if val == nil {
		// If the cahce miss read from the db then update the cache
		log.Printf("Miss for key %s", key)
		message, err := h.db.QueryMessage(key)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error querying the DB: %s", err)
			return
		}
		val, err := json.Marshal(message)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error encoding: %s", err)
			return
		}
		messageJson = val
		h.cache.SetValue(key, val)
	} else {
		// If the cache hits retur its value
		log.Printf("Hit for key %s", key)
		messageJson = val
	}

	// Generate the response
	rw.WriteHeader(http.StatusOK)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(messageJson) // Ignore write/connection errors
}

func main() {
	// Handle termination signals
	errors := make(chan error, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)

	// Setup and run this example
	http.Handle("/", NewJsonCacheHandler(MessageIdHash))
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			errors <- err
		}
	}()

	// Log that the service is now running on a certain port
	log.Printf("Welcome:\n" +
		"Connect to http://localhost:8080/?id=<value> for getting a message\n")

	// Waits until the service fails or it is terminated.
	select {
	case err := <-errors:
		// Handles the error from http.ListenAndServe
		log.Printf("Error: %v\n", err)
		break
	case sig := <-signals:
		// Handles shotdown signals
		log.Printf("Signal: %v\n", sig)
		break
	}
}
