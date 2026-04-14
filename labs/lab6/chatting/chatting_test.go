package chatting

import (
	"actors"
	"reflect"
	"testing"
)

var defaultQueueSize uint = 10

// Test 2.0
func TestBasicUser(t *testing.T) {
	bob := &UserActor{name: "Bob"}
	alice := &UserActor{name: "Alice"}

	bobRef := actors.New(bob, defaultQueueSize)
	aliceRef := actors.New(alice, defaultQueueSize)

	bobRef.Send(&ChatMessage{sender: "Alice", text: "Hello Bob!"})
	aliceRef.Send(&ChatMessage{sender: "Bob", text: "Wassup Alice!!!!!"})

	bobRef.Kill()
	aliceRef.Kill()

	bobExpected := []string{"From Alice: Hello Bob!"}
	if !reflect.DeepEqual(bob.messagesReceived, bobExpected) {
		t.Errorf("Got Bob's messages: %v, expected: %v", bob.messagesReceived, bobExpected)
	}
	aliceExpected := []string{"From Bob: Wassup Alice!!!!!"}
	if !reflect.DeepEqual(alice.messagesReceived, aliceExpected) {
		t.Errorf("Got Alice's messages: %v, expected: %v", alice.messagesReceived, aliceExpected)
	}
}

// Test 2.1
func TestUnkownUser(t *testing.T) {
	john := &UserActor{}
	doe := &UserActor{}

	johnRef := actors.New(john, defaultQueueSize)
	doeRef := actors.New(doe, defaultQueueSize)

	johnRef.Send(&ChatMessage{text: "Hello stranger!"})
	doeRef.Send(&ChatMessage{text: "Hello unknown..."})

	johnRef.Kill()
	doeRef.Kill()

	johnExpected := []string{"From Unknown: Hello stranger!"}
	if !reflect.DeepEqual(john.messagesReceived, johnExpected) {
		t.Errorf("Got John's messages: %v, expected: %v", john.messagesReceived, johnExpected)
	}
	doeExpected := []string{"From Unknown: Hello unknown..."}
	if !reflect.DeepEqual(doe.messagesReceived, doeExpected) {
		t.Errorf("Got Doe's messages: %v, expected: %v", doe.messagesReceived, doeExpected)
	}
}

// Test 2.2
func TestBasicChatRoom(t *testing.T) {
	bob := &UserActor{name: "Bob"}
	alice := &UserActor{name: "Alice"}
	bobRef := actors.New(bob, defaultQueueSize)
	aliceRef := actors.New(alice, defaultQueueSize)

	chatRoom := &ChatRoomActor{name: "Awesome-sauce"}
	chatRoomRef := actors.New(chatRoom, defaultQueueSize)

	chatRoomRef.Send(&ChatRoomMessage{command: Join, from: bobRef, from_name: bob.name})
	chatRoomRef.Send(&ChatRoomMessage{command: Join, from: aliceRef, from_name: alice.name})

	chatRoomRef.Send(&ChatRoomMessage{command: SendMessage, from: bobRef, from_name: "Bob", text: "Hey Alice it's Bob!"})

	chatRoomRef.Send(&ChatRoomMessage{command: Leave, from: bobRef, from_name: bob.name})
	chatRoomRef.Send(&ChatRoomMessage{command: Leave, from: aliceRef, from_name: alice.name})

	chatRoomRef.Kill()
	bobRef.Kill()
	aliceRef.Kill()

	logsExpected := []string{
		"Bob joined room Awesome-sauce",
		"Alice joined room Awesome-sauce",
		"Bob's message sent sucessfully",
		"Bob left room Awesome-sauce",
		"Alice left room Awesome-sauce",
	}
	if !reflect.DeepEqual(chatRoom.logs, logsExpected) {
		t.Errorf("Got chat logs: %v, expected: %v", chatRoom.logs, logsExpected)
	}
	aliceExpected := []string{"From Bob: Hey Alice it's Bob!"}
	if !reflect.DeepEqual(alice.messagesReceived, aliceExpected) {
		t.Errorf("Got Alice's messages: %v, expected: %v", alice.messagesReceived, aliceExpected)
	}
	if len(bob.messagesReceived) != 0 {
		t.Errorf("Bob should not have received any messages")
	}
}

// Test 2.3
func TestMultipleChatRooms(t *testing.T) {
	bryan := &UserActor{name: "Bryan"}
	grace := &UserActor{name: "Grace"}
	alex := &UserActor{name: "Alex"}
	olivia := &UserActor{name: "Olivia"}
	bryanRef := actors.New(bryan, defaultQueueSize)
	graceRef := actors.New(grace, defaultQueueSize)
	alexRef := actors.New(alex, defaultQueueSize)
	oliviaRef := actors.New(olivia, defaultQueueSize)

	chatRoom1 := &ChatRoomActor{name: "1"}
	chatRoom1Ref := actors.New(chatRoom1, defaultQueueSize)
	chatRoom2 := &ChatRoomActor{name: "2"}
	chatRoom2Ref := actors.New(chatRoom2, defaultQueueSize)

	chatRoom1Ref.Send(&ChatRoomMessage{command: Leave, from: bryanRef, from_name: "Bryan"})

	chatRoom1Ref.Send(&ChatRoomMessage{command: Join, from: bryanRef, from_name: "Bryan"})
	chatRoom1Ref.Send(&ChatRoomMessage{command: Join, from: graceRef, from_name: "Grace"})

	chatRoom2Ref.Send(&ChatRoomMessage{command: Join, from: graceRef, from_name: "Grace"})
	chatRoom2Ref.Send(&ChatRoomMessage{command: Join, from: alexRef, from_name: "Alex"})
	chatRoom2Ref.Send(&ChatRoomMessage{command: Join, from: oliviaRef, from_name: "Olivia"})

	chatRoom1Ref.Send(&ChatRoomMessage{command: SendMessage, from: alexRef, from_name: "Alex", text: "Hello room 1!"})
	chatRoom2Ref.Send(&ChatRoomMessage{command: SendMessage, from: alexRef, from_name: "Alex", text: "Hello room 2!"})
	chatRoom2Ref.Send(&ChatRoomMessage{command: SendMessage, from: graceRef, from_name: "Grace", text: "Hey Alex!"})
	chatRoom1Ref.Send(&ChatRoomMessage{command: SendMessage, from: graceRef, from_name: "Grace", text: "Hey people in room 1!!"})

	chatRoom1Ref.Kill()
	chatRoom2Ref.Kill()
	bryanRef.Kill()
	graceRef.Kill()
	alexRef.Kill()
	oliviaRef.Kill()

	logsExpected1 := []string{
		"Bryan is not in room 1",
		"Bryan joined room 1",
		"Grace joined room 1",
		"Alex is not in room 1",
		"Grace's message sent sucessfully",
	}
	if !reflect.DeepEqual(chatRoom1.logs, logsExpected1) {
		t.Errorf("Got chat logs for room 1: %v, expected: %v", chatRoom1.logs, logsExpected1)
	}
	bryanMsgExpected := []string{"From Grace: Hey people in room 1!!"}
	if !reflect.DeepEqual(bryan.messagesReceived, bryanMsgExpected) {
		t.Errorf("Bryan's messages got: %v, expected %v", bryan.messagesReceived, bryanMsgExpected)
	}

	logsExpected2 := []string{
		"Grace joined room 2",
		"Alex joined room 2",
		"Olivia joined room 2",
		"Alex's message sent sucessfully",
		"Grace's message sent sucessfully",
	}
	if !reflect.DeepEqual(chatRoom2.logs, logsExpected2) {
		t.Errorf("Got chat logs for room 2: %v, expected: %v", chatRoom2.logs, logsExpected2)
	}

	firstMsgExpected := []string{"From Alex: Hello room 2!"}
	secondMsgExpected := []string{"From Grace: Hey Alex!"}
	if !reflect.DeepEqual(grace.messagesReceived, firstMsgExpected) {
		t.Errorf("Got Grace's messages: %v, expected: %v", grace.messagesReceived, firstMsgExpected)
	}
	if !reflect.DeepEqual(alex.messagesReceived, secondMsgExpected) {
		t.Errorf("Got Alex's messages: %v, expected: %v", alex.messagesReceived, secondMsgExpected)
	}
	oliviaMsgExpected := []string{firstMsgExpected[0], secondMsgExpected[0]}
	if !reflect.DeepEqual(olivia.messagesReceived, oliviaMsgExpected) {
		t.Errorf("Got Olivia's messages: %v, expected: %v", olivia.messagesReceived, oliviaMsgExpected)
	}
}
