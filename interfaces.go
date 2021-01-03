package main

type Compare interface {
	// 小于传入item时返回负数，等于时返回0，大于时返回正数
	CompareWith(item Compare) int
}
