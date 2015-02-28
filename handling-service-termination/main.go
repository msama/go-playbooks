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

// Factory method for helloWorldHandler
func NewHelloWorldHandler() http.Handler {
	return &helloWorldHandler{}
}

// An http.Handler which says hello world
type helloWorldHandler struct{}

func (h *helloWorldHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Simulate Delay
	time.Sleep(1 * time.Second)
	io.WriteString(rw, "hello, world!\n")
}

// Factory method for panicNowHandler
func NewPanicNowHandler() http.Handler {
	return &panicNowHandler{}
}

// An http.Handler which simulates panic
type panicNowHandler struct{}

func (h *panicNowHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, "Starting to respond before the failure. \n")
	panic("simulated panic!")
}

// Factory method for panicHandler
func NewPanicHandler(delegate http.Handler) http.Handler {
	return &panicHandler{
		delegate: delegate,
	}
}

// An http.Handler which can recover from panic events
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

type serviceStateValue int

const (
	serviceStateBootstrapping serviceStateValue = iota
	serviceStateRunning
	serviceStateTerminating
	serviceStateTerminated
)

func NewServerState() *ServiceState {
	return &ServiceState{
		State: serviceStateBootstrapping,
	}
}

type ServiceState struct {
	State serviceStateValue
}

// Factory method for terminationHandler
func NewTerminationHandler(state *ServiceState, delegate http.Handler) http.Handler {
	return &terminationHandler{
		state:    state,
		delegate: delegate,
	}
}

// Stop serving new requests if the server is terminating
type terminationHandler struct {
	state    *ServiceState
	delegate http.Handler
}

func (h *terminationHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch h.state.State {
	case serviceStateRunning:
		h.delegate.ServeHTTP(rw, req)
	case serviceStateBootstrapping:
		io.WriteString(rw, "Service is bootstrapping\n")
		rw.WriteHeader(http.StatusServiceUnavailable)
	case serviceStateTerminating:
		io.WriteString(rw, "Service is terminating\n")
		rw.WriteHeader(http.StatusServiceUnavailable)
	case serviceStateTerminated:
		rw.WriteHeader(http.StatusServiceUnavailable)
	}
}

// Used to print the uptime at the end of service execution
func PrintUptime(start time.Time) {
	now := time.Now()
	log.Printf("Service was running for %v\n", now.Sub(start))
	log.Printf("Service was terminated at %v\n", now)

}

// Release all the resources and completes service termination
func ReleaseResourcesAndTerminate(state *ServiceState) {
	// Nothing else to do in this example
	state.State = serviceStateTerminated
}

func main() {
	// The service state
	serviceState := NewServerState()

	// At the end of the execution prints how long the service was running.
	defer PrintUptime(time.Now())
	defer ReleaseResourcesAndTerminate(serviceState)

	// Handle termination signals
	errors := make(chan error, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)

	// Setup and run this example
	http.Handle("/hello", NewTerminationHandler(serviceState, NewPanicHandler(NewHelloWorldHandler())))
	http.Handle("/panic", NewTerminationHandler(serviceState, NewPanicHandler(NewPanicNowHandler())))
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			errors <- err
		}
	}()

	// Log that the service is now running on a certain port
	log.Printf("Welcome:\n" +
		"Connect to http://localhost:8080/hello for hello world\n" +
		"Connect to http://localhost:8080/panic for simulating a panic\n")

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

	// Gracefully terminates after few seconds hoping that
	// most of the handlers would have terminated.
	// This is just a simple example. A proper implementation would
	// depend on the nature of the service. In most of cases there is
	// no need wrapp the boolean with a Mutex unless you need to exact
	// control on when requests are no longer being served.
	serviceState.State = serviceStateTerminating
	i := 3
	for i > 0 {
		log.Printf("Terminating in %d\n", i)
		time.Sleep(1 * time.Second)
		i = i - 1
	}
}
