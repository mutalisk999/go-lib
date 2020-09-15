package blocking_priority_queue

import (
	"fmt"
	"testing"
)

func intCmp(l interface{}, r interface{}) int {
	intLeft := l.(int)
	intRight := r.(int)
	if intLeft < intRight {
		return -1
	} else if intLeft == intRight {
		return 0
	} else {
		return 1
	}
}

func TestSmaller(t *testing.T) {
	queue := new(BlockingPriorityQueue)
	_ = queue.Initialise(1, intCmp, 0)

	_ = queue.Push(8)
	_ = queue.Push(5)
	_ = queue.Push(10)
	_ = queue.Push(2)
	_ = queue.Push(1)

	for {
		if queue.QueueSize() > 0 {
			v, _ := queue.Pop()
			fmt.Println("Pop: ", v.(int))
		} else {
			break
		}
	}
}

func TestBigger(t *testing.T) {
	queue := new(BlockingPriorityQueue)
	_ = queue.Initialise(2, intCmp, 0)

	_ = queue.Push(8)
	_ = queue.Push(5)
	_ = queue.Push(10)
	_ = queue.Push(2)
	_ = queue.Push(1)

	for {
		if queue.QueueSize() > 0 {
			v, _ := queue.Pop()
			fmt.Println("Pop: ", v.(int))
		} else {
			break
		}
	}
}
