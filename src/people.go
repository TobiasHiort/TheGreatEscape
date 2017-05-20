package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"sync"
)

var step = float32(0)

type Person struct {
	alive bool
	safe  bool
	hp    int
	path  []*tile
	plan  []*tile
	dir Direction
	time float32
}

func (p *Person)getStats(aslice *[]int) {
	*aslice = append(*aslice, p.currentTile().yCoord)
	*aslice = append(*aslice, p.currentTile().xCoord)
	*aslice = append(*aslice, p.hp)
}

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
	person.path = append(person.path, t)
	person.hp = 100 
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

func (t *tile) getDamage() int {
	damage := int(0)
	damage += 10*int(t.fireLevel) 
//	if t.smoke > 2 {damage += 1}
	return damage
}

func (p *Person) moveTo(t *tile) bool {
	if canGo(t) && t.occupied == nil {
		p.path = append(p.path, t)
		return true
	} else {
		return false
	}
}

func (p *Person) followDir() bool{
	if p.currentTile() == p.plan[0] {
		if len(p.plan) > 1 {p.dir = getDir(p.plan[0], p.plan[1])}
		p.plan = p.plan[1:]
	}
	return p.moveTo(p.nextTile())   
}

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

func (p *Person) followPlan() {
	if p.path[len(p.path) - 1] == nil { return} // TODO updatestats
	if p.path[len(p.path) - 1].door {   // standing at the exit -> leave
		(p.path[len(p.path) - 1].occupied) = nil
		p.path = append(p.path, nil)  // replace with safezone?
		p.updateTime()
		p.save()
	} else /*if len(p.plan) > 0 */{ // follow tha plan!
		if p.followDir() {   // next step in plan is available -> move		
			p.updateTime()  
		} else { // next step in plan is occupied -> redirect or w8
			if !p.redirect() {p.wait()}
			p.updateTime()	
		}
	}/*else {		
		p.kill()
		// TODO: no valid path! panic behavior? lay down and w8 for death?
		// insert a followdir(randomDir)/move to 'safest' nearby tile?
	}*/
}

func (p *Person)wait() { // just chillin'
	p.path = append(p.path, p.path[len(p.path) - 1])
}

func (p *Person)IsWaiting() bool{
	if p == nil {return false}
	if len(p.path) <= 1 {
		return false
	} else {return p.path[len(p.path) - 1] == p.path[len(p.path) - 2]}
}

func (p *Person) kill() {
	p.alive = false
//	p.currentTile().occupied = nil // If u wanna run over corpses.
}

func (p *Person) save() {
	p.safe = true
	p.path[len(p.path) - 1] = &tile{}
}

func (p *Person) updatePlan(m *[][]tile) {  //OBS: Function has been reduced greatly, is more like 'updateDir' right now.., 
	if len(p.plan) > 0 {
		p.dir = getDir(p.currentTile(), p.plan[0])
	} else {
		xDir := rand.Intn(1) - rand.Intn(1)
		yDir := rand.Intn(1) - rand.Intn(1)
		p.dir = Direction{xDir, yDir}
	}
} 

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

func MovePeople(m *[][]tile, ppl []*Person) {
	var wg sync.WaitGroup
	InitPlans(m)
	for !CheckFinish(ppl) {
		wg.Add(len(ppl))
		print("\033[H\033[2J")
		PrintTileMapP(*m)
		for _, pers := range ppl {
			if pers.path[0].xCoord == 32 && pers.path[0].yCoord == 86 {fmt.Println("buggy's at: ", pers.currentTile())}
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
		time.Sleep(70 * time.Millisecond)
	}	
}

func (p *Person)currentTile() *tile{
	if len(p.path) == 0 {return nil}
	return p.path[len(p.path) - 1]
}

func (p *Person)updateTime() {
	if p.wasDiag() {
		p.time += float32(math.Sqrt(2))	
	} else {p.time += 1}

	if p.currentTile() != nil && p.currentTile().smoke > 0 {
		smoke := float32(p.currentTile().smoke/50)
		p.time += smoke}
}

func (p *Person)DiagonalStep() bool{ // Is the next step a diagonal one?
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

func InitPlans(m *[][]tile) {
	doors := []*tile{}
	for i, list := range *m {
		for j, _ := range list {	
			if (*m)[i][j].door {
				doors = append(doors, &(*m)[i][j])}
		}
	}
	getPath3(m, doors)	
}

func (p *Person) redirectOld(m *[][]tile) {
	//	newPlan, ok := getPPath(m, p.currentTile(), p.plan[0])
	current := p.currentTile()
	//----

	if p.dir == nw {
		if p.moveTo(current.neighborNorth) {
			p.plan = append([]*tile{current.neighborWest}, p.plan...)
			//p.plan = append( p.plan, current.neighborWest)
		} else if p.moveTo(current.neighborWest) {
			p.plan = append([]*tile{current.neighborNorth}, p.plan...)
			//p.plan = append( p.plan, current.neighborNorth)
		}
	}
	if p.dir == ne {
		//	if current.neighborNorth.occupied == nil { }
		if p.moveTo(current.neighborNorth) {
			p.plan = append([]*tile{current.neighborEast}, p.plan...)
			//p.plan = append( p.plan, current.neighborEast)
		} else if p.moveTo(current.neighborEast) {
			p.plan = append([]*tile{current.neighborNorth}, p.plan...)
			//p.plan = append( p.plan, current.neighborNorth)
		}
	}
	if p.dir == se {
		//	if current.neighborNorth.occupied == nil { }
		if p.moveTo(current.neighborSouth) {
			p.plan = append([]*tile{current.neighborEast}, p.plan...)
			//p.plan = append( p.plan, current.neighborEast)
		} else if p.moveTo(current.neighborEast) {
			p.plan = append([]*tile{current.neighborSouth}, p.plan...)
			//p.plan = append( p.plan, current.neighborSouth)
		}
	}
	if p.dir == sw {
		//	if current.neighborNorth.occupied == nil { }
		if p.moveTo(current.neighborSouth) {
			p.plan = append([]*tile{current.neighborWest}, p.plan...)
			//p.plan = append( p.plan, current.neighborWest)
		} else if p.moveTo(current.neighborWest) {
			p.plan = append([]*tile{current.neighborSouth}, p.plan...)
			//p.plan = append( p.plan, current.neighborSouth)
		}
	}
	
	//----
}

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
