package websocket

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type ClientMessage struct{
	message *Message
	client  *Client
}
// Chat server.
type Server struct {
	pattern   string
	messages  []*Message
	clients   map[int]*Client
	clientsByEmail   map[string]*Client
	packages   map[string]*WebsocketPackage
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *ClientMessage
	doneCh    chan bool
	errCh     chan error
}



// Create new chat server.
func NewServer(pattern string) *Server {
	messages := []*Message{}
	clients := make(map[int]*Client)
	clientsByEmail := make(map[string]*Client)
	packages := make(map[string]*WebsocketPackage)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *ClientMessage)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		messages,
		clients,
		clientsByEmail,
		packages,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *Message, cli *Client) {
	s.sendAllCh <-  &ClientMessage { msg, cli }
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.messages {
		c.Write(msg)
	}
}

func (s *Server) sendAll(msg *Message) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

func (s *Server) loginAsAClient(email string, client *Client) {
	log.Println("Email login...", email)
	s.clientsByEmail[email] = client
}

func (s *Server) createPackage(name string, client *Client) {
	log.Println("createPackage", name)
	s.packages[name] = NewWebsocketPackage(name)
}


func (s *Server) on(namespace string, event string, client *Client) {
	s.packages[namespace].on(event, client)
}
func (s *Server) emit(name string, event string, body *Message, client *Client) {
	log.Println("emit", name)
	if _, has := s.packages[name]; !has {
		log.Println("emit in inexitend package")
	}else{
		s.packages[name].emit(event, body,client)
	}

}





// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
		log.Println("Dispach event %v", msg.message);
			switch msg.message.Event {
					case "connectAsClient":
						s.loginAsAClient(msg.message.Body, msg.client)
					case "createRoom":
							//
					case "createPackage":
						s.createPackage(msg.message.Body, msg.client)
					case "joinRoom":
						s.on(msg.message.Package, msg.message.Body, msg.client)
					case "on":
						s.on(msg.message.Package, msg.message.Body, msg.client)
					case "emit":
						s.emit(msg.message.Package, msg.message.Room, msg.message,  msg.client)
					default:
						log.Println(msg.message.Event+" Send all:", msg.message)
						s.emit(msg.message.Package, msg.message.Event, msg.message,  msg.client)
			}
		case err := <-s.errCh:
			log.Println("Error:", err.Error())
		case <-s.doneCh:
			return
		}
	}
}
