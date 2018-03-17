package adn

import (
	"fmt"
	"strconv"
)

type UID string

func (uid UID) String() string {
	return string(uid)
}

func (uid UID) Int() int64 {
	i, _ := strconv.ParseInt(uid.String(), 10, 64)
	return i
}

func ParseUID(str string) (UID, error) {
	return UID(str), nil
}

type ChanID int64

func (i ChanID) String() string {
	return fmt.Sprintf("%d", i)
}

func (i ChanID) Int() int64 {
	return int64(i)
}

func ParseChanID(str string) (ChanID, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return ChanID(i), err
}
