package model

type State string

const (
	StateWin  State = "win"
	StateLost State = "lost"
)

var States = []string{
	string(StateWin),
	string(StateLost),
}
