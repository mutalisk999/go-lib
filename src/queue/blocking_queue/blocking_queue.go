package blocking_queue

import (
	"errors"
	"sync"
)

type BlockingQueue struct {
	mutex     *sync.RWMutex
	emptyCond *sync.Cond
	fullCond  *sync.Cond
	fullSize  uint64
	queue     []interface{}
}

func (b *BlockingQueue) Initialise(maxSize uint64) {
	b.mutex = new(sync.RWMutex)
	b.emptyCond = sync.NewCond(b.mutex)
	if maxSize != 0 {
		b.fullCond = sync.NewCond(b.mutex)
	} else {
		b.fullCond = nil
	}
	b.fullSize = maxSize
	b.queue = make([]interface{}, 0)
}

func (b *BlockingQueue) QueueSize() int {
	b.mutex.Lock()
	queueSize := len(b.queue)
	b.mutex.Unlock()
	return queueSize
}

func (b *BlockingQueue) PopFront() (interface{}, error) {
	b.mutex.Lock()
	queueSize := len(b.queue)
	if queueSize <= 0 {
		b.mutex.Unlock()
		return nil, errors.New("PopFront: empty queue")
	} else {
		queueElem := b.queue[0]
		b.queue = b.queue[1:]
		b.mutex.Unlock()
		return queueElem, nil
	}
	b.mutex.Unlock()
	return nil, nil
}

func (b *BlockingQueue) PopFrontBlocking() interface{} {
	b.mutex.Lock()
	for {
		queueSize := len(b.queue)
		if queueSize <= 0 {
			b.emptyCond.Wait()
		} else {
			break
		}
	}
	queueElem := b.queue[0]
	b.queue = b.queue[1:]

	if b.emptyCond != nil {
		b.emptyCond.Broadcast()
	}

	b.mutex.Unlock()
	return queueElem
}

func (b *BlockingQueue) PushBack(elem interface{}) error {
	b.mutex.Lock()
	queueSize := len(b.queue)
	if uint64(queueSize) >= b.fullSize {
		b.mutex.Unlock()
		return errors.New("PushBack: full queue")
	} else {
		b.queue = append(b.queue, elem)
		b.mutex.Unlock()
		return nil
	}
	b.mutex.Unlock()
	return nil
}

func (b *BlockingQueue) PushBackBlocking(elem interface{}) {
	b.mutex.Lock()
	if b.fullCond != nil {
		for {
			queueSize := len(b.queue)
			if uint64(queueSize) >= b.fullSize {
				b.fullCond.Wait()
			} else {
				break
			}
		}
	}
	b.queue = append(b.queue, elem)
	b.emptyCond.Broadcast()
	b.mutex.Unlock()
}

func (b *BlockingQueue) Destroy() {
	b.mutex = nil
	b.emptyCond = nil
	b.fullCond = nil
	b.fullSize = 0
	b.queue = make([]interface{}, 0)
}
