package logic

type broadcaster struct {
	users map[string]*User

	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	checkUserChannel      chan string
	checkUserCanInChannel chan bool
}

var Broadcaster = &broadcaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message, 1024),

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname
	return <-b.checkUserCanInChannel
}
