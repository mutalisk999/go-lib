package blocking_queue

import "testing"
import (
	. "github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"strconv"
	"fmt"
	"time"
)

func produce(g Goroutine, q interface{}) {
	queue := q.(*BlockingQueue)
	index := 1
	for {
		queue.PushBackBlocking(index)
		fmt.Println(g.GetName(), "produce:", index)
		index = index + 1
		time.Sleep(1*time.Second)
	}
}

func consume(g Goroutine, q interface{}) {
	queue := q.(*BlockingQueue)
	for {
		elem := queue.PopFrontBlocking()
		fmt.Println(g.GetName(), "consume:", elem.(int))
	}
}

func TestAll(t *testing.T) {
	manager := new(GoroutineManager)
	manager.Initialise("mgr1")

	queue := new(BlockingQueue)
	queue.Initialise(100)

	manager.GoroutineCreateP1("goroutine_producer", produce, queue)

	for i:=0; i<4; i++ {
		manager.GoroutineCreateP1("goroutine_consumer"+strconv.Itoa(i), consume, queue)
	}

	manager.GoroutineDump()
	fmt.Println("goroutineCount", manager.GoroutineCount())

	for {
		time.Sleep(1*time.Second)
	}
}
