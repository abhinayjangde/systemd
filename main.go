package main

import (
	"fmt"
	"sync"
)

type Singleton struct{}

var (
	instance *Singleton
	once     sync.Once
)

func GetSingleton() *Singleton {
	once.Do(func() {
		instance = &Singleton{}
	})

	return instance
}

func main() {
	s1 := GetSingleton()
	s2 := GetSingleton()

	fmt.Println(s1 == s2)
}
