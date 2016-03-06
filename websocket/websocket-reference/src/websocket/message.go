package websocket

type Message struct {
	Event  string `json:"event"`
	Room  string `json:"room"`
	Package  string `json:"package"`
	Body   string `json:"body"`
}

func (self *Message) String() string {
	return "Package : "+self.Package + " Body: " + self.Body +" Event: "+self.Event
}
