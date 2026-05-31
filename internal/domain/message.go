package domain

import "time"

type Direction string

const (
	DirectionIn  Direction = "in"
	DirectionOut Direction = "out"
)

type Message struct {
	ID        int64
	UID       int64
	Direction Direction
	Text      string
	TgMsgID   int64
	At        time.Time
}
