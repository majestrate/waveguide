package apub


const Context = "https://www.w3.org/ns/activitystreams"

type Message interface {
	Type() MessageType
	MarshalJSON() ([]byte, error)
}
