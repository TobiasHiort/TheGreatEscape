package main

import (
	"fmt"
//	"github.com/beefsack/go-astar"
	
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

/*
func setPlan(p *Person, m Map) { // implement A*
	//You shall not pass!
}
*/

func kill(p *Person) {
	p.alive = false
}

func main() {
//	pq := PriorityQueue{}

	
//	var p1 = Person{}
/*	var p1 = makePerson()
	
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
	fmt.Println(p1)*/

	t1 := Tile{0, 0}
	t2 := Tile{0, 1}
	t3 := Tile{1, 0}
	t4 := Tile{1, 1}

	
	m := []Tile{t1, t2, t3, t4}

/*	s1 := Thing{t1,1}
	var q = Queue{s1}

	fmt.Println(q)
	(&q).Add(t2, 2)
	fmt.Println(q)

	fmt.Println((&q).Pop()) */
	
	fmt.Println(getNeighbours(m, t1))
	
	path := getpath(m, t1, t4)
	fmt.Println(path)
	fmt.Println("\nThey shall all burn!");

}

