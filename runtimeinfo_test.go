// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"fmt"
	"testing"
	"time"
)

func fib(n int64) int64 {
	if n < 2 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}
func TestRuntimeInfo(t *testing.T) {
	perf := &runtimePerf{}
	begin := time.Now()
	start := begin
	perf.Init()
	go func() {
		for i := 0; i < 50; i++ {
			time.Sleep(time.Second * 2)
			go func() {
				for i := 0; i < 100; i++ {
					time.Sleep(time.Second)
				}
			}()
		}
	}()
	startPoing := begin
	for now := time.Now(); now.Before(startPoing.Add(500 * time.Second)); now = time.Now() {
		time.Sleep(1 * time.Second)
		if !start.Add(20 * time.Second).After(now) {
			perf.Snap()
			begin = now
			start = start.Add(20 * time.Second)
			p := &runtimeBlock{}
			p.Read(perf)
			fmt.Println("NumGoroutine     :", p.NumGoroutine.IntSlice())
			fmt.Println("NumCgoCall       :", p.NumCgoCall.IntSlice())
			fmt.Println("NumGC            :", p.NumGC.IntSlice())
			fmt.Println("Frees            :", p.Frees.IntSlice())
			fmt.Println("Mallocs          :", p.Mallocs.IntSlice())
			fmt.Println("Lookups          :", p.Lookups.IntSlice())
			fmt.Println("PauseTotalNs     :", p.PauseTotalNs.IntSlice())
			fmt.Println("GCTime           :", p.GCTime.IntSlice())

			fmt.Println("MemTotalSys      :", p.MemTotalSys.FloatSlice())
			fmt.Println("MemHeapSys       :", p.MemHeapSys.FloatSlice())
			fmt.Println("MemStackSys      :", p.MemStackSys.FloatSlice())
			fmt.Println("MSpanSys         :", p.MSpanSys.IntSlice())
			fmt.Println("MSpanInuse       :", p.MSpanInuse.IntSlice())
			fmt.Println("MCacheSys        :", p.MCacheSys.IntSlice())
			fmt.Println("MCacheInuse      :", p.MCacheInuse.IntSlice())
			fmt.Println("BuckHashSys      :", p.BuckHashSys.IntSlice())
			fmt.Println("HeapInuse        :", p.HeapInuse.FloatSlice())
			fmt.Println("StackInuse       :", p.StackInuse.FloatSlice())
			fmt.Println("UserTime         :", p.UserTime.IntSlice())
			fmt.Println("UserUtilization  :", p.UserUtilization.IntSlice())
			fmt.Println("mem              :", p.mem.IntSlice())
			perf.Reset()
		}
		fmt.Printf("---------------\nfib(%d)=%d\n---------\n", 45, fib(45))
	}
	fmt.Printf("done\n")

}
