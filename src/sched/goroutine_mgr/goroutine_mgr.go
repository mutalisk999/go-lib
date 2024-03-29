package goroutine_mgr

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Goroutine struct {
	goroutineId         uint64
	goroutineName       string
	goroutineCreateTime time.Time
	fromManager         *GoroutineManager
	goroutineCtx        *context.Context
	goroutineCancelFunc *context.CancelFunc
}

func (g *Goroutine) OnQuit() {
	g.fromManager.GoroutineRemove(g)
}

func (g *Goroutine) GetId() uint64 {
	return g.goroutineId
}

func (g *Goroutine) GetName() string {
	return g.goroutineName
}

func (g *Goroutine) GetContext() *context.Context {
	return g.goroutineCtx
}

func (g *Goroutine) GetCancelFunc() *context.CancelFunc {
	return g.goroutineCancelFunc
}

type GoroutineManager struct {
	mutex        *sync.RWMutex
	counter      uint64
	managerName  string
	goroutineMap map[uint64]*Goroutine
}

func (g *GoroutineManager) Initialise(managerName string) {
	g.mutex = new(sync.RWMutex)
	g.counter = 0
	g.managerName = managerName
	g.goroutineMap = make(map[uint64]*Goroutine)
}

func (g *GoroutineManager) constructGoroutine(goroutineName string) *Goroutine {
	goroutine := new(Goroutine)
	g.mutex.Lock()
	g.counter++
	goroutine.goroutineId = g.counter
	g.mutex.Unlock()
	goroutine.goroutineName = goroutineName
	goroutine.goroutineCreateTime = time.Now()
	goroutine.fromManager = g
	g.mutex.Lock()
	g.goroutineMap[goroutine.goroutineId] = goroutine
	g.mutex.Unlock()
	return goroutine
}

func (g *GoroutineManager) constructGoroutineWithCancel(goroutineName string) *Goroutine {
	goroutine := new(Goroutine)
	g.mutex.Lock()
	g.counter++
	goroutine.goroutineId = g.counter
	g.mutex.Unlock()
	goroutine.goroutineName = goroutineName
	goroutine.goroutineCreateTime = time.Now()
	ctx, cbCancel := context.WithCancel(context.Background())
	goroutine.goroutineCtx = &ctx
	goroutine.goroutineCancelFunc = &cbCancel
	goroutine.fromManager = g
	g.mutex.Lock()
	g.goroutineMap[goroutine.goroutineId] = goroutine
	g.mutex.Unlock()
	return goroutine
}

func (g *GoroutineManager) GoroutineCreatePn(goroutineName string,
	goroutineFunc func(Goroutine, ...interface{}),
	argFunc ...interface{}) uint64 {
	goroutine := g.constructGoroutine(goroutineName)
	go goroutineFunc(*goroutine, argFunc...)
	return goroutine.goroutineId
}

func (g *GoroutineManager) GoroutineCreateWithCancelPn(goroutineName string,
	goroutineFunc func(Goroutine, ...interface{}),
	argFunc ...interface{}) uint64 {
	goroutine := g.constructGoroutineWithCancel(goroutineName)
	go goroutineFunc(*goroutine, argFunc...)
	return goroutine.goroutineId
}

func (g *GoroutineManager) GoroutineCreateP0(goroutineName string,
	goroutineFunc func(Goroutine)) uint64 {
	goroutine := g.constructGoroutine(goroutineName)
	go goroutineFunc(*goroutine)
	return goroutine.goroutineId
}

func (g *GoroutineManager) GoroutineCreateWithCancelP0(goroutineName string,
	goroutineFunc func(Goroutine)) (uint64, context.CancelFunc) {
	goroutine := g.constructGoroutineWithCancel(goroutineName)
	go goroutineFunc(*goroutine)
	return goroutine.goroutineId, *goroutine.goroutineCancelFunc
}

func (g *GoroutineManager) GoroutineCreateP1(goroutineName string,
	goroutineFunc func(Goroutine, interface{}),
	argFunc1 interface{}) uint64 {
	goroutine := g.constructGoroutine(goroutineName)
	go goroutineFunc(*goroutine, argFunc1)
	return goroutine.goroutineId
}

func (g *GoroutineManager) GoroutineCreateWithCancelP1(goroutineName string,
	goroutineFunc func(Goroutine, interface{}),
	argFunc1 interface{}) (uint64, context.CancelFunc) {
	goroutine := g.constructGoroutineWithCancel(goroutineName)
	go goroutineFunc(*goroutine, argFunc1)
	return goroutine.goroutineId, *goroutine.goroutineCancelFunc
}

func (g *GoroutineManager) GoroutineCreateP2(goroutineName string,
	goroutineFunc func(Goroutine, interface{}, interface{}),
	argFunc1 interface{}, argFunc2 interface{}) uint64 {
	goroutine := g.constructGoroutine(goroutineName)
	go goroutineFunc(*goroutine, argFunc1, argFunc2)
	return goroutine.goroutineId
}

func (g *GoroutineManager) GoroutineCreateWithCancelP2(goroutineName string,
	goroutineFunc func(Goroutine, interface{}, interface{}),
	argFunc1 interface{}, argFunc2 interface{}) (uint64, context.CancelFunc) {
	goroutine := g.constructGoroutineWithCancel(goroutineName)
	go goroutineFunc(*goroutine, argFunc1, argFunc2)
	return goroutine.goroutineId, *goroutine.goroutineCancelFunc
}

func (g *GoroutineManager) GoroutineCreateP3(goroutineName string,
	goroutineFunc func(Goroutine, interface{}, interface{}, interface{}),
	argFunc1 interface{}, argFunc2 interface{}, argFunc3 interface{}) uint64 {
	goroutine := g.constructGoroutine(goroutineName)
	go goroutineFunc(*goroutine, argFunc1, argFunc2, argFunc3)
	return goroutine.goroutineId
}

func (g *GoroutineManager) GoroutineCreateWithCancelP3(goroutineName string,
	goroutineFunc func(Goroutine, interface{}, interface{}, interface{}),
	argFunc1 interface{}, argFunc2 interface{}, argFunc3 interface{}) (uint64, context.CancelFunc) {
	goroutine := g.constructGoroutineWithCancel(goroutineName)
	go goroutineFunc(*goroutine, argFunc1, argFunc2, argFunc3)
	return goroutine.goroutineId, *goroutine.goroutineCancelFunc
}

func (g *GoroutineManager) GoroutineCreateP4(goroutineName string,
	goroutineFunc func(Goroutine, interface{}, interface{}, interface{}, interface{}),
	argFunc1 interface{}, argFunc2 interface{}, argFunc3 interface{}, argFunc4 interface{}) uint64 {
	goroutine := g.constructGoroutine(goroutineName)
	go goroutineFunc(*goroutine, argFunc1, argFunc2, argFunc3, argFunc4)
	return goroutine.goroutineId
}

func (g *GoroutineManager) GoroutineCreateWithCancelP4(goroutineName string,
	goroutineFunc func(Goroutine, interface{}, interface{}, interface{}, interface{}),
	argFunc1 interface{}, argFunc2 interface{}, argFunc3 interface{}, argFunc4 interface{}) (uint64, context.CancelFunc) {
	goroutine := g.constructGoroutineWithCancel(goroutineName)
	go goroutineFunc(*goroutine, argFunc1, argFunc2, argFunc3, argFunc4)
	return goroutine.goroutineId, *goroutine.goroutineCancelFunc
}

func (g *GoroutineManager) GoroutineRemove(goroutine *Goroutine) {
	g.mutex.Lock()
	delete(g.goroutineMap, goroutine.goroutineId)
	g.mutex.Unlock()
}

func (g *GoroutineManager) Destroy() {
	g.mutex = nil
	g.counter = 0
	g.managerName = ""
	g.goroutineMap = make(map[uint64]*Goroutine)
}

func (g GoroutineManager) GetGoroutineById(goroutineId uint64) *Goroutine {
	return g.goroutineMap[goroutineId]
}

func (g GoroutineManager) GoroutineCount() int {
	g.mutex.RLock()
	length := len(g.goroutineMap)
	g.mutex.RUnlock()
	return length
}

func (g GoroutineManager) GoroutineDump() {
	g.mutex.RLock()
	fmt.Println("-------------GoroutineDump Begin-------------")
	fmt.Println("Goroutine_count", len(g.goroutineMap))
	for _, v := range g.goroutineMap {
		fmt.Println("goroutineId:", v.goroutineId, "|",
			"goroutineName:", v.goroutineName, "|",
			"goroutineCreateTime:", v.goroutineCreateTime.Unix())
	}
	fmt.Println("-------------GoroutineDump End-------------")
	g.mutex.RUnlock()
}
