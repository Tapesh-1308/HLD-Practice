package main

import (
	"errors"
	"fmt"
	"sync"
)

func add(a int, b int) int {
	return a + b
}

func stats(a int, b int) (int, int) {
	sum := a + b
	diff := a - b
	return sum, diff
}

type User struct {
	Name  string
	Email string
	Age   int
}

func (u User) Greet() {
	fmt.Printf("Hello, my name is %s and I am %d years old.\n", u.Name, u.Age)
}

type Animal interface {
	Speak() string
}

// implementing the Animal interface
type Dog struct{}

func (d Dog) Speak() string {
	return "Woof!"
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// defer
func deferredFunction() {
	defer fmt.Println("Deferred statement executed")
	fmt.Println("Function execution")
}

func main() {
	fmt.Println("Hello, World!")

	//  ways to declare variables in Go
	// var a int = 10;
	// b := 20

	if x > 10 {
		fmt.Println("x is greater than 10")
	} else {
		fmt.Println("x is less than or equal to 10")
	}

	for i := 0; i < 5; i++ {
		fmt.Println("Iteration:", i)
	}

	arr := [5]int{1, 2, 3, 4, 5}
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	slice = append(slice, 11, 12, 13)

	a := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	a["four"] = 4
	delete(a, "two")

	user := User{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   30,
	}

	user.Greet()

	x := 10
	ptr := &x

	*ptr = 20

	dog := Dog{}
	fmt.Println(dog.Speak())

	result, err := divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// go routines
	go func() {
		fmt.Println("This is a goroutine")
	}()

	// channels
	ch := make(chan string)

	go func() {
		ch <- "Hello from goroutine"
	}()

	msg := <-ch
	fmt.Println(msg)

	// wait group
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("Goroutine with WaitGroup")
	}()

	wg.Wait()
}
