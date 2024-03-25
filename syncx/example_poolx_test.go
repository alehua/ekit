package syncx

import "fmt"

func ExampleNew() {
	type A struct {
		Name string
	}
	p := NewPool[A](func() A {
		return A{"A"}
	})

	res := p.Get()
	fmt.Print(res.Name)
	// Output:
	// A
}
