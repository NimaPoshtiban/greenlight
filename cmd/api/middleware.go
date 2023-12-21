package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

// recoverPanic recovers from panics in the HTTP handler chain.
//
// It takes an http.Handler as input and returns an http.Handler. When a panic occurs
// during the execution of the handler, this middleware recovers from the panic,
// sets the "Connection" header to "close", and calls the serverErrorResponse function
// to handle the error.
//
// Parameters:
// - next: The next http.Handler in the chain.
//
// Return type:
// - http.Handler: The recovered HTTP handler.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// rateLimit implements rate limiting for HTTP requests.
//
// It takes a http.Handler `next` as input and returns a http.Handler.
// The rateLimit function uses a map to keep track of clients and their
// respective rate limiters. It also uses a goroutine to periodically
// remove clients that haven't been seen in the last 3 minutes.
// The rateLimit function checks if rate limiting is enabled in the
// application's configuration. If it is, it extracts the client's IP
// address from the request and checks if the client already exists in
// the map. If the client doesn't exist, a new rate limiter is created
// and added to the map. The function then updates the last seen time
// for the client and checks if the rate limiter allows the request.
// If the request is not allowed, a rate limit exceeded response is
// sent. Finally, the function calls the `next` http.Handler to handle
// the request.
func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}
