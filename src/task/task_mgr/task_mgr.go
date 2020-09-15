package task_mgr

import (
	"github.com/mutalisk999/go-lib/src/queue/blocking_queue"
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"strconv"
)

type Task struct {
	taskFunc func(goroutine_mgr.Goroutine, interface{}) error
	taskArgs interface{}
}

type TaskMgr struct {
	taskMgrName       string
	goroutineCount    uint64
	goroutineMgr      *goroutine_mgr.GoroutineManager
	goroutineQuit     []chan error
	taskQueue         *blocking_queue.BlockingQueue
	taskQueueBlocking bool
}

func (t *TaskMgr) Initialise(taskMgrName string, goroutineCount uint64, taskQueueBlocking bool, taskQueueMaxSize uint64) {
	t.taskMgrName = taskMgrName
	t.goroutineCount = goroutineCount
	t.goroutineMgr = new(goroutine_mgr.GoroutineManager)
	t.goroutineMgr.Initialise(t.taskMgrName + ".GoroutineMgr")
	t.goroutineQuit = make([]chan error, 0)
	t.taskQueue = new(blocking_queue.BlockingQueue)
	t.taskQueue.Initialise(taskQueueMaxSize)
	t.taskQueueBlocking = taskQueueBlocking
}

func (t *TaskMgr) PushTask(task Task) error {
	if t.taskQueueBlocking {
		t.taskQueue.PushBackBlocking(task)
		return nil
	} else {
		err := t.taskQueue.PushBack(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TaskMgr) PopTask() (Task, error) {
	if t.taskQueueBlocking {
		elem := t.taskQueue.PopFrontBlocking()
		return elem.(Task), nil
	} else {
		elem, err := t.taskQueue.PopFront()
		if err != nil {
			return Task{}, err
		}
		return elem.(Task), nil
	}
}

func (t *TaskMgr) runCallBackBlocking(g goroutine_mgr.Goroutine) {
	defer g.OnQuit()
	for {
		task, _ := t.PopTask()
		_ = task.taskFunc(g, task.taskArgs)
	}
}

func (t *TaskMgr) runCallBack(g goroutine_mgr.Goroutine, chanIndex interface{}, detach interface{}) {
	defer g.OnQuit()
	for {
		task, err := t.PopTask()
		if err != nil {
			break
		}
		_ = task.taskFunc(g, task.taskArgs)
	}
	if !detach.(bool) {
		t.goroutineQuit[chanIndex.(int)] <- nil
	}
}

func (t *TaskMgr) Run(detach bool) {
	if t.taskQueueBlocking {
		for i := 0; i < int(t.goroutineCount); i++ {
			t.goroutineMgr.GoroutineCreateP0(t.taskMgrName+".Goroutine"+strconv.Itoa(i), t.runCallBackBlocking)
		}
	} else {
		for i := 0; i < int(t.goroutineCount); i++ {
			if !detach {
				ch := make(chan error)
				t.goroutineQuit = append(t.goroutineQuit, ch)
			}
			t.goroutineMgr.GoroutineCreateP2(t.taskMgrName+".Goroutine"+strconv.Itoa(i), t.runCallBack, i, detach)
		}
	}
}

func (t *TaskMgr) Wait() {
	for i := 0; i < len(t.goroutineQuit); i++ {
		<-t.goroutineQuit[i]
	}
}

func (t *TaskMgr) RunAndWait(detach bool) {
	t.Run(detach)
	t.Wait()
}

func (t *TaskMgr) Destroy() {
	t.taskMgrName = ""
	t.goroutineCount = 0
	t.goroutineMgr.Destroy()
	t.goroutineMgr = nil
	t.goroutineQuit = make([]chan error, 0)
	t.taskQueue.Destroy()
	t.taskQueue = nil
	t.taskQueueBlocking = false
}
