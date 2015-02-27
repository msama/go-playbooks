package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func HelloWorldHandler() http.Handler {
	return &helloWorldHandler{}
}

type helloWorldHandler struct{}

func (h *helloWorldHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Simulate Delay
	time.Sleep(1 * time.Second)
	io.WriteString(rw, "hello, world!\n")
}

func PanicNowHandler() http.Handler {
	return &panicNowHandler{}
}

type panicNowHandler struct{}

func (h *panicNowHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	panic("simulated panic!")
}

// Creates a new handler which handles panic
func PanicHandler(delegate http.Handler) http.Handler {
	return &panicHandler{
		delegate: delegate,
	}
}

// An handler which can recover from panic events
type panicHandler struct {
	delegate http.Handler
}

func (h *panicHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		panicData := recover()
		if panicData != nil {
			// Handle panics
			log.Printf("Just recovered from panic: %v\n", panicData)
			rw.WriteHeader(500)
			io.WriteString(rw, "An error occurred and was internally handled :)\n")
		}
	}()
	h.delegate.ServeHTTP(rw, req)
}

type ServiceState struct {
	Running bool
}

// Stop serving new requests if the server is terminating
func TerminationHandler(state *ServiceState, delegate http.Handler) http.Handler {
	return &terminationHandler{
		state:    state,
		delegate: delegate,
	}
}

type terminationHandler struct {
	state    *ServiceState
	delegate http.Handler
}

func (h *terminationHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !h.state.Running {
		io.WriteString(rw, "Service is terminating\n")
		rw.WriteHeader(503)
		return
	}
	h.delegate.ServeHTTP(rw, req)
}

// Used to print the uptime at the end of service execution
func PrintUptime(start time.Time) {
	log.Printf("Service was running for %v\n", time.Now().Sub(start))
}

func main() {
	// At the end of the execution prints how long the service was running.
	defer PrintUptime(time.Now())

	// The service state
	serviceState := &ServiceState{
		Running: true,
	}

	// Handle termination signals
	errors := make(chan error, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)

	// Setup and run this example
	http.Handle("/hello", TerminationHandler(serviceState, PanicHandler(HelloWorldHandler())))
	http.Handle("/panic", TerminationHandler(serviceState, PanicHandler(PanicNowHandler())))
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			errors <- err
		}
	}()

	log.Printf("Welcome:\n" +
		"Connect to http://localhost:8080/hello for hello world\n" +
		"Connect to http://localhost:8080/panic for simulating a panic\n")

	// Waits until the service fails or it is terminated.
	select {
	case err := <-errors:
		log.Printf("Error: %v\n", err)
		break
	case sig := <-signals:
		log.Printf("Signal: %v\n", sig)
		break
	}

	// Gracefully terminates after few seconds hoping that
	// most of the handlers would have terminated.
	// This is just a simple example. A proper implementation would
	// depend on the nature of the service. In most of cases there is
	// no need wrapp the boolean with a Mutex unless you need to exact
	// control on when requests are no longer being served.
	serviceState.Running = false
	i := 5
	for i > 0 {
		log.Printf("Terminating in %d\n", i)
		time.Sleep(1 * time.Second)
		i = i - 1
	}
}
