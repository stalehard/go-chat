package message

import (
	"time"
	"encoding/json"
	"errors"
)

type ClientMessage struct {
	Name string
	Text string
}

type ServerMessage struct {
	Type int
	Text string
	Name string
	Time time.Time
	Users []string
}

type ServerUserList struct {
	Users []string
	Type int
}

func NewMessage(data []byte) (ServerMessage, error) {
	var parsed ClientMessage
	var users [0]string

	if err := json.Unmarshal(data, &parsed); err != nil {
		return ServerMessage{10, "Error parse", "", time.Now(), users[:]}, err
	} else {
		if len(parsed.Name) == 0 {
			return ServerMessage{10, "Empty name", "", time.Now(), users[:]}, errors.New("Invalid message")
		}
		return ServerMessage{0, parsed.Text, parsed.Name, time.Now(), users[:]}, nil
	}
}