// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"encoding/json"

	"github.com/stalehard/chat/message"
)

type SingleChannelMsg struct {
	c *Client
	data []byte
}

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]string

	// Inbound messages from the clients.
	broadcast chan []byte

	singlecast chan SingleChannelMsg

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		singlecast:  make(chan SingleChannelMsg),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]string),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = "Anonymous"

			var users []string
			for c := range h.clients {
				users = append(users, h.clients[c])
			}
			usersList := message.ServerUserList{users, 1}

			if str, err := json.Marshal(usersList); err != nil {
				log.Println("error: v%", err)
			} else {
				for c := range h.clients {
					select {
					case c.send <- str:
					default:
						close(c.send)
						delete(h.clients, c)
					}
				}
			}


		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				var users []string
				for c := range h.clients {
					users = append(users, h.clients[c])
				}
				usersList := message.ServerUserList{users, 1}

				if str, err := json.Marshal(usersList); err != nil {
					log.Println("error: v%", err)
				} else {
					for c := range h.clients {
						select {
						case c.send <- str:
						default:
							close(c.send)
							delete(h.clients, c)
						}
					}
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case singleMessage := <-h.singlecast:
			select {
			case singleMessage.c.send <- singleMessage.data:
			default:
				close(singleMessage.c.send)
				delete(h.clients, singleMessage.c)
			}

		}
	}
}
