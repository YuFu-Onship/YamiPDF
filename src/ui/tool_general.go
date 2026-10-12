package ui

import (
	"fmt"
	"image"
	"os"
	"runtime"
)

// 加载图片
func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// 得到内存占用情况
func getGoMemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Alloc: 当前仍在使用的堆内存大小
	// Sys: 从操作系统申请的总内存大小
	allocMB := float64(m.Alloc) / 1024 / 1024
	sysMB := float64(m.Sys) / 1024 / 1024

	return fmt.Sprintf("Alloc: %.2f MB | Sys: %.2f MB\n", allocMB, sysMB)
}

// clamp1 limits v to range [0..1].
func clamp1(v float32) float32 {
	if v >= 1 {
		return 1
	} else if v <= 0 {
		return 0
	} else {
		return v
	}
}

// 泛型三元
func Ternary[T any](cond bool, trueVal, falseVal T) T {
	// 泛型中, 强求这的两个参数是同样类型
	if cond {
		return trueVal
	}
	return falseVal
}
