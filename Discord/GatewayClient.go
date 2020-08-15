package Discord

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"main/Discord/Entities/Gateway"
	"time"
)

var connectionLostError = errors.New("Connection lost")

type GatewayClient struct {
	GatewayUrl string

	connection *websocket.Conn
	heartbeat chan struct{}

	heartbeatConfirmed bool
	lastSequenceId *int
}

// Setup

func (ws *GatewayClient) listenToSocket() {
	if ws.connection == nil { return }

	for {
		msg := new(Gateway.Message)
		err := ws.connection.ReadJSON(msg)

		if err != nil {
			log.Println("Discord gateway error: ", err)
			return
		}

		ws.onMessageReceived(msg)
	}
}

func (ws *GatewayClient) startHearthbeat(interval float64) {
	log.Println("Heartbeat in ", interval, " ms")

	heartbeat := time.NewTicker(time.Duration(interval) * time.Millisecond)
	ws.heartbeat = make(chan struct{})
	ws.heartbeatConfirmed = true

	go func() {
		for {
			select {
				case <- heartbeat.C:
					err := ws.sendHearthbeat()

					if err == connectionLostError {
						log.Println("Connection lost. Reconnecting...")
						ws.Disconnect()
						ws.Connect()
					} else if err != nil {
						heartbeat.Stop()
						return
					}
				case <- ws.heartbeat:
					heartbeat.Stop()
					return
			}
		}
	}()
}

func (ws *GatewayClient) sendHearthbeat() error {
	if !ws.heartbeatConfirmed {
		return connectionLostError
	}

	ws.heartbeatConfirmed = false

	message := Gateway.Message{Op: 1, D: ws.lastSequenceId}
	return ws.SendMessage(message)
}

func (ws *GatewayClient) stopHeartbeat() {
	close(ws.heartbeat)
}

// Actions

func (ws *GatewayClient) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(ws.GatewayUrl, nil)

	if err != nil {
		return err
	}

	log.Println("Connected")

	ws.connection = conn
	ws.listenToSocket()

	return nil
}

func (ws *GatewayClient) Disconnect() {
	ws.stopHeartbeat()
	ws.connection.Close()

	ws.heartbeat = nil
	ws.connection = nil

	log.Println("Disconnected")
}

func (ws *GatewayClient) SendMessage(message Gateway.Message) error {
	data, _ := json.Marshal(message)
	log.Println("Sending message: ", string(data))

	return ws.connection.WriteJSON(message)
}

// Message receive

func (ws *GatewayClient) onMessageReceived(message *Gateway.Message) {
	data, _ := json.Marshal(message)
	log.Println("Received message: ", string(data))

	ws.lastSequenceId = message.S

	switch message.Op {
	case 1:
		msg := Gateway.Message{Op: 11}
		ws.SendMessage(msg)
	case 10:
		params := message.D.(map[string]interface{})
		ws.startHearthbeat(params["heartbeat_interval"].(float64))
	case 11:
		ws.heartbeatConfirmed = true
	}
}
