package graph

import "errors"

var (
	ErrNotRound = errors.New("argument not a round struct type")
	ErrNotNode  = errors.New("argument not a Neo4j node type")
	ErrNotSlice = errors.New("argument not a slice type")
)
