package main

import "math/rand"

const (
	MaxLevel = 32
	P        = 0.25
)

type Skiplist struct {
	head *SkiplistNode
	tail *SkiplistNode

	length uint64
	level  int
}

func NewSkiplist() *Skiplist {
	head := NewSkiplistNode(nil, 0, MaxLevel)
	for i := 0; i < MaxLevel; i++ {
		head.level[i].walk = 0
		head.level[i].forward = nil
	}
	head.backward = nil
	return &Skiplist{
		length: 0,
		level:  0,
		head:   head,
		tail:   nil,
	}
}

func randomLevel() int {
	level := 1
	for (float32((rand.Int() & 0xFFF)) < (P * 0xFFFF)) && level <= MaxLevel {
		level++
	}
	return level
}

func (list *Skiplist) InsertNode(score float64, val Compare) {
	// 每一层上的backward
	backwards := make([]*SkiplistNode, MaxLevel)
	levelWalks := make([]uint64, MaxLevel)
	current := list.head
	for i := list.level - 1; i >= 0; i-- {
		if i == list.level {
			levelWalks[i] = 0
		} else {
			levelWalks[i] = levelWalks[i+1]
		}
		for current.level[i].forward != nil &&
			(current.level[i].forward.score < score ||
				(current.level[i].forward.score == score &&
					val.CompareWith(current.level[i].forward.val) > 0)) {
			levelWalks[i] += current.level[i].walk
			current = current.level[i].forward
		}
		backwards[i] = current
	}
	level := randomLevel()
	if level > list.level {
		for i := list.level; i < level; i++ {
			levelWalks[i] = 0
			backwards[i] = list.head
			backwards[i].level[i].walk = list.length
		}
		list.level = level
	}
	current = NewSkiplistNode(val, score, level)
	for i := 0; i < level; i++ {
		current.level[i].forward = backwards[i].level[i].forward
		backwards[i].level[i].forward = current

		current.level[i].walk = backwards[i].level[i].walk - (levelWalks[0] - levelWalks[i])
		backwards[i].level[i].walk = (levelWalks[0] - levelWalks[i]) + 1
	}

	for i := level; i < list.level; i++ {
		backwards[i].level[i].walk++
	}
	if backwards[0] == list.head {
		current.backward = nil
	} else {
		current.backward = backwards[0]
	}
	if current.level[0].forward != nil {
		current.level[0].forward.backward = current
	} else {
		list.tail = current
	}
	list.length++
}

func (list *Skiplist) Delete(score float64, val Compare) {
	backwards := make([]*SkiplistNode, MaxLevel)
	current := list.head
	for i := list.level - 1; i >= 0; i-- {
		for current.level[i].forward != nil &&
			(current.level[i].forward.score < score ||
				(current.level[i].forward.score == score &&
					val.CompareWith(current.level[i].forward.val) > 0)) {
			current = current.level[i].forward
		}
		backwards[i] = current
	}
	current = current.level[0].forward
	if current.score != score || val.CompareWith(current.val) != 0 {
		return
	}
	for i := 0; i < list.level; i++ {
		backwards[i].level[i].walk--
		if backwards[i].level[i].forward == current {
			backwards[i].level[i].walk += current.level[i].walk
			backwards[i].level[i].forward = current.level[i].forward
		}
	}
	for i := 0; i < len(current.level); i++ {
		backwards[i].level[i].forward = current.level[i].forward
	}
	if current.level[0].forward != nil {
		current.level[0].forward.backward = current.backward
	} else {
		list.tail = current.backward
	}

	// 跳表当前层数减小
	for list.level > 1 && list.head.level[list.level-1].forward == nil {
		list.level--
	}
	// 跳表长度减小
	list.length--
}

func (list *Skiplist) Find() {

}
