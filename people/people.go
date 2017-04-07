package main

import (
	"fmt"
)

type Tile struct {
	X int
	Y int
}

type Person struct {
	alive bool
	path []Tile
	plan []Tile	
}

func makePerson() Person{
	var person = Person{}
	person.alive = true
	return person
}

func moveTo(p *Person, t Tile) bool{
	p.path = append(p.path, t)
	return true
}

func followPlan(p *Person) {
	moveTo(p, p.plan[0])
	p.plan = p.plan[1:]
}

func kill(p *Person) {
	p.alive = false
}

func main() {
//	var p1 = Person{}
	var p1 = makePerson()
	
	fmt.Println(p1)
	t1 := Tile{1,1}
	t2 := Tile{2,2}
	
	moveTo(&p1, t1)
	fmt.Println(p1)
	p1.plan = append(p1.plan, t2)
	fmt.Println(p1)

	followPlan(&p1)
	fmt.Println(p1)

	kill(&p1)
	fmt.Println(p1)
	
	fmt.Println("\nThey shall all burn!");
}

