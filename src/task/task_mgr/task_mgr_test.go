package task_mgr

import (
	"fmt"
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"testing"
	"time"
)

func taskCallBack(g goroutine_mgr.Goroutine, i interface{}) error {
	fmt.Println("goroutine:", g.GetName(), "| print i:", i)
	return nil
}

func TestAll(t *testing.T) {
	taskMgr := new(TaskMgr)
	taskMgr.Initialise("TestTask", 5, true, 0)
	taskMgr.Run(true)

	index := 0
	for {
		task := new(Task)
		task.taskFunc = taskCallBack
		task.taskArgs = index
		taskMgr.PushTask(*task)
		index = index + 1
		time.Sleep(1 * time.Second)
	}
}
