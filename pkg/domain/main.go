package domain

import "fmt"

type Node struct {
	Path  string
	Value string
	Type  string
}

type ValueStatus int

const (
	Eq ValueStatus = iota
	Ne
	Wt
)

func (src ValueStatus) ToString() string {
	switch src {
	case Eq:
		return "Eq"
	case Ne:
		return "Ne"
	case Wt:
		return "WrongType"
	default:
		panic(fmt.Sprintf("Unknown ValueStatus type %v", src))
	}
}

type LinkStatus int

const (
	Present LinkStatus = iota
	NoLeft
	NoRight
)

func (src LinkStatus) ToString() string {
	switch src {
	case Present:
		return "Present"
	case NoLeft:
		return "NoLeft"
	case NoRight:
		return "NoRight"
	default:
		panic(fmt.Sprintf("Unknown LinkStatus type %v", src))
	}
}

type Link struct {
	Path  string
	Left  Node
	Right Node
}

func (src Link) GetValueStatus() ValueStatus {
	if src.Left.Type != src.Right.Type {
		return Wt
	}
	if src.Left.Value != src.Right.Value {
		return Ne
	}
	return Eq
}

func (src Link) GetLinkStatus() LinkStatus {
	if src.Left == (Node{}) {
		return NoLeft
	}
	if src.Right == (Node{}) {
		return NoRight
	}
	return Present
}
