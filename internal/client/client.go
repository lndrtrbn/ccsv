package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Client struct {
	// Name of the client user.
	name string

	// Chat server address.
	serverAddr string
	// Websocket connection with server.
	connection *websocket.Conn
}

func (client *Client) Start(onMessageFunc func(msg Message)) {
	go client.readServer(onMessageFunc)
}

func (client *Client) Send(content string) {
	message := Message{
		Name:    client.name,
		Content: content,
	}
	str, err := json.Marshal(message)
	if err != nil {
		log.Printf("| Error while marshaling the message %s, %v", message, err)
		return
	}

	req, err := http.NewRequest(
		"POST",
		client.serverAddr+"/publish",
		bytes.NewBuffer([]byte(str)),
	)
	if err != nil {
		log.Printf("| Error while sending the message %s, %v", message, err)
		return
	}

	req.Header.Add("X-ACCESS-TOKEN", os.Getenv("ACCESS_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("| Error while sending the message %s, %v", message, err)
		return
	}
	defer resp.Body.Close()
}

func NewClient(address string, name string) *Client {
	headers := http.Header{}
	headers.Add("X-ACCESS-TOKEN", os.Getenv("ACCESS_TOKEN"))

	connection, _, err := websocket.Dial(
		context.Background(),
		address+"/subscribe",
		&websocket.DialOptions{
			HTTPHeader: headers,
		})
	if err != nil {
		log.Fatalf("| Error while trying to connect the client %v", err)
	}

	return &Client{
		name:       name,
		serverAddr: address,
		connection: connection,
	}
}

func (client *Client) readServer(onMessageFunc func(msg Message)) {
	for {
		var msg Message
		err := wsjson.Read(context.Background(), client.connection, &msg)
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Printf("| Error while reading message %v", err)
			continue
		}

		onMessageFunc(msg)
	}
}
