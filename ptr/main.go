package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

type Person struct {
	id   int
	name string
	age  int
}

func main() {
	num := 10
	numCpy := num
	fmt.Printf("Num: %d\nNum Copy: %d\n\n", num, numCpy)
	numCpy = 15
	fmt.Printf("Num after Num Copy Change: %d\nNum Copy after Change: %d\n\n", num, numCpy)
	var numPtr *int = &num
	numPtrAuto := &num
	fmt.Printf("Explicit ptr init %d, Implicit ptr init: %d\n\n", *numPtr, *numPtrAuto)

	*numPtrAuto = 20
	fmt.Printf("Num after Ptr Change: %d, Num Copy after Ptr Change: %d, Num Ptr Imp after Ptr Change: %d\n", num, numCpy, *numPtrAuto)
	fmt.Printf("Address of Num: %x, Address of Num Ptr: %x, Address of Num Copy: %x\n", &num, numPtr, &numCpy)

	color.Yellow("Please enter a number of iterations")
	input, err := fmt.Scanln(numPtr)
	if input < 1 || err != nil {
		fmt.Printf("Inputs scanned %d, Error found %s\n", input, err.Error())
		panic(1)
	}
	rand.Seed(time.Now().UnixNano())
	end := rand.Intn(1000000)
	fmt.Println("Start: ", end)
	iterate(&end, numPtr)
	fmt.Println("End: ", end)
	mapFunc()
	sliceFunc()
	structFunc()
	staticFunc(num)
	color.Red("Done!")
}

func iterate(start, end *int) {
	for i := 0; i < *end; i++ {
		(*start)++
	}
}

func mapFunc() {
	mp := map[string]int{"foo": 1}
	updateMap(mp)
	fmt.Println(mp)
}

func updateMap(mp map[string]int) {
	mp["foo"] = 4
	mp["baz"] = 5
}

func sliceFunc() {
	sl := []int{1, 2, 4, 5}
	updateSlice(sl)
	fmt.Println(sl)
	sli := make([]int, 4)
	sli = append(sli, 2)
	updateSlice(sli)
	fmt.Println(sli)
}

func updateSlice(sl []int) []int {
	// slices are pointer windows into an array basically
	// when we pass in our slice
	for i := 0; i < len(sl); i++ {
		if sl[i] == 0 {
			break
		}
		sl[i] = sl[i] + 1
	}

	sl = append(sl, 7)
	return sl
}

func structFunc() {
	var p Person
	p.age = 33
	p.name = "Christian"
	p2 := Person{name: "Bob", id: 1, age: 0}
	fmt.Println(p, p2)
	p2Ptr := &p2
	p2Cpy := p2

	p2Cpy.age = 12
	// no change, need to pass by reference
	incAge(p2Cpy)
	p2Cpy.name = "Copy"
	p2Ptr.age = 30
	fmt.Println(p2Cpy, *p2Ptr, p2)
}

func incAge(p Person) {
	p.age += 1
}

// always needs var (type optional)
var number = 14

func staticFunc(a int) {
	fmt.Println("Static number", number)
	fmt.Println("Scoped func param", a)
	//fmt.Println("Inaccesible field", num)
}
