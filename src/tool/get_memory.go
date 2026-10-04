// package memmonitor 提供了获取 Windows 进程内存占用的封装。
package tool

import (
	"fmt"
	"math"
	"runtime"
	"syscall"
	"unsafe"
)

// Win32 API 结构体，不对外暴露 (私有)
type processMemoryCounters struct {
	cb                         uint32
	pageFaultCount             uint32
	peakWorkingSetSize         uintptr
	workingSetSize             uintptr
	quotaPeakPagedPoolUsage    uintptr
	quotaPagedPoolUsage        uintptr
	quotaPeakNonPagedPoolUsage uintptr
	quotaNonPagedPoolUsage     uintptr
	pagefileUsage              uintptr
	peakPagefileUsage          uintptr
}

// Monitor 内存监控对象 (对外暴露)
type Monitor struct {
	pid                      uint32
	handle                   uintptr
	procGetProcessMemoryInfo *syscall.LazyProc
}

// New 创建并初始化内存监控实例
func NewMemoryMonitor() (*Monitor, error) {
	m := &Monitor{}
	if err := m.init(); err != nil {
		return nil, err
	}
	return m, nil
}

// init 初始化内部 DLL 句柄
func (m *Monitor) init() error {
	modKernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetCurrentProcess := modKernel32.NewProc("GetCurrentProcess")
	procGetCurrentProcessId := modKernel32.NewProc("GetCurrentProcessId")

	modPsapi := syscall.NewLazyDLL("psapi.dll")
	m.procGetProcessMemoryInfo = modPsapi.NewProc("GetProcessMemoryInfo")

	pidRet, _, _ := procGetCurrentProcessId.Call()
	m.pid = uint32(pidRet)

	handleRet, _, _ := procGetCurrentProcess.Call()
	m.handle = handleRet

	return nil
}

// GetMem 获取当前进程的各项内存指标（单位：MB，保留 2 位小数）
//
// 返回值说明：
//   - workingSetMB: 任务管理器显示的“内存(物理工作集)”，包含 Go + Cgo/C++ + DLL 等所有原生内存
//   - privateMB   : 进程私有提交内存 (Private Commit Size)，用于分析进程申请的总虚存/内存泄漏
//   - goAllocMB   : Go 运行时正在使用的堆内存大小 (m.Alloc)
//   - goSysMB     : Go 运行时向操作系统申请的总内存大小 (m.Sys)
//   - pid         : 当前进程 PID
func (m *Monitor) GetMem() (workingSetMB, privateMB, goAllocMB, goSysMB float64, pid uint32, err error) {
	// 1. 获取 Go 内部内存
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	// 2. 获取 Windows 系统总进程内存
	var counters processMemoryCounters
	counters.cb = uint32(unsafe.Sizeof(counters))

	ret, _, errNo := m.procGetProcessMemoryInfo.Call(
		m.handle,
		uintptr(unsafe.Pointer(&counters)),
		uintptr(counters.cb),
	)

	if ret == 0 {
		return 0, 0, 0, 0, m.pid, fmt.Errorf("GetProcessMemoryInfo 失败: %v", errNo)
	}

	// 3. 计算并精确截断为 2 位小数
	workingSetMB = round2(float64(counters.workingSetSize) / 1024 / 1024)
	privateMB = round2(float64(counters.pagefileUsage) / 1024 / 1024)
	goAllocMB = round2(float64(ms.Alloc) / 1024 / 1024)
	goSysMB = round2(float64(ms.Sys) / 1024 / 1024)

	return workingSetMB, privateMB, goAllocMB, goSysMB, m.pid, nil
}

// 保留 2 位小数辅助函数
func round2(val float64) float64 {
	return math.Round(val*100) / 100
}
