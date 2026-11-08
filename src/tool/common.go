package tool

// 泛型三元
func Ternary[T any](cond bool, trueVal, falseVal T) T {
	// 泛型中, 强求这的两个参数是同样类型
	if cond {
		return trueVal
	}
	return falseVal
}

// 惰性泛型三元(避免提前取值)
func TernaryLazy[T any](cond bool, trueVal func() T, falseVal func() T) T {
	if cond {
		return trueVal()
	}
	return falseVal()
}
