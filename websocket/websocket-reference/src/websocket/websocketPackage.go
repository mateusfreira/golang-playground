package websocket

import(
	"log"
)
type Room struct {
	name string
	clients   map[int]*Client
}

func NewRoom(name string) *Room {
	clients := make(map[int]*Client)
	return &Room {
	name,
	clients};
}

type WebsocketPackage struct {
	clients   map[int]*Client
	rooms   map[string]*Room
}


func NewWebsocketPackage(name string) *WebsocketPackage {
	clients := make(map[int]*Client)
	rooms := make(map[string] *Room)
	return &WebsocketPackage {
	clients,
	rooms};
}

func (self *WebsocketPackage) emit(event string, body *Message, client *Client) {
	log.Printf("emit clients : %s", event)
	if _, has := self.rooms[event]; !has {
		log.Printf("event emit, in a inexitent room!")
	}else{
		for _, c := range self.rooms[event].clients {
			c.Write(body)
		}
	}
	/*
	//, len(self.rooms[event].clients)
	*/
}


func (self *WebsocketPackage) addClient(client *Client) {
	self.clients[client.id] = client
}

func (self *WebsocketPackage) on(room string, client *Client) {
	if _, has := self.rooms[room]; !has {
		self.rooms[room] = NewRoom(room)
	}
	self.rooms[room].clients[client.id] = client
}
/*
func (self *WebsocketPackage) emit(event string, msg *Message, client *Client) {
for _, c := range self.clientsByEvent {
c.Write(msg)
}
}
*/
