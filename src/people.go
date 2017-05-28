package main

import (
//	"fmt"
	"math"
//	"math/rand"
	"time"
	"sync"
)

// ariable to keep track of the current 'frame',
// used to decide wether a person is allowed to moved yet or not
var step = float32(0)

type Person struct {
	alive bool
	safe  bool
	screwed bool
	hp    int
	smokedmg int
	path  []*tile
	plan  []*tile
	dir Direction
	time float32
}

// returns the stats for 1 person,
// used to pipe the correct data to gui.
func (p *Person)getStats(aslice *[]int) {
	*aslice = append(*aslice, p.currentTile().yCoord)
	*aslice = append(*aslice, p.currentTile().xCoord)
	*aslice = append(*aslice, p.hp)
}

// Used to get the starting positions for all ppl
func StartStats(ppl []*Person) [][]int{
	lst := [][]int{}
	for _, pers := range ppl {
		templst := []int{}
		pers.getStats(&templst)
		lst = append(lst, templst)
	}
	return lst
}

func makePerson(t *tile) *Person {
	var person = Person{}
	person.alive = true
	person.screwed = false
	person.path = append(person.path, t)
	person.hp = 100 
	t.occupied = &person
	return &person
}

func (p *Person) updateStats() {
	currentTile := p.currentTile()
	currentTile.occupied = p
//	if len(p.path) > 1 {
//		if p.path[len(p.path) - 2] != currentTile {p.path[len(p.path)-2].occupied = nil}
//	}
	p.hp = p.hp - currentTile.getDamage()
	if p.hp <= 0 {
		p.kill()
	}	
}

// calculates and returns the amount of damage a tile deals
func (t *tile) getDamage() int {
	damage := int(0)
	damage += 10*int(t.fireLevel) 
	if t.smoke > 10 {
		damage += 1
		t.occupied.smokedmg += 1
	}
	return damage
}

// attempts to move a person to a tile
// returns true/false accordingly
func (p *Person) moveTo(t *tile) bool {
	if t == nil {return false}
	if canGo(t) && t.occupied == nil {
		p.path[len(p.path) - 1].occupied = nil
		t.occupied = p
		p.path = append(p.path, t)
		return true
	} else {
		return false
	}
}

// attempts to move a person one step foward in their current direction
func (p *Person) followDir() bool{
	if len(p.plan) == 0 {return p.moveTo(p.currentTile().safestTile())}
	if p.currentTile() == p.plan[0] {
		if len(p.plan) > 1 {p.dir = getDir(p.plan[0], p.plan[1])} //else {panic("plan?")}
		p.plan = p.plan[1:]
	}
	return p.moveTo(p.nextTile())   
}

// returns the next tile for a person based on their current direction
func (p *Person) nextTile() *tile{
	if p.dir == n {return p.currentTile().neighborNorth}
	if p.dir == e {return p.currentTile().neighborEast}
	if p.dir == s {return p.currentTile().neighborSouth}
	if p.dir == w {return p.currentTile().neighborWest}

	if p.dir == nw {return p.currentTile().neighborNW}
	if p.dir == ne {return p.currentTile().neighborNE}
	if p.dir == se {return p.currentTile().neighborSE}
	if p.dir == sw {return p.currentTile().neighborSW}

	return nil // default?
}

// moves a person according to their current state
// leaving the building if standing on a door,
// moving forward in current direction if possible,
// redirecting otherwise
// then updating the persons internal clock
func (p *Person) followPlan() {
	if p.path[len(p.path) - 1] == nil {return} // person escaped, this case should never happen though
	if p.path[len(p.path) - 1].door {   // standing at the exit -> leave
		p.currentTile().occupied = nil
		p.path = append(p.path, nil)
		p.save()
	} else { // follow tha plan!
		if p.followDir() {   // next step in plan is available -> move
		} else if !p.redirect() {p.wait()}  // next step in plan is occupied -> redirect or w8
	}
	p.updateTime()	
}

// makes a person 'wait' for 1 step
func (p *Person)wait() { // just chillin'
	p.path = append(p.path, p.path[len(p.path) - 1])
}

// checks if a person waited during the previous step
func (p *Person)IsWaiting() bool{
	if p == nil {return false}
	if len(p.path) <= 1 {
		return false
	} else {return p.path[len(p.path) - 1] == p.path[len(p.path) - 2]}
}

// markes a person as dead
func (p *Person) kill() {
	p.alive = false
	p.screwed = true
//	p.hp = 0
//	p.currentTile().occupied = nil // If u wanna run over corpses.
}

// markes a person as saved, used when they leave the building
func (p *Person) save() {
	p.safe = true
	p.path[len(p.path) - 1] = &tile{}
}

// modifies a persons direction, used if a new plan has been set
// selects a 'safest tile' for a person if they are screwed
func (p *Person) updatePlan(m *[][]tile) {  //OBS: Function has been reduced greatly, is more like 'updateDir' right now..,
	if len(p.plan) == 0 || (len(p.plan) > 0 && !canGo(p.plan[0])) {
		p.screwed = true
		sf := p.currentTile().safestTile()
		if sf != nil {
			p.plan = []*tile{sf}		
		}
	}
	if len(p.plan) > 0 {p.dir = getDir(p.currentTile(), p.plan[0])}
} 

// main function to move ppl forward
// updates the persons plan,
// moves them if possible
// and updates their stats accordingly
func (p *Person) MovePerson(m *[][]tile) {	
	if p == nil {return}
	if p.safe || !p.alive {
		return
	}
	if p.time <= step {
		p.updatePlan(m)
		p.followPlan()
	} 
	p.updateStats()
}

// runs a simulation through the terminal
// used for bug catching etc without using the GUI
func MovePeople(m *[][]tile, ppl []*Person) {
	var wg sync.WaitGroup
	InitPlans(m)
	for !CheckFinish(ppl) {
		wg.Add(len(ppl))
		print("\033[H\033[2J")  // empties the terminal print for a 'prettier' display
		PrintTileMapP(*m)
		for _, pers := range ppl {
		//	if pers.path[0].xCoord == 32 && pers.path[0].yCoord == 86 {fmt.Println("buggy's at: ", pers.currentTile())}
			go func(p *Person){
				defer wg.Done()				
				p.MovePerson(m)
			}(pers)
		}
		step++	
		wg.Wait()
		if math.Mod(float64(step),2) == 0 {
			FireSpread(*m)
			SmokeSpread(*m)
			InitPlans(m)
		}
		time.Sleep(700 * time.Millisecond)
	}	
}

// finds the current tile of a person
func (p *Person)currentTile() *tile{
	if len(p.path) == 0 {return nil}
	return p.path[len(p.path) - 1]
}

// updates a persons internal clock
// time added depends on wether the step was straight/diagonal and if the tile has any smoke
func (p *Person)updateTime() {
	if p.wasDiag() {
		p.time += float32(math.Sqrt(2))	
	} else {p.time += 1}

	if p.currentTile() != nil && p.currentTile().smoke > 0 {
		smoke := float32(p.currentTile().smoke/50)
		p.time += smoke}
}

// checks if the next step a diagonal one
func (p *Person)DiagonalStep() bool{ 
	if p == nil {return false}
	if len(p.plan) < 1 {return false}
	return Diagonal (p.path[len(p.path) - 1], p.plan[0])
}

// checks if the last step was a diagonal one
func (p *Person)wasDiag() bool{
	if len(p.path) < 2 {return false}
	return Diagonal(p.path[len(p.path) - 1], p.path[len(p.path) - 2])
}

// auxiliary function for diagonalStep() and wasdiag()
func Diagonal(t1, t2 *tile) bool {
	if t1 == nil {return false}
	if t2 == nil {return false}
	if t1.neighborNW == t2 {return true}
	if t1.neighborNE == t2 {return true}
	if t1.neighborSE == t2 {return true}
	if t1.neighborSW == t2 {return true}
	return false
}

// main function in people
// sets the path for all people who are currently occupying a tile in the map
func InitPlans(m *[][]tile) {
	doors := []*tile{}
	for i, list := range *m {
		for j, _ := range list {	
			if (*m)[i][j].door {
				doors = append(doors, &(*m)[i][j])}
		}
	}
	getPath(m, doors)	
}

// returns the stats for 1 person,
// used to pipe the correct data to gui.
func (p *Person) GetStats() []int {
	aslice := make([]int, 0)
	aslice = append(aslice, p.currentTile().yCoord)
	aslice = append(aslice, p.currentTile().xCoord)
	aslice = append(aslice, p.hp)
	//	*aslice = append(*aslice, p.currentTile().yCoord)
	//	*aslice = append(*aslice, p.currentTile().xCoord)
	//	*aslice = append(*aslice, p.hp)
	//aslice[2] = p.hp
	return aslice
}
