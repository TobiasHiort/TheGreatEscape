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
	hp    int
	path  []*tile
	plan  []*tile  // plan of jps
	//	pPlan []*tile  // partial plan
	dir Direction
	time float32
	w8ed int
}

type Stats struct {
	x int
	y int
//	hp float32
}
/*

func (p *Person)getStats() []int {
	//if p == nil || p.currentTile() == nil {return Stats{}}
	//if len(p.path) < 1 {return []int{0,0}}
//	return Stats{p.path[len(p.path) - 1].xCoord, p.path[len(p.path) - 1].yCoord}
//	return []int{p.currentTile().xCoord, p.currentTile().yCoord}//, p.hp}
	return []int{p.currentTile().yCoord, p.currentTile().xCoord}//, p.hp}
======= */
func (p *Person)getStats(aslice *[]int) {
//	aslice[0] = p.currentTile().xCoord
//	aslice[1] = p.currentTile().yCoord
	*aslice = append(*aslice, p.currentTile().yCoord)
	*aslice = append(*aslice, p.currentTile().xCoord)
	*aslice = append(*aslice, p.hp)
	
  //aslice[2] = p.hp
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

func (t *tile) getDamage() int {
	damage := int(0)
	damage += 10*int(t.fireLevel) // man dör om man kliver i elden, right...
	//damage = damage + int(t.heat)   // TODO: how much does the fire hurt??
	// damage = damage + effect from smoke'n stuff
//	smoke := int(t.smoke)
	if t.smoke > 1 {damage += 1}
//	damage += smoke //int(t.smoke)
	return damage
}

func (p *Person) moveTo(t *tile) bool {
	//if validTile(t) && t.occupied == nil {
	if canGo(t) && t.occupied == nil {
		p.path = append(p.path, t)
	//	p.updateStats()
		return true
	} else {
		return false
	}
}

func (p *Person) followDir() bool{
	if p.currentTile() == p.plan[0] {
		if len(p.plan) > 1 {p.dir = getDir(p.plan[0], p.plan[1])}
		p.plan = p.plan[1:]
	//	return 
		// new jp
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
/*	if len(p.plan) > 0 { // follow tha plan!		
	//	if p.moveTo(p.plan[0]) {   // next step in plan is available -> move		
		if p.followDir() {
			//p.plan = p.plan[1:]
			p.updateTime()  
		} else {                   // next step in plan is occupied -> w8
			p.wait()
			p.updateTime()
		}          
	} else if p.path[len(p.path) - 1].door {   // standing at the exit -> leave
		(p.path[len(p.path) - 1].occupied) = nil
		p.path = append(p.path, nil)  // replace with safezone?
		p.updateTime()
		p.save()*/
	if p.path[len(p.path) - 1].door {   // standing at the exit -> leave
		(p.path[len(p.path) - 1].occupied) = nil
		p.path = append(p.path, nil)  // replace with safezone?
		p.updateTime()
		p.w8ed = 0 // testing testing
		p.save()

	} else if len(p.plan) > 0 { // follow tha plan!

		//	if len(p.plan) < 2 {fmt.Println(p.currentTile())}
		
		//	if p.moveTo(p.plan[0]) {   // next step in plan is available -> move		
		if p.followDir() {
			//p.plan = p.plan[1:]
			p.w8ed = 0 // testing testing
			p.updateTime()  
		} else {
			//fmt.Println("Redir....")// next step in plan is occupied -> w8
			if !p.redirect() {p.wait()}
			//fmt.Println("Redired!")
			p.updateTime()	
		}
	}else {
	/*	fmt.Println(p.path)
		fmt.Println(p.currentTile())
		fmt.Println("door?", p.currentTile().door)
		fmt.Println("wall?", p.currentTile().wall)
		fmt.Println("you're screwed!")*/
		p.kill()
		// TODO: no valid path! panic behavior? lay down and w8 for death?
		// idea: don't update last plan-path, follow it despite fire etc?
	}
}

func (p *Person)wait() {
	p.w8ed++
	p.path = append(p.path, p.path[len(p.path) - 1])
//	p.updateStats()	
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
	//plan, ok := getPath(m, p.path[len(p.path)-1])
	

	if !p.currentTile().door && (p.foo() || len(p.plan) < 1) {
		if len(p.plan) < 1 {fmt.Println(p.plan)}
		//plan, _ := getPath2(m, p.path[len(p.path)-1])  //changed!		if ok {			
		//p.plan = plan[1:]
		//	p.w8ed = 0
		//p.plan = p.plan[1:]
	}
//	fmt.Println("why??2")
	if len(p.plan) > 0 {p.dir = getDir(p.currentTile(), p.plan[0])} //TODO: fixa till!

//	fmt.Println("why??3")
} 
/*	ret := p.foo()
	if !p.currentTile().door && (ret || len(p.plan) < 1) {
		p.redirect(m)
		p.updateTime()
		if len(p.plan) > 0 {p.dir = getDir(p.currentTile(), p.plan[0])} //TODO: fixa till!
	}
	return ret*/
//}

func (p *Person) foo() bool{  //TODO!! checka om planen bör updates lr ej
	return false //p.w8ed > 7
}

func (p *Person) MovePerson(m *[][]tile) {	
	if p == nil {return}
	if p.safe || !p.alive {
		return
	}
	if p.time <= step {
		
	//	if step > 100 &&  p.path[0].yCoord == 9 {fmt.Println("\n path:")
	//		printPath(p.path)}//{fmt.Println(p.currentTile().xCoord, p.currentTile().yCoord)} //|| !validTile(p.plan[0]) {p.updatePlan(m)}
//
		p.updatePlan(m)
//		fmt.Print("1 done")	
		p.followPlan()
	} /*else if p.time + 1 <= step{
		p.updateStats()*/

	p.updateStats()   // DOING, outcommented in 2 other places
	

	
}

func MovePeople(m *[][]tile, ppl []*Person) {
	var wg sync.WaitGroup
	fmt.Println("next lap..")
//	SmokeSpread(*m)
//	SmokeSpread(*m)
	for !CheckFinish(ppl) {
		wg.Add(len(ppl))
		//	print("\033[H\033[2J")

		PrintTileMapP(*m)
		fmt.Print("\n")
		InitPlans(m)

	//	fmt.Println("\n", (ppl[0].plan))
	//	fmt.Println("\n", (ppl[0].plan[1]))
		
		//for _, p := range ppl[0].plan {
		//	fmt.Println(p.xCoord, p.yCoord)}
		//	time.Sleep(1000 * time.Millisecond)
	
		for _, pers := range ppl {
			if pers.path[0].xCoord == 32 && pers.path[0].yCoord == 86 {fmt.Println("buggy's at: ", pers.currentTile())}
			go func(p *Person){
				defer wg.Done()				
				p.MovePerson(m)
			}(pers)
		}
		step++	
		wg.Wait()
	//	if math.Mod(float64(step),2) == 0 {FireSpread(*m)}
		//if math.Mod(float64(step),2) == 0 {SmokeSpread(*m)}
		FireSpread(*m)
		SmokeSpread(*m)
		time.Sleep(70 * time.Millisecond)
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

	if p.currentTile() != nil && p.currentTile().smoke > 0 {
		smoke := float32(p.currentTile().smoke/50)
		
		p.time += smoke}
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

func InitPlans(m *[][]tile) {
	doors := []*tile{}
	for i, list := range *m {
		for j, _ := range list {
		
			if (*m)[i][j].door {
			//	fmt.Println("door?", (*m)[i][j])
				doors = append(doors, &(*m)[i][j])}
		}
	}

	//for _, d := range doors {
	//	fmt.Println("??", d)
	//}
//	fmt.Println("got the doors")
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
	p.w8ed = 0
/*	if ok {
		p.plan = append(newPlan, p.plan[1:]...)
//		path = append([]*tile{parentOf[current]}, path...)
	}*/
}




// merging

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
