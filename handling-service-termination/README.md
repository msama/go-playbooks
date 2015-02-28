#Handling service termination

When developers design component of a SOA they assume that they will be running all the time and they normally overlook they way those components terminates. I hope this won't come as a shock to anybody but in a SOA microservices gets terminated all the times and making sure that they gracefully terminates is important. In this post I will present some common patterns that can be adopted to make service termination more graceful. The full solution is published in github.

The most common reasons why services terminates are:

* Autoscaling events: a common architecture for scaling microservices is to have horizontal replication with a traffic multiplexer and a scaling policy. In AWS that is normally achieved with an Elastic Loadbalancer and with an Autoscaling Group. Whenever the environment decides that the traffic is low and it is time to downscale one or more instances of a microservice are terminated.
* Continuous development: one of the paradigm of agile server side development is develop-commit-deploy. However every time a the service are redeployed the previous running version is terminated and replaced with a new one.
Un-handled errors: every time the process running the service abnormally terminates then also the service itself terminates. I have seen environments with services crushing every minute and upstart reloading them.
* Leaks: that is unfortunate but services with known memory or resource leaks are often terminated while a patch is under development.
* Canary testing: the easiest way to test new version of a micro service is to deploy instances of the new one within the production cluster. This process imply starting and then terminating a canary instance.

All the above should be enough to convince you that also in healthy environments microservices terminates all the time and being aware of how they terminate is important. Imagine a payment system abnormally terminating without waiting for all the transactions in progress to complete. This playbook shows examples of how to handle graceful termination in GO.


##Gracefully terminating

The first step for gracefully terminating is to handle errors from the `ListenAndServe` routine as well as signals.

```
func main() {

	// The service state
	serviceState := NewServerState()

	// Handle termination signals
	errors := make(chan error, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)

	// Initialise the handlers
	// ...

	// Start serving
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			errors <- err
		}
	}()
  
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
	serviceState.State = serviceStateTerminating
	i := 3
	for i > 0 {
		log.Printf("Terminating in %d\n", i)
		time.Sleep(1 * time.Second)
		i = i - 1
	}
}
```

Once errors and signals have been captured they can be properly handled. This playbook suggests to turn off the service gracefully by stopping serving new requests and by waiting for the ones in progress to terminate. Once again this can be done by wrapping the handler
with a state controller delegate that enforces some sort of service lifecycle. The state object in the playbook is only an example and can be improved.

Just a quick note about synchronization. In the playbook the `running` variable is not synchronised and in the common case there is no reason for it to be. Synchronising it would make the implementation slower and therefore it should only be done if there is a real reason. In the playbook example it doesn't really matter when requests stops being served as long as that happens quickly enough. In other systems it might be necessary to stop serving requests immediately as soon as other subcomponents are terminated and therefore the 
lifecycle need to be synchronised.

##Releasing resources upon termination

Services normally hold up resources such as connections to databases, caches or other services. Although those connections will eventually time out and being reassigned it is often a good practice to release them.

```
// Release all the resources and complete service termination
func ReleaseResourcesAndTerminate(state *ServiceState) {
	// Nothing else to do in this example
	state.State = serviceStateTerminated
}

func main() {
	// The service state
	serviceState := NewServerState()

	defer ReleaseResourcesAndTerminate(serviceState)
}
```

##Knowing uptime and termination time

It is often interesting to measure how long the service has been running. Also while examinging the logs of a malfunctioning service it is often usefult to know if the main process is still running or if it has terminated and when. This can be done by deferring a function before invoking `ListenAndServe`.

```
// Used to print the uptime at the end of service execution
func PrintUptime(start time.Time) {
	now := time.Now()
	log.Printf("Service was running for %v\n", now.Sub(start))
	log.Printf("Service was terminated at %v\n", now)
}

func main() {
	// At the end of the execution prints how long the service was running.
	defer PrintUptime(time.Now())
}

```

##Gracefully recovering from panic

GO automatically isolates panic inside each handler so that when a panic occurs it only affects a single request. However the default configuration simply terminates the request without returning any response to the caller. However it is possible to handle the panic by still responding to client. Notice that this is only possible unless the response writer has already been generated. In the playbook for instance the response is partially generated by the original handler and half by the panic handler.

A way to handle panics is to wrap the real handler with a wapper handler which will defer and recover in case of panic.

```
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
```