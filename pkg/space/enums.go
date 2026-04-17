package space

import "encoding/json"

type SubjectType int

const (
	SubjectUser         SubjectType = 1
	SubjectMessageBoard SubjectType = iota * 10
	SubjectArticle
	SubjectComment
)

//------------------------------------------------------------------------------

type DataEventType int

const (
	DataEventTypeCreate DataEventType = iota + 1
	DataEventTypeRead
	DataEventTypeUpdate
	DataEventTypeDelete
)

type DataEvent[T any] struct {
	EventType DataEventType
	Data      T
}

func (de *DataEvent[T]) MarshalBinary() (data []byte, err error) {
	return json.Marshal(de)
}

func (de *DataEvent[T]) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, de)
}
