package chatting

import (
	"actors"
)

// Defines the actions associated with the chat room
type ChatCommand int

const (
	Join ChatCommand = iota
	Leave
	SendMessage
)

// Messages that UserActors can send to the ChatRoomActor
type ChatRoomMessage struct {
	command   ChatCommand      // Command telling the ChatRoomActor what to do
	from      *actors.ActorRef // Reference of user sending this message
	from_name string           // Name of user sending message (done this way to not expose state of actor)
	text      string           // The actual contents of the chat message (only required if command == SendMessage)
}

type ChatRoomActor struct {
	name    string             // Name of this chat rooom
	members []*actors.ActorRef // the user actors in this room
	logs    []string           // A log of every action this chat room has had
}

// TODO: Implement this function to handle the 3 different commands in the
// correct output format to pass the autograder
func (room *ChatRoomActor) Receive(msg any, self *actors.ActorRef) {
	panic("TODO")
}
