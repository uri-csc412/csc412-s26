package actors

// Interface that defines the behavior for an Actor
type Actor interface {
	// Defines how the Actor handles an arbitrary message,
	// concrete implementations of actors will need to implement this
	// however they please
	Receive(msg any, self *ActorRef)
}

// Manages sending and receiving message for the Actor
type ActorRef struct {
	messageQueue chan any      // The queue of messages this actor needs to process
	actor        Actor         // The underlying actor provided by caller
	killed       chan struct{} // Whether or not this actor has been killed or not
}

// TODO: Create a new ActorRef to manage receiving messages in a separate goroutine
func New(actor Actor, messageQueueSize uint) *ActorRef {
	panic("TODO")
}

// TODO: Put a message onto the message queue channel (asynchronously) of this actor
func (ref *ActorRef) Send(msg any) error {
	panic("TODO")
}

// TODO: Shut down the give actor
func (ref *ActorRef) Kill() {
	panic("TODO")
}
