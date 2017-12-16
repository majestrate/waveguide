package apub

import (
	"errors"
	"strings"
)

type MessageType string

func (t MessageType) Equals(other MessageType) bool {
	return t.Lower() == other.Lower()
}

func (t MessageType) Is(str string) bool {
	return t.Lower() == strings.ToLower(str)
}

func (t MessageType) Lower() string {
	return strings.ToLower(t.String())
}

func (t MessageType) String() string {
	return string(t)
}

const TypeEmpty = MessageType("")

const TypePerson = MessageType("Person")
const TypeNote = MessageType("Note")
const TypeCreate = MessageType("Create")
const TypeLike = MessageType("Like")
const TypeCollection = MessageType("Collection")

var MessageTypes = []MessageType{TypePerson, TypeNote, TypeCreate, TypeLike, TypeCollection}
var ErrBadMessageType = errors.New("Bad message type")

// parse string to const MessageType
func TypeFromString(str string) (MessageType, error) {
	for _, t := range MessageTypes {
		if t.Is(str) {
			return t, nil
		}
	}
	return TypeEmpty, ErrBadMessageType
}
