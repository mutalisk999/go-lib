package goroutine_mgr

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func print(g Goroutine, a ...interface{}) {
	defer g.OnQuit()
	fmt.Println("print", a)
	time.Sleep(100000)
}

func print1(g Goroutine, a interface{}) {
	defer g.OnQuit()
	fmt.Println("print1", a)
	time.Sleep(100000)
}

func print2(g Goroutine, a interface{}, b interface{}) {
	defer g.OnQuit()
	fmt.Println("print2", a, b)
	time.Sleep(100000)
}

func print3(g Goroutine, a interface{}, b interface{}, c interface{}) {
	defer g.OnQuit()
	fmt.Println("print3", a, b, c)
	time.Sleep(100000)
}

func print4(g Goroutine, a interface{}, b interface{}, c interface{}, d interface{}) {
	defer g.OnQuit()
	fmt.Println("print4", a, b, c, d)
	time.Sleep(100000)
}

func TestAll(t *testing.T) {
	manager := new(GoroutineManager)
	manager.Initialise("mgr1")

	gid := manager.GoroutineCreatePn("goroutine"+strconv.Itoa(0), print, 1, 2, 3, 4, 5)
	fmt.Println("gid:", gid)

	gid1 := manager.GoroutineCreateP1("goroutine"+strconv.Itoa(0), print1, 1)
	fmt.Println("gid1:", gid1)

	gid2 := manager.GoroutineCreateP2("goroutine"+strconv.Itoa(0), print2, 1, 2)
	fmt.Println("gid2:", gid2)

	gid3 := manager.GoroutineCreateP3("goroutine"+strconv.Itoa(0), print3, 1, 2, 3)
	fmt.Println("gid3:", gid3)

	gid4 := manager.GoroutineCreateP4("goroutine"+strconv.Itoa(0), print4, 1, 2, 3, 4)
	fmt.Println("gid4:", gid4)

	g := manager.GetGoroutineById(1)
	fmt.Println("g", g)

	manager.GoroutineDump()
	time.Sleep(10000000)
	fmt.Println("goroutineCount", manager.GoroutineCount())
}
