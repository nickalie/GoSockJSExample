package ng

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type Room struct {
	id      int
	max     int
	clients []sockjs.Session
}

func NewRoom(id int, max int) *Room {
	return &Room{id: id, max: max}
}

func (this *Room) AddClient(client sockjs.Session) {
	if this.Len() >= this.max {
		return
	} else if indexOf(this.clients, client) == -1 {
		message := "{\"type\":\"userConnected\"}"
		this.sendToAll(message)

		for i := 0; i < this.Len(); i++ {
			client.Send(message)
		}

		this.clients = append(this.clients, client)
	}
}

func (this *Room) RemoveClient(client sockjs.Session) {
	index := indexOf(this.clients, client)
	if index != -1 {
		this.clients = append(this.clients[:index], this.clients[index+1:]...)
		this.sendToAll("{\"type\":\"userDisconnected\"}")
	}
}

func (this *Room) Send(message string, sender sockjs.Session) {
	for _, t := range this.clients {
		if t != sender {
			t.Send(message)
		}
	}
}

func (this *Room) sendToAll(message string) {
	for _, t := range this.clients {
		t.Send(message)
	}
}

func (this *Room) ID() int {
	return this.id
}

func (this *Room) Len() int {
	return len(this.clients)
}

func indexOf(slice []sockjs.Session, value sockjs.Session) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
