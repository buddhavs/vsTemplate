package bdd

import (
	"fmt"
	"os"
	"testing"
)

// https://golang.org/pkg/testing
// https://github.com/smartystreets/gunit
// https://godoc.org/github.com/smartystreets/gunit

func ExamplePerm() {
	for _, value := range []int{1, 2, 4, 3, 0} {
		fmt.Println(value)
	}
	// Unordered output: 4
	// 2
	// 1
	// 3
	// 0
}

func ExampleSalutations() {
	fmt.Println("hello, and")
	fmt.Println("goodbye")
	// Output:
	// hello, and
	// goodbye
}

func TestFoo(t *testing.T) {
	// <setup code>
	t.Run("A=1", func(t *testing.T) {})
	t.Run("A=2", func(t *testing.T) {})
	t.Run("B=1", func(t *testing.T) {})
	// <tear-down code>
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
