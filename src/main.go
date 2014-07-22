package main

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"net/http"
	"strconv"
	"strings"
	"flag"
	"./ng"
)
//room id -> Room
var rooms map[int]*ng.Room

//client -> Room
var clientToRooms map[sockjs.Session]*ng.Room

func main() {
	rooms = make(map[int]*ng.Room)
	clientToRooms = make(map[sockjs.Session]*ng.Room)
	handler := sockjs.NewHandler("/sockjs", sockjs.DefaultOptions, sockjsHandler)
	log.Fatal(http.ListenAndServe("127.0.0.1:" + *flag.String("port", "8081", "Port"), handler))
}

func sockjsHandler(client sockjs.Session) {
	//new client connected
	for {
		if msg, err := client.Recv(); err == nil {
			//message received
			if strings.HasPrefix(msg, "join") {

				room, ok := clientToRooms[client]

				if ok {
					removeClientFromRoom(room, client)
				}

				roomId, err := strconv.Atoi(strings.TrimPrefix(msg, "join"))

				if err != nil {
					continue
				}

				room, ok = rooms[roomId]

				//create new room if not exists
				if !ok {
					room = ng.NewRoom(roomId, 2)
					rooms[roomId] = room
				}

				clientToRooms[client] = room
				room.AddClient(client)

			} else if room, ok := clientToRooms[client]; ok {
				room.Send(msg, client)
			}
		} else {
			break
		}
	}

	if room, ok := clientToRooms[client]; ok {
		removeClientFromRoom(room, client)
	}
}

func removeClientFromRoom(room *ng.Room, client sockjs.Session) {
	room.RemoveClient(client)
	delete(clientToRooms, client)
	if room.Len() == 0 {
		delete(rooms, room.ID())
	}
}
