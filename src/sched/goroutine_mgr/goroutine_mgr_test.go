package goroutine_mgr

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func printn(g Goroutine, a ...interface{}) {
	defer g.OnQuit()
	fmt.Println("print", fmt.Sprintf("%d %d %d %d %d", a...))
	time.Sleep(100 * time.Microsecond)
}

func print0(g Goroutine) {
	defer g.OnQuit()
	fmt.Println("print0")
	time.Sleep(100 * time.Microsecond)
}

func print1(g Goroutine, a interface{}) {
	defer g.OnQuit()
	fmt.Println("print1", a)
	time.Sleep(100 * time.Microsecond)
}

func print2(g Goroutine, a interface{}, b interface{}) {
	defer g.OnQuit()
	fmt.Println("print2", a, b)
	time.Sleep(100 * time.Microsecond)
}

func print3(g Goroutine, a interface{}, b interface{}, c interface{}) {
	defer g.OnQuit()
	fmt.Println("print3", a, b, c)
	time.Sleep(100 * time.Microsecond)
}

func print4(g Goroutine, a interface{}, b interface{}, c interface{}, d interface{}) {
	defer g.OnQuit()
	fmt.Println("print4", a, b, c, d)
	time.Sleep(100 * time.Microsecond)
}

func TestAll(t *testing.T) {
	manager := new(GoroutineManager)
	manager.Initialise("mgr1")

	gid := manager.GoroutineCreatePn("goroutine"+strconv.Itoa(0), printn, 1, 2, 3, 4, 5)
	fmt.Println("gid:", gid)

	gid0 := manager.GoroutineCreateP0("goroutine"+strconv.Itoa(0), print0)
	fmt.Println("gid0:", gid0)

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

func select0(g Goroutine) {
	defer g.OnQuit()
	fmt.Println("select0")

	select {
	case <-(*g.GetContext()).Done():
		fmt.Println("done/cancelled")
	case <-time.After(3 * time.Second):
		fmt.Println("timeout")
	}
}

func TestCancel(t *testing.T) {
	manager := new(GoroutineManager)
	manager.Initialise("mgr2")

	gid0, cbCancel := manager.GoroutineCreateWithCancelP0("goroutine"+strconv.Itoa(0), select0)
	fmt.Println("gid0:", gid0)

	time.Sleep(1 * time.Second)
	cbCancel()
}
