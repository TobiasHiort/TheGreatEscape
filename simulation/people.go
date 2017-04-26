package main

import (
	"fmt"
)

type Person struct {
	alive bool
	safe  bool
	hp    float32
	path  []*tile
	plan  []*tile
}

func makePerson(t *tile) *Person {
	var person = Person{}
	person.alive = true
	person.path = append(person.path, t)
	person.hp = 100 // TODO: default health??
	t.occupied = &person
	return &person
}

func (p *Person) updateStats() {
	currentTile := p.path[len((p.path))-1]
	(p.path[len(p.path)-1]).occupied = p
	if len(p.path) > 1 {
		p.path[len(p.path)-2].occupied = nil
	}
	p.hp = p.hp - currentTile.getDamage()
	if p.hp <= 0 {
		p.kill()
	}
}

func (t *tile) getDamage() float32 {
	damage := float32(0)
	damage = 100 * float32(t.fireLevel) // man dör om man kliver i elden, right...
	damage = damage + float32(t.heat)   // TODO: how much does the fire hurt??
	// damage = damage + effect from smoke'n stuff
	return damage
}

func (p *Person) moveTo(t *tile) bool {
	if validTile(t) && t.occupied == nil {
		p.path = append(p.path, t)
		p.updateStats()
		return true
	} else {
		return false
	}
}

func (p *Person) followPlan() {
	if len(p.plan) > 0 {
		if p.moveTo(p.plan[0]) {
			p.plan = p.plan[1:]
		}
	} else if p.path[len(p.path)-1].door {
		(p.path[len(p.path)-1].occupied) = nil
		p.path = append(p.path, nil) // replace with safezone?
		p.save()
	} else {
		fmt.Println("you're screwed!")
		p.kill()
		// TODO: no valid path! panic behavior? lay down and w8 for death?
	}
}

func (p *Person) kill() {
	p.alive = false
}

func (p *Person) save() {
	p.safe = true
	// TODO: maybe p.movetosafezone?
}

func (p *Person) updatePlan(m *[][]tile) {
	plan, ok := getPath(m, p.path[len(p.path)-1])
	if ok {
		p.plan = plan[1:]
	}
}

func (p *Person) MovePerson(m *[][]tile) {
	if p.safe || !p.alive {
		return
	}
	p.updatePlan(m)
	p.followPlan()
}

func MainPeople() {

	matrix := [][]int{
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{1, 0, 1, 1, 1, 1, 1},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 2, 0, 0, 0}}
	testmap := TileConvert(matrix)

	start1 := &testmap[1][0]
	start2 := &testmap[1][2]
	var p1 = *makePerson(start1)
	var p2 = *makePerson(start2)

	for !p1.safe && p1.alive || !p2.safe && p2.alive {
		if !p1.safe {
			fmt.Println("p1:", p1.path[len(p1.path)-1])
			p1.MovePerson(&testmap)
		}
		if !p2.safe {
			fmt.Println("p2:", p2.path[len(p2.path)-1])
			p2.MovePerson(&testmap)
		}
		fmt.Println("- - - - - - -")
	}
	fmt.Println("p1:", p1.path[len(p1.path)-1])
	fmt.Println("p2:", p2.path[len(p2.path)-1])
}
