package priority_task_mgr

import (
	"github.com/mutalisk999/go-lib/src/queue/blocking_priority_queue"
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"strconv"
)

type PriorityTask struct {
	taskFunc     func(goroutine_mgr.Goroutine, interface{}) error
	taskArgs     interface{}
	taskPriority interface{}
}

type PriorityTaskMgr struct {
	taskMgrName       string
	goroutineCount    uint64
	goroutineMgr      *goroutine_mgr.GoroutineManager
	goroutineQuit     []chan error
	taskQueue         *blocking_priority_queue.BlockingPriorityQueue
	taskQueueBlocking bool
}

func (t *PriorityTaskMgr) Initialise(taskMgrName string, goroutineCount uint64, taskCmpFunc func(l interface{}, r interface{}) int,
	taskQueueBlocking bool, taskQueueMaxSize uint64) {
	t.taskMgrName = taskMgrName
	t.goroutineCount = goroutineCount
	t.goroutineMgr = new(goroutine_mgr.GoroutineManager)
	t.goroutineMgr.Initialise(t.taskMgrName + ".GoroutineMgr")
	t.goroutineQuit = make([]chan error, 0)
	t.taskQueue = new(blocking_priority_queue.BlockingPriorityQueue)
	_ = t.taskQueue.Initialise(2, taskCmpFunc, taskQueueMaxSize)
	t.taskQueueBlocking = taskQueueBlocking
}

func (t *PriorityTaskMgr) PushTask(task PriorityTask) error {
	if t.taskQueueBlocking {
		t.taskQueue.PushBlocking(task)
		return nil
	} else {
		err := t.taskQueue.Push(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *PriorityTaskMgr) PopTask() (PriorityTask, error) {
	if t.taskQueueBlocking {
		elem := t.taskQueue.PopBlocking()
		return elem.(PriorityTask), nil
	} else {
		elem, err := t.taskQueue.Pop()
		if err != nil {
			return PriorityTask{}, err
		}
		return elem.(PriorityTask), nil
	}
}

func (t *PriorityTaskMgr) runCallBackBlocking(g goroutine_mgr.Goroutine) {
	defer g.OnQuit()
	for {
		task, _ := t.PopTask()
		_ = task.taskFunc(g, task.taskArgs)
	}
}

func (t *PriorityTaskMgr) runCallBack(g goroutine_mgr.Goroutine, chanIndex interface{}, detach interface{}) {
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

func (t *PriorityTaskMgr) Run(detach bool) {
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

func (t *PriorityTaskMgr) Wait() {
	for i := 0; i < len(t.goroutineQuit); i++ {
		<-t.goroutineQuit[i]
	}
}

func (t *PriorityTaskMgr) RunAndWait(detach bool) {
	t.Run(detach)
	t.Wait()
}

func (t *PriorityTaskMgr) Destroy() {
	t.taskMgrName = ""
	t.goroutineCount = 0
	t.goroutineMgr.Destroy()
	t.goroutineMgr = nil
	t.goroutineQuit = make([]chan error, 0)
	t.taskQueue.Destroy()
	t.taskQueue = nil
	t.taskQueueBlocking = false
}
