package main

type SkiplistLevel struct {
	forward *SkiplistNode
	walk    uint64
}
type SkiplistNode struct {
	val   Compare
	score float64

	backward *SkiplistNode

	level []SkiplistLevel
}

func NewSkiplistNode(val Compare, score float64, level int) *SkiplistNode {
	n := &SkiplistNode{
		val:   val,
		score: score,
	}
	n.level = make([]SkiplistLevel, level)
	return n
}
