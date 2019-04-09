package priority_task_mgr

import (
	"fmt"
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"testing"
	"time"
)

func taskCallBack(g goroutine_mgr.Goroutine, i interface{}) error {
	fmt.Println("goroutine:", g.GetName(), "| print i:", i)
	time.Sleep(1 * time.Millisecond)
	return nil
}

func taskCmpFunc(l interface{}, r interface{}) int {
	taskLeft := l.(PriorityTask)
	taskRight := r.(PriorityTask)
	taskPriorityLeft := taskLeft.taskPriority.(int)
	taskPriorityRight := taskRight.taskPriority.(int)
	if taskPriorityLeft < taskPriorityRight {
		return -1
	} else if taskPriorityLeft == taskPriorityRight {
		return 0
	} else {
		return 1
	}
}

func TestAll(t *testing.T) {
	taskMgr := new(PriorityTaskMgr)
	taskMgr.Initialise("TestTask", 5, taskCmpFunc, true, 0)
	taskMgr.Run(true)

	for index := 0; index < 10; index++ {
		task := new(PriorityTask)
		task.taskFunc = taskCallBack
		task.taskArgs = index
		task.taskPriority = index
		taskMgr.PushTask(*task)
	}

	time.Sleep(1 * time.Second)
}
