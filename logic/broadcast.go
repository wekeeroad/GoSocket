package logic

type broadcaster struct {
	users map[string]*User

	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	checkUserChannel      chan string
	checkUserCanInChannel chan bool
	requestUsersChannel   chan struct{}
	usersChannel          chan []*User
}

var Broadcaster = &broadcaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message, 1024),

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
	requestUsersChannel:   make(chan struct{}),
	usersChannel:          make(chan []*User),
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname
	return <-b.checkUserCanInChannel
}

func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			b.users[user.NickName] = user
			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			delete(b.users, user.NickName)
			user.CloseMessageChannel()
		case msg := <-b.messageChannel:
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
			OfflineProcessor.Save(msg)
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case <-b.requestUsersChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}
			b.usersChannel <- userList
		}
	}
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	b.messageChannel <- msg
}

func (b *broadcaster) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}
