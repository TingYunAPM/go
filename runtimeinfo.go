// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"runtime"
)

type runtimePerf struct {
	NumGoroutine structPerformance
	NumCgoCall   structPerformance
	PauseTotalNs structPerformance
	NumGC        structPerformance
	GCTime       structPerformance //1分钟
	MemTotalSys  structPerformance
	MemStackSys  structPerformance
	MSpanSys     structPerformance
	MCacheSys    structPerformance
	MemHeapSys   structPerformance
	BuckHashSys  structPerformance
	Frees        structPerformance
	Mallocs      structPerformance
	Lookups      structPerformance
	HeapInuse    structPerformance
	StackInuse   structPerformance
	MSpanInuse   structPerformance
	MCacheInuse  structPerformance
	//for this process
	FDSize     structPerformance
	Threads    structPerformance
	VMPeakSize structPerformance
	RssPeak    structPerformance

	numCgoCall float64
	memState   *runtime.MemStats

	sys sysInfo

	UserTime        structPerformance
	UserUtilization structPerformance
	mem             structPerformance

	//for this process
	fdSize     uint64
	vmPeakSize uint64
	vmCurrent  uint64
	rssPeak    uint64
	rssCurrent uint64
}

func (p *runtimePerf) Reset() {
	p.NumGoroutine.Reset()
	p.NumCgoCall.Reset()
	p.PauseTotalNs.Reset()
	p.NumGC.Reset()
	p.GCTime.Reset()
	p.MemTotalSys.Reset()
	p.MemStackSys.Reset()
	p.MSpanSys.Reset()
	p.MCacheSys.Reset()
	p.MemHeapSys.Reset()
	p.BuckHashSys.Reset()
	p.Frees.Reset()
	p.Mallocs.Reset()
	p.Lookups.Reset()
	p.HeapInuse.Reset()
	p.StackInuse.Reset()
	p.MSpanInuse.Reset()
	p.MCacheInuse.Reset()
	//for this process
	p.FDSize.Reset()
	p.Threads.Reset()
	p.VMPeakSize.Reset()
	p.RssPeak.Reset()

	p.UserTime.Reset()
	p.mem.Reset()
	p.UserUtilization.Reset()
}

func (p *runtimePerf) Init() {
	p.numCgoCall = float64(runtime.NumCgoCall())
	p.memState = &runtime.MemStats{}
	runtime.ReadMemStats(p.memState)
	p.NumGoroutine.AddValue(float64(runtime.NumGoroutine()), 0)
	p.MemTotalSys.AddValue(float64(p.memState.Sys)/1048576, 0)
	p.MemHeapSys.AddValue(float64(p.memState.HeapSys)/1048576, 0)
	p.MemStackSys.AddValue(float64(p.memState.StackSys)/1048576, 0)
	p.MSpanSys.AddValue(float64(p.memState.MSpanSys), 0)
	p.MSpanInuse.AddValue(float64(p.memState.MSpanInuse), 0)
	p.MCacheSys.AddValue(float64(p.memState.MCacheSys), 0)
	p.MCacheInuse.AddValue(float64(p.memState.MCacheInuse), 0)
	p.BuckHashSys.AddValue(float64(p.memState.BuckHashSys), 0)
	p.HeapInuse.AddValue(float64(p.memState.HeapInuse)/1048576, 0)
	p.StackInuse.AddValue(float64(p.memState.StackInuse)/1048576, 0)

	p.FDSize.Reset()
	p.Threads.Reset()
	p.VMPeakSize.Reset()
	p.RssPeak.Reset()

	p.UserTime.Reset()
	p.UserUtilization.Reset()
	p.mem.Reset()

	p.sys.Init()
}
func (p *runtimePerf) Snap() {

	numCgoCall := p.numCgoCall
	memState := p.memState
	lastSys := p.sys
	p.Init()
	if runtime.GOOS == "linux" {
		p.FDSize.AddValue(float64(p.sys.FdSize), 0)
	}
	if p.sys.Threads > 0 {
		p.Threads.AddValue(float64(p.sys.Threads), 0)
	}
	if lastSys.err == nil && p.sys.err == nil {
		userUse := p.sys.cpuProcess.ProcessUse() - lastSys.cpuProcess.ProcessUse()
		fullUse := p.sys.cpuSystem.FullUse() - lastSys.cpuSystem.FullUse()
		UserTime := float64(userUse) / 100
		p.UserTime.AddValue(UserTime, 0)
		UserUtilization := float64(userUse) * 100 / float64(fullUse)
		p.UserUtilization.AddValue(UserUtilization, 0)
		vmRss := float64(p.sys.vmRss) / 1024
		p.mem.AddValue(vmRss, 0)
	}

	p.NumCgoCall.AddValue(p.numCgoCall-numCgoCall, 0)
	p.Frees.AddValue(float64(p.memState.Frees-memState.Frees), 0)
	p.Mallocs.AddValue(float64(p.memState.Mallocs-memState.Mallocs), 0)
	lookups := p.memState.Lookups - memState.Lookups
	p.Lookups.AddValue(float64(lookups), 0)

	numGC := p.memState.NumGC - memState.NumGC
	p.PauseTotalNs.AddValue(float64(p.memState.PauseTotalNs-memState.PauseTotalNs)/1000000, 0)
	p.NumGC.AddValue(float64(p.memState.NumGC-memState.NumGC), 0)
	loopCount := numGC
	cacheCount := uint32(len(memState.PauseEnd))
	if loopCount > cacheCount {
		loopCount = cacheCount
	}
	index := ((numGC % cacheCount) + cacheCount - 1)
	for i := uint32(0); i < loopCount; i, index = i+1, index-1 {
		p.GCTime.AddValue(float64(p.memState.PauseNs[index%cacheCount])/1000000, 0)
	}
}

type runtimeBlock struct {
	NumGoroutine structPerformance
	NumCgoCall   structPerformance
	PauseTotalNs structPerformance
	NumGC        structPerformance
	GCTime       structPerformance //1分钟
	MemTotalSys  structPerformance
	MemStackSys  structPerformance
	MSpanSys     structPerformance
	MCacheSys    structPerformance
	MemHeapSys   structPerformance
	BuckHashSys  structPerformance
	Frees        structPerformance
	Mallocs      structPerformance
	Lookups      structPerformance
	HeapInuse    structPerformance
	StackInuse   structPerformance
	MSpanInuse   structPerformance
	MCacheInuse  structPerformance
	//for this process
	UserTime        structPerformance
	UserUtilization structPerformance
	mem             structPerformance
	FDSize          structPerformance
	Threads         structPerformance
	VMPeakSize      structPerformance
	RssPeak         structPerformance
}

func (p *runtimeBlock) Merge(q *runtimeBlock) {
	p.NumGoroutine.Append(&q.NumGoroutine)
	p.NumCgoCall.Append(&q.NumCgoCall)
	p.PauseTotalNs.Append(&q.PauseTotalNs)
	p.NumGC.Append(&q.NumGC)
	p.GCTime.Append(&q.GCTime)
	p.MemTotalSys.Append(&q.MemTotalSys)
	p.MemStackSys.Append(&q.MemStackSys)
	p.MSpanSys.Append(&q.MSpanSys)
	p.MCacheSys.Append(&q.MCacheSys)
	p.MemHeapSys.Append(&q.MemHeapSys)
	p.BuckHashSys.Append(&q.BuckHashSys)
	p.Frees.Append(&q.Frees)
	p.Mallocs.Append(&q.Mallocs)
	p.Lookups.Append(&q.Lookups)
	p.HeapInuse.Append(&q.HeapInuse)
	p.StackInuse.Append(&q.StackInuse)
	p.MSpanInuse.Append(&q.MSpanInuse)
	p.MCacheInuse.Append(&q.MCacheInuse)

	p.UserTime.Append(&q.UserTime)
	p.UserUtilization.Append(&q.UserUtilization)
	p.mem.Append(&q.mem)
	p.FDSize.Append(&q.FDSize)
	p.Threads.Append(&q.Threads)
	p.VMPeakSize.Append(&q.VMPeakSize)
	p.RssPeak.Append(&q.RssPeak)
}

func (p *runtimeBlock) Read(r *runtimePerf) {
	p.NumGoroutine.Reset().Append(&r.NumGoroutine)
	p.NumCgoCall.Reset().Append(&r.NumCgoCall)
	p.PauseTotalNs.Reset().Append(&r.PauseTotalNs)
	p.NumGC.Reset().Append(&r.NumGC)
	p.GCTime.Reset().Append(&r.GCTime)
	p.MemTotalSys.Reset().Append(&r.MemTotalSys)
	p.MemStackSys.Reset().Append(&r.MemStackSys)
	p.MSpanSys.Reset().Append(&r.MSpanSys)
	p.MCacheSys.Reset().Append(&r.MCacheSys)
	p.MemHeapSys.Reset().Append(&r.MemHeapSys)
	p.BuckHashSys.Reset().Append(&r.BuckHashSys)
	p.Frees.Reset().Append(&r.Frees)
	p.Mallocs.Reset().Append(&r.Mallocs)
	p.Lookups.Reset().Append(&r.Lookups)
	p.HeapInuse.Reset().Append(&r.HeapInuse)
	p.StackInuse.Reset().Append(&r.StackInuse)
	p.MSpanInuse.Reset().Append(&r.MSpanInuse)
	p.MCacheInuse.Reset().Append(&r.MCacheInuse)

	p.UserTime.Reset().Append(&r.UserTime)
	p.UserUtilization.Reset().Append(&r.UserUtilization)
	p.mem.Reset().Append(&r.mem)

	p.FDSize.Reset().Append(&r.FDSize)
	p.Threads.Reset().Append(&r.Threads)
	p.VMPeakSize.Reset().Append(&r.VMPeakSize)
	p.RssPeak.Reset().Append(&r.RssPeak)
}
