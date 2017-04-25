package main

import (
	"fmt"	
)

type Person struct {
	alive bool
	safe bool     
	hp float32
	path []*tile  
	plan []*tile 
}

func makePerson(t *tile) Person{
	var person = Person{}
	person.alive = true
	person.path = append(person.path, t)
	person.hp = 100   // TODO: default health??
	return person
}

func (p *Person)updateStats() {
	currentTile := p.path[len((p.path)) - 1]
	(p.path[len(p.path) - 2]).occupied = false
	currentTile.occupied = true
	p.hp = p.hp - currentTile.getDamage()
	if p.hp <= 0 {
		p.kill()		
	}
}

func (t *tile)getDamage() float32{
	damage := float32(0)
	damage = 100*float32(t.fireLevel)  // man dör om man kliver i elden, right...
	damage = damage + float32(t.heat)  // TODO: how much does the fire hurt??
	// damage = damage + effect from smoke'n stuff
	return damage
}

func (p *Person)moveTo(t *tile) bool{
	if validTile(t) && !t.occupied {
		p.path = append(p.path, t)
		p.updateStats()	
		return true
	} else {return false}
}

func (p *Person)followPlan() {
	if len(p.plan) > 0 {
		if p.moveTo(p.plan[0]) {
			p.plan = p.plan[1:]
		}
	} else /*(if p.reachedGoal())*/{ // TODO: empty plan can mean deadend!
		p.path[len(p.path) - 1].occupied = false
		p.path = append(p.path, nil)  // replace with safezone?
		p.save()
	}
}


func (p *Person)kill() {
	p.alive = false
}

func (p *Person)save() {
	p.safe = true
	// TODO: maybe p.movetosafezone?
}

func (p *Person) updatePath(m *[][]tile) {
	path, ok := getPath(m, p.path[len(p.path) - 1], p.plan[len(p.plan) - 1])
	if ok {
		p.path = path
	}
}

func MainPeople() {

	matrix := [][]int {
		{0,0,0,1,0,0,0},
		{0,0,0,1,0,0,0},
		{1,0,1,1,0,0,0},
		{0,0,0,1,0,0,0},
		{0,0,0,1,0,0,0},
		{0,0,0,0,0,0,0}, 
		{0,0,0,0,0,0,0}}
	testmap := TileConvert(matrix)

	start1 := &testmap[1][0]
	start2 := &testmap[1][2]
	var p1 = makePerson(start1)
	var p2 = makePerson(start2)

	goal := &testmap[6][3]

	plan1, _ := getPath(&testmap, start1, goal)
	plan2, _ := getPath(&testmap, start2, goal)

	p1.plan = plan1
	p2.plan = plan2

	fmt.Println(plan1[0])
	fmt.Println(len(plan1))
	
	for !p1.safe || !p2.safe{
		if !p1.safe {
			p1.followPlan() 
			fmt.Println("p1:", p1.path[len(p1.path) - 1])
		}
		if !p2.safe {
			p2.followPlan() 		 
			fmt.Println("p2:", p2.path[len(p2.path) - 1])
		}
		 fmt.Println("- - - - - - -")
	 }
}

