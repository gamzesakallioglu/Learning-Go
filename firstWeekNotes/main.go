package main

import "fmt"

func main() {

	// As usual, let's start with hello world
	fmt.Println("Hello World")

	// Declare a variable
	// var variable_name type
	var flag bool

	// Initialize a variable
	// 1. var variable_name type = value
	// 2. var variable_name = value
	// 3. variable := value
	var number int8 = 10
	var day = "sunday"
	numberOfSomething := 50

	//Also..
	var x, y = 80, 40.2
	var (
		a    int
		b    float64
		c    int    = 4
		d, e string // You can't say string,int..
		f    string = "Happy"
	)

	fmt.Println(flag, number, day, numberOfSomething, x, y, a, b, c, d, e, f)

	// Const
	// Constants are some variables that we will Never want to change the value. Just like value of Pi never changes..
	// Constants are initialized and the value never changes
	const pi = 3.14
	// const pi float32 = 3.14

	// ARRAYS
	// Decleare an array
	var grades [3]int // An array which lenght and capacity equals to 3.

	// Initialize an array
	var grades2 [3]int = [3]int{10, 20, 30}
	grades3 := [...]int{10, 20, 30}
	fmt.Println(len(grades), cap(grades))   // 3 3
	fmt.Println(len(grades2), cap(grades2)) // 3 3
	fmt.Println(len(grades3), cap(grades3)) // 3 3
	fmt.Println(grades2 == grades3)         // True

	// Slices
	// A slice of an array. Capacity changes if you append values.
	var slice1 []int         // Decleare
	slice2 := []int{1, 2, 3} //Initialize

	// Add values to a slice
	slice1 = append(slice1, 10, 20, 30)
	fmt.Println(slice1)      // [10 20 30]
	fmt.Println(len(slice1)) // 3
	fmt.Println(cap(slice1)) // 3

	fmt.Println(slice2)      // [1 2 3]
	fmt.Println(len(slice2)) // 3
	fmt.Println(cap(slice2)) // 3
	slice2 = append(slice2, 4, 5)
	fmt.Println(slice2)      // [1 2 3 4 5]
	fmt.Println(len(slice2)) // 5
	fmt.Println(cap(slice2)) // 6

	// You can append a slice into a slice
	slice1 = append(slice1, slice2...)
	fmt.Println(slice1) // [10 20 30 1 2 3 4 5]

	// Make - create a slice  (The elements will be all zeros)
	slice3 := make([]int, len(slice2))
	fmt.Println(slice3)                   // [0 0 0 0 0]
	fmt.Println(len(slice3), cap(slice3)) // 5 5

	fmt.Println(slice1)      // [10 20 30 1 2 3 4 5]
	fmt.Println(slice1[1])   // 20
	fmt.Println(slice1[1:4]) // [20 30 1]  (1st is included, 4th is excluded)
	fmt.Println(slice1[:])   // [10 20 30 1 2 3 4 5]  (This can look in vein, but its not. You can turn an array into a slice since it returns a slice.)

	// Copy - copy a slice
	slice4 := []int{80, 24, 60}
	slice5 := make([]int, len(slice4))
	copy(slice5, slice4)
	num := copy(slice5, slice4)
	fmt.Println(slice5) // [80 24 60]
	fmt.Println(num)    // 3

	// Maps
	var map1 map[string]int // key -> string, value -> int
	fmt.Println(map1)       // map[]
	map2 := map[int]int{}
	fmt.Println(map2) // map[]
	// Initialize a map
	dictionary := map[string]string{
		"a": "meaningOfa",
		"b": "meaningOfb", // dont forget about the last comma
	}
	fmt.Println(dictionary) // map[a:meaningOfa b:meaningOfb]

	// The values type might be a slice, but of course the key cannot.
	gradesOfStudents := map[string][]int{}
	gradesOfStudents["Ahmet"] = []int{50, 80, 100}
	gradesOfStudents["Ayşe"] = []int{70, 85, 90}
	fmt.Println(gradesOfStudents) // map[Ahmet:[50 80 100] Ayşe:[70 85 90]]

	// make can be used in maps too
	ages := make(map[int][]int, 10)
	fmt.Println(ages) // map[]

	// How to check if something is map's element?
	v, ok := gradesOfStudents["Ahmet"] // value, isExists
	fmt.Println(v, ok)                 // [50 80 100] true
	var v2, ok2 = gradesOfStudents["Mehmet"]
	fmt.Println(v2, ok2) // [] false

	// How to delete an element of a map
	delete(gradesOfStudents, "Ahmet")
	fmt.Println(gradesOfStudents) // map[Ayşe:[70 85 90]]

	// Once you iterate through a map, the order is unordered
	gradesOfStudents["Selen"] = []int{80, 60, 40}
	gradesOfStudents["Beyza"] = []int{40, 60, 100}
	fmt.Println(gradesOfStudents) // The order was Ayşe, Beyza, Selen in the console. Which Beyza should be last

	// Structs
	// Structs helps us to
	type Person struct {
		name string
		age  int
	}
	type Pet struct {
		name  string
		owner Person
		age   int
	} // The Pet's owner is a Person

	var gamze Person
	gamze.name = "Gamze"
	gamze.age = 22
	fmt.Println(gamze)

	merve := Person{}
	merve.name = "Merve"
	merve.age = 22
	fmt.Println(merve)

	serhat := Person{"Serhat", 30}
	fmt.Println(serhat)

	petito := Pet{name: "Petito", owner: gamze, age: 1}
	fmt.Println(petito) // {Petito {Gamze 22} 1}

	// Anonymus Articles
	var student struct {
		name string
		age  int
	}
	student.name = "Gamze"
	student.age = 17 // You access to student directly, rather than creating a variable

	// The other way
	teacher := struct {
		name string
		age  int
	}{
		"Gamze",
		30,
	}
	fmt.Println(teacher)

	// IF ELSE
	count := 0
	// some operation
	if count > 30 {
		fmt.Println("You've made it")
	} else if count > 15 {
		fmt.Println("You should try again")
	} else {
		fmt.Println("You lost")
	}

	// FOR
	// C-Style
	for j := 0; j <= 3; j++ {
		fmt.Println("Hello", j)
	}
	// Condition-Only
	i := 1
	for i <= 3 {
		fmt.Println("Hello", i)
		i *= 2
	}

	ages2 := []int{20, 3, 56, 10}
	for i, v := range ages2 {
		fmt.Println(i, v)
	}

}
