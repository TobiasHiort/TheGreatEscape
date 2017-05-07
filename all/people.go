package main

import (
	"fmt"
	"math"
	"time"
	"sync"
)

var step = float32(0)

type Person struct {
	alive bool
	safe  bool
	hp    float32
	path  []*tile
	plan  []*tile
	time float32
}

type Stats struct {
	x int
	y int
//	hp float32
}

func (p *Person)getStats() []int {
	//if p == nil || p.currentTile() == nil {return Stats{}}
	//if len(p.path) < 1 {return []int{0,0}}
//	return Stats{p.path[len(p.path) - 1].xCoord, p.path[len(p.path) - 1].yCoord}
//	return []int{p.currentTile().xCoord, p.currentTile().yCoord}//, p.hp}
	return []int{p.currentTile().yCoord, p.currentTile().xCoord}//, p.hp}
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
	currentTile := p.currentTile()
	(p.path[len(p.path)-1]).occupied = p
	if len(p.path) > 1 {
		if p.path[len(p.path) - 2] != currentTile {p.path[len(p.path)-2].occupied = nil}
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
	if p.path[len(p.path) - 1] == nil { return} // TODO updatestats
	if len(p.plan) > 0 { // follow tha plan!		
		if p.moveTo(p.plan[0]) {   // next step in plan is available -> move		
			p.plan = p.plan[1:]
			p.updateTime()  
		} else {                   // next step in plan is occupied -> w8
			p.wait()
			p.updateTime()
		}          
	} else if p.path[len(p.path) - 1].door {   // standing at the exit -> leave
		(p.path[len(p.path) - 1].occupied) = nil
		p.path = append(p.path, nil)  // replace with safezone?
		p.updateTime()
		p.save()
	} else {
		fmt.Println("you're screwed!")
		p.kill()
		// TODO: no valid path! panic behavior? lay down and w8 for death?
		// idea: don't update last plan-path, follow it despite fire etc?
	}
}

func (p *Person)wait() {
	p.path = append(p.path, p.path[len(p.path) - 1])
	p.updateStats()	
}

func (p *Person)IsWaiting() bool{
	if p == nil {return false}
	if len(p.path) <= 1 {
		return false
	} else {return p.path[len(p.path) - 1] == p.path[len(p.path) - 2]}
}

func (p *Person) kill() {
	p.alive = false
}

func (p *Person) save() {
	p.safe = true
	p.path[len(p.path) - 1] = &tile{}
	// TODO: maybe p.movetosafezone?
}

func (p *Person) updatePlan(m *[][]tile) {
	plan, ok := getPath(m, p.path[len(p.path)-1])
	if ok {
		p.plan = plan[1:]
	}
}

func (p *Person) MovePerson(m *[][]tile) {
	if p == nil {return}
	if p.safe || !p.alive {
		return
	}
	if p.time <= step {
		//if p.plan[0].occupied != nil || !validTile(p.plan[0]) {p.updatePlan(m)}
		p.updatePlan(m)
		p.followPlan()	
	}
}

func MovePeople(m *[][]tile, ppl []*Person) {
	var wg sync.WaitGroup
	
	for !CheckFinish(ppl) {
		wg.Add(len(ppl))
		print("\033[H\033[2J")
		PrintTileMapP(*m)
		fmt.Print("\n")
		time.Sleep(1000 * time.Millisecond)
		for _, pers := range ppl {			
			go func(p *Person){
				defer wg.Done()
				p.MovePerson(m)
			}(pers)
		}
		step++
		wg.Wait()
		FireSpread(*m)
	}
}

func (p *Person)currentTile() *tile{
	if len(p.path) == 0 {return nil}
	return p.path[len(p.path) - 1]
}

func (p *Person)updateTime() {
	if p.wasDiag() {//p.DiagonalStep() {
		p.time += float32(math.Sqrt(2))	
	} else {p.time += 1}
}

func (p *Person)DiagonalStep() bool{    // Is the next step a diagonal one?
	if p == nil {return false}
	if len(p.plan) < 1 {return false}
	return Diagonal (p.path[len(p.path) - 1], p.plan[0])
}

func (p *Person)wasDiag() bool{ // was the last step a diagonal one?
	if len(p.path) < 2 {return false}
	return Diagonal(p.path[len(p.path) - 1], p.path[len(p.path) - 2])
}

func Diagonal(t1, t2 *tile) bool {
	if t1 == nil {return false}
	if t2 == nil {return false}
	if t1.neighborNW == t2 {return true}
	if t1.neighborNE == t2 {return true}
	if t1.neighborSE == t2 {return true}
	if t1.neighborSW == t2 {return true}
	return false
}

func MainPeople() {

	matrix := [][]int{
		{0,0,0,0},
		{0,0,0,0},
		{0,0,0,0},
		{0,0,2,0}}
	testmap := TileConvert(matrix)
	
	start1 := &testmap[0][1]
	start2 := &testmap[2][0]
	start3 := &testmap[1][3]
	var p1 = *makePerson(start1)
	var p2 = *makePerson(start2)
	var p3 = *makePerson(start3)

	for !p1.safe || !p2.safe || !p3.alive  {
		if !p1.safe {
			fmt.Println("p1:", p1.path[len(p1.path)-1])
			p1.MovePerson(&testmap)
		}
		if !p2.safe {
			fmt.Println("p2:", p2.path[len(p2.path)-1])
			p2.MovePerson(&testmap)
		}
		if !p3.safe {
			fmt.Println("p3:", p3.path[len(p3.path)-1])
			p3.MovePerson(&testmap)
		}
		
		fmt.Println("- - - - - - -")
	}
	fmt.Println("p1:", p1.path[len(p1.path) - 1])
	fmt.Println("p2:", p2.path[len(p2.path) - 1])
	fmt.Println("- - - - - - -")
	fmt.Println("- - - - - - -")
	fmt.Println("p1")
	printPath(p1.path)
	fmt.Println("p2")
	printPath(p2.path)
}
