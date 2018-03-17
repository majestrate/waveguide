package adn

import (
	"fmt"
	"strconv"
)

type UID int64

func (i UID) String() string {
	return fmt.Sprintf("%d", i)
}

func (i UID) Int() int64 {
	return int64(i)
}

func ParseUID(str string) (UID, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return UID(i), err
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
