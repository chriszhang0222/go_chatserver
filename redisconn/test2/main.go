package main

import (
	"sync"
	"time"
)

type Glimit struct {
	n int
	c chan struct{}
}

// initialization Glimit struct
func New(n int) *Glimit {
	return &Glimit{
		n: n,
		c: make(chan struct{}, n),
	}
}
// Run f in a new goroutine but with limit.
func (g *Glimit) Run(f func()) {
	g.c <- struct{}{}
	go func() {
		f()
		<-g.c
	}()
}

var wg = sync.WaitGroup{}
var c = make(chan int, 1)
func main() {
	for i :=0 ;i < 10 ;i++ {
		go func() {
			time.Sleep(3 * time.Second)
			c <- 2
		}()
	}
	<-c

	//number := 10
	//g := New(1)
	//for i := 0; i < number; i++ {
	//	wg.Add(1)
	//	value :=i
	//	goFunc := func() {
	//		// 做一些业务逻辑处理
	//		fmt.Printf("go func: %d\n", value)
	//		time.Sleep(time.Second)
	//		wg.Done()
	//	}
	//	g.Run(goFunc)
	//}
	//wg.Wait()
}
