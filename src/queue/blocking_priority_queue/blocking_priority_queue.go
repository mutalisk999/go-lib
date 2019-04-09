package blocking_priority_queue

import (
	"errors"
	"fmt"
	"sync"
)

type BlockingPriorityQueue struct {
	mutex     *sync.RWMutex
	emptyCond *sync.Cond
	fullCond  *sync.Cond
	fullSize  uint64
	queueType int
	cmpFunc   func(l interface{}, r interface{}) int
	queue     []interface{}
}

func swapQueueElem(queue []interface{}, leftIndex int, rightIndex int) error {
	if leftIndex < 0 || leftIndex >= len(queue) {
		return errors.New(fmt.Sprintf("invalid leftIndex: %d", leftIndex))
	}
	if rightIndex < 0 || rightIndex >= len(queue) {
		return errors.New(fmt.Sprintf("invalid rightIndex: %d", rightIndex))
	}
	queue[leftIndex], queue[rightIndex] = queue[rightIndex], queue[leftIndex]
	return nil
}

func (b *BlockingPriorityQueue) Initialise(queueType int, cmpFunc func(l interface{}, r interface{}) int, maxSize uint64) error {
	if queueType != 1 && queueType != 2 {
		return errors.New(fmt.Sprintf("invalid queueType: %d", queueType))
	}
	b.mutex = new(sync.RWMutex)
	b.emptyCond = sync.NewCond(b.mutex)
	if maxSize != 0 {
		b.fullCond = sync.NewCond(b.mutex)
	} else {
		b.fullCond = nil
	}
	b.fullSize = maxSize
	// 1 -- small root heap
	// 2 -- big root heap
	b.queueType = queueType
	b.cmpFunc = cmpFunc
	b.queue = make([]interface{}, 0)
	return nil
}

func (b BlockingPriorityQueue) QueueSize() int {
	b.mutex.RLock()
	queueSize := len(b.queue)
	b.mutex.RUnlock()
	return queueSize
}

func (b BlockingPriorityQueue) Top() (interface{}, error) {
	b.mutex.Lock()
	queueSize := len(b.queue)
	if queueSize <= 0 {
		b.mutex.Unlock()
		return nil, errors.New("Top: empty queue")
	} else {
		queueElem := b.queue[0]
		b.mutex.Unlock()
		return queueElem, nil
	}
	b.mutex.Unlock()
	return nil, nil
}

func (b *BlockingPriorityQueue) siftDown() {
	if len(b.queue) > 0 {
		parentIndex := 0
		for {
			leftChildIndex := (parentIndex+1)*2 - 1
			rightChildIndex := (parentIndex + 1) * 2

			if leftChildIndex >= len(b.queue) {
				// invalid leftChildIndex and invalid rightChildIndex
				break
			} else if rightChildIndex >= len(b.queue) {
				// valid leftChildIndex and invalid rightChildIndex
				if b.queueType == 1 {
					if b.cmpFunc(b.queue[parentIndex], b.queue[leftChildIndex]) > 0 {
						swapQueueElem(b.queue, parentIndex, leftChildIndex)
						parentIndex = leftChildIndex
					} else {
						break
					}
				} else if b.queueType == 2 {
					if b.cmpFunc(b.queue[parentIndex], b.queue[leftChildIndex]) < 0 {
						swapQueueElem(b.queue, parentIndex, leftChildIndex)
						parentIndex = leftChildIndex
					} else {
						break
					}
				}
			} else {
				// valid leftChildIndex and valid rightChildIndex
				if b.queueType == 1 {
					smallerIndex := -1
					if b.cmpFunc(b.queue[leftChildIndex], b.queue[rightChildIndex]) < 0 {
						smallerIndex = leftChildIndex
					} else {
						smallerIndex = rightChildIndex
					}
					if b.cmpFunc(b.queue[parentIndex], b.queue[smallerIndex]) > 0 {
						swapQueueElem(b.queue, parentIndex, smallerIndex)
						parentIndex = smallerIndex
					} else {
						break
					}
				} else if b.queueType == 2 {
					biggerIndex := -1
					if b.cmpFunc(b.queue[leftChildIndex], b.queue[rightChildIndex]) > 0 {
						biggerIndex = leftChildIndex
					} else {
						biggerIndex = rightChildIndex
					}
					if b.cmpFunc(b.queue[parentIndex], b.queue[biggerIndex]) < 0 {
						swapQueueElem(b.queue, parentIndex, biggerIndex)
						parentIndex = biggerIndex
					} else {
						break
					}
				}
			}
		}
	}
}

func (b *BlockingPriorityQueue) Pop() (interface{}, error) {
	b.mutex.Lock()
	queueSize := len(b.queue)
	if queueSize <= 0 {
		b.mutex.Unlock()
		return nil, errors.New("Pop: empty queue")
	} else {
		swapQueueElem(b.queue, 0, len(b.queue)-1)
		queueElem := b.queue[len(b.queue)-1]
		b.queue = b.queue[:len(b.queue)-1]
		// sift down
		b.siftDown()
		b.mutex.Unlock()
		return queueElem, nil
	}
	b.mutex.Unlock()
	return nil, nil
}

func (b *BlockingPriorityQueue) PopBlocking() interface{} {
	b.mutex.Lock()
	for {
		queueSize := len(b.queue)
		if queueSize <= 0 {
			b.emptyCond.Wait()
		} else {
			break
		}
	}

	swapQueueElem(b.queue, 0, len(b.queue)-1)
	queueElem := b.queue[len(b.queue)-1]
	b.queue = b.queue[:len(b.queue)-1]
	// sift down
	b.siftDown()

	if b.fullCond != nil {
		b.fullCond.Broadcast()
	}

	b.mutex.Unlock()
	return queueElem
}

func (b *BlockingPriorityQueue) siftUp() {
	if len(b.queue) > 0 {
		childIndex := len(b.queue) - 1
		for {
			parentIndex := (childIndex+1)/2 - 1
			if parentIndex >= 0 {
				if b.queueType == 1 {
					if b.cmpFunc(b.queue[parentIndex], b.queue[childIndex]) > 0 {
						swapQueueElem(b.queue, parentIndex, childIndex)
						childIndex = parentIndex
					} else {
						break
					}
				} else if b.queueType == 2 {
					if b.cmpFunc(b.queue[parentIndex], b.queue[childIndex]) < 0 {
						swapQueueElem(b.queue, parentIndex, childIndex)
						childIndex = parentIndex
					} else {
						break
					}
				}
			} else {
				break
			}
		}
	}
}

func (b *BlockingPriorityQueue) Push(elem interface{}) error {
	b.mutex.Lock()
	queueSize := len(b.queue)
	if b.fullSize != 0 && uint64(queueSize) >= b.fullSize {
		b.mutex.Unlock()
		return errors.New("Push: full queue")
	} else {
		b.queue = append(b.queue, elem)
		// sift up
		b.siftUp()
		b.mutex.Unlock()
		return nil
	}
	b.mutex.Unlock()
	return nil
}

func (b *BlockingPriorityQueue) PushBlocking(elem interface{}) {
	b.mutex.Lock()
	if b.fullCond != nil {
		for {
			queueSize := len(b.queue)
			if b.fullSize != 0 && uint64(queueSize) >= b.fullSize {
				b.fullCond.Wait()
			} else {
				break
			}
		}
	}
	b.queue = append(b.queue, elem)
	// sift up
	b.siftUp()
	b.emptyCond.Broadcast()
	b.mutex.Unlock()
}

func (b *BlockingPriorityQueue) Destroy() {
	b.mutex = nil
	b.emptyCond = nil
	b.fullCond = nil
	b.fullSize = 0
	b.queueType = 0
	b.cmpFunc = nil
	b.queue = make([]interface{}, 0)
}
