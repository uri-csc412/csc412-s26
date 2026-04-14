package actors

import (
	"testing"
	"time"
)

// ============================ EchoActor ==========================
// Implements a concret actor implementation that is simply updating a string
type EchoActor struct {
	currentStr string        // state of the actor
	done       chan struct{} // If the current message has been sent
}

func (a *EchoActor) Receive(msg any, self *ActorRef) {
	a.currentStr = msg.(string)

	if a.done != nil {
		a.done <- struct{}{}
	}
}

func helperSendMessage(ref *ActorRef, t *testing.T, str string) {
	if err := ref.Send(str); err != nil {
		t.Fatal(err)
	}

	// Wait for actor to finish
	<-ref.actor.(*EchoActor).done

	if ref.actor.(*EchoActor).currentStr != str {
		t.Errorf("Got: %v, expected: %v", ref.actor.(*EchoActor).currentStr, str)
	}
}

// ============================ CounterActor =======================
// Another concret actor implementation implementing a counter

type CounterActor struct {
	counter int           // state of the actor
	done    chan struct{} // If the current send operation is done
}

func (a *CounterActor) Receive(msg any, self *ActorRef) {
	a.counter += msg.(int)
}

func helperAddCounter(ref *ActorRef, t *testing.T, count int) {
	if err := ref.Send(count); err != nil {
		t.Fatal(err)
	}
}

// ============================== Start of Test ====================

// Test 1.0: Creates an Echo actor with various message
func TestEchoActor(t *testing.T) {
	actor := &EchoActor{done: make(chan struct{}, 10)}
	ref := New(actor, 10)

	helperSendMessage(ref, t, "Hello from this test")
	helperSendMessage(ref, t, "Wow let's do this thing again")
	helperSendMessage(ref, t, "Last time I promise")

	ref.Kill()
}

// Test 1.1: Create a Counter actor changing the value of the counter
func TestCounterActor(t *testing.T) {
	actor := &CounterActor{}
	ref := New(actor, 5)

	helperAddCounter(ref, t, 5)
	helperAddCounter(ref, t, 10)
	helperAddCounter(ref, t, 1000000)
	helperAddCounter(ref, t, -10)
	helperAddCounter(ref, t, 0)

	// Wait for all the routines to be done
	ref.Kill()

	expected := 5 + 10 + 1000000 - 10
	if actor.counter != expected {
		t.Errorf("Incorrect counter. Got: %v, Expected: %v", actor.counter, expected)
	}
}

// Test 1.2: Checks to see how your implementation handles a full message queue
func TestFullQueue(t *testing.T) {
	actor := &CounterActor{}
	ref := New(actor, 3)

	// Sending 3 message is ok
	helperAddCounter(ref, t, 5)
	helperAddCounter(ref, t, 5)
	helperAddCounter(ref, t, 5)

	// On the 4th there's no more room
	if err := ref.Send(5); err == nil {
		t.Error("An error should be returned after sending a message to a full queue!")
	}

	ref.Kill()
}

type SlowActor struct{}

func (*SlowActor) Receive(msg any, self *ActorRef) {
	// Does nothing but wait for 1 second
	time.Sleep(time.Second)
}

// Test 1.3: Check to make sure you're Actor model is asynchronous
func TestAsync(t *testing.T) {
	actor := &SlowActor{}
	ref := New(actor, 5)

	start := time.Now()
	ref.Send(nil)
	elapsed := time.Since(start)

	if elapsed > 50*time.Millisecond {
		t.Error("Your receive function returned right away. Your implementation is not asynchronous!")
	}

	ref.Kill()
}
