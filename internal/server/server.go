package server

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

	"nhooyr.io/websocket"
)

// Chat server to relay messages between clients.
type Server struct {
	// Routes endpoints to handlers.
	serveMux http.ServeMux

	// List of all subscribers.
	subscribers map[*websocket.Conn]struct{}
	// Mutex to manipulate the list of subscribers.
	subscribersMu sync.Mutex

	// API endpoints
	subscribeEndpoint   string
	publishEndpoint     string
	healthcheckEndpoint string
}

// Function implemented to be able to use our Server struct as a http.Handler.
func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.serveMux.ServeHTTP(writer, request)
}

// To start our http server that listens on the given address.
func (server *Server) ListenAndServe(address string) {
	http.Handle("/", server)
	log.Printf("| HTTP server listening on %s\n", address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("| Error trying to start http server on: %s", address)
		return
	}
}

// Creates a new Server and initializes the endpoints.
//
// We have two endpoints:
// - one for clients to receive messages (subscribe),
// - one for clients to send messages (publish).
func NewServer() *Server {
	server := &Server{
		subscribers:         make(map[*websocket.Conn]struct{}),
		subscribeEndpoint:   "/subscribe",
		publishEndpoint:     "/publish",
		healthcheckEndpoint: "/healthcheck",
	}
	server.serveMux.HandleFunc(server.subscribeEndpoint, server.subscribeHandler)
	server.serveMux.HandleFunc(server.publishEndpoint, server.publishHandler)
	server.serveMux.HandleFunc(server.healthcheckEndpoint, server.healthcheckHandlerwriter)
	return server
}

// Registers a subscriber.
func (server *Server) addSubscriber(connection *websocket.Conn) {
	server.subscribersMu.Lock()
	server.subscribers[connection] = struct{}{}
	server.subscribersMu.Unlock()
}

// Deletes the given subscriber.
func (server *Server) deleteSubscriber(connection *websocket.Conn) {
	server.subscribersMu.Lock()
	delete(server.subscribers, connection)
	server.subscribersMu.Unlock()
}

func (server *Server) subscribeHandler(writer http.ResponseWriter, request *http.Request) {
	connection, err := websocket.Accept(writer, request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("| /subscribe: Error trying to create the ws connection: %v", err)
		return
	}
	defer connection.Close(websocket.StatusNormalClosure, "")

	ctx := connection.CloseRead(request.Context())

	server.addSubscriber(connection)
	defer server.deleteSubscriber(connection)

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			if errors.Is(err, context.Canceled) {
				return
			}
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				return
			}
			if err != nil {
				log.Printf("| /subscribe: Error with context: %v", err)
				return
			}
		}
	}
}

// Handler when a request on /publish endpoint is received.
//
// Parse the received message and send it to every subscribers.
func (server *Server) publishHandler(writer http.ResponseWriter, request *http.Request) {
	// Check we well received a POST http request.
	if request.Method != "POST" {
		log.Printf("| /publish: Error because using method '%s'\n", request.Method)
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Get and parse the request content.
	body := http.MaxBytesReader(writer, request.Body, 2048)
	msg, err := io.ReadAll(body)
	if err != nil {
		log.Printf("| /publish: Error because request is too large")
		http.Error(writer, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}
	log.Printf("| /publish: %s", string(msg))

	// Send the message to everyone.
	server.subscribersMu.Lock()
	defer server.subscribersMu.Unlock()
	for sub := range server.subscribers {
		sub.Write(request.Context(), websocket.MessageText, msg)
	}

	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusAccepted)
}

// Handler when a request on /healthcheck endpoint is received.
//
// Parse the received message and send it to every subscribers.
func (server *Server) healthcheckHandlerwriter(writer http.ResponseWriter, request *http.Request) {
	// Check we well received a GET http request.
	if request.Method != "GET" {
		log.Printf("| /healthcheck: Error because using method '%s'\n", request.Method)
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("true"))
}
