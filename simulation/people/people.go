package main

import (
	"fmt"
	//"math/random"
)

// In this struct is all the data we might want to collect from a person.
// Implementing the stats like this might lead to some ugly code though. We'll se.
type person struct {
	ID int
	escapeTime int

	alive bool
	timeOfDeath int

	//startLocation tile
	//escapeDoor tile
}

func sayHello() {
	fmt.Println("Hello, I am ")
}

// Here I initiate an example person
func personInit(personID int) person {
	newPerson := person{personID, -1, true, -1}
	return newPerson
}

func act(dude person) {
	sayHello()
	fmt.Println("a person\n")
}

func main() {
	personAmount := 10
	personID := 0
	//slice is like an array, only a bit different
	personSlice := []person{}
	for i := 0; i < personAmount; i++ {
		personSlice = append (personSlice, personInit(personID))
		personID++
	}

	for _, person := range personSlice {
		go act(person)
	}
}
