package main

import (
	"fmt"
	
)


type Person struct {
	alive bool
	safe bool     
	hp float32
	path []*tile  // TODO: gör om till pekare!!
	plan []*tile  //   --''--	
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
	if currentTile == p.plan[len(p.plan) - 1] {	
		// borde man ta skada på sista tilen?
		// ska man ockupera dörr-tilen lr bara gå ut?
		// kan dörren brinna etc
		p.save()		
	} else {
		(p.path[len(p.path) - 2]).occupied = false
		currentTile.occupied = true
		p.hp = p.hp - currentTile.getDamage()
		if p.hp <= 0 {
			p.kill()
		}
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
	if p.moveTo(p.plan[0]) {
		p.plan = p.plan[1:]
	}
}

/*
func setPlan(p *Person, m Map) { // implement A*
	//You shall not pass!
}
*/

func (p *Person)kill() {
	p.alive = false
}

func (p *Person)save() {
	p.safe = true
}

func (p *Person) updatePath(m *[][]tile) {
	path, ok := getPath(m, p.path[len(p.path) - 1], p.plan[len(p.plan) - 1])
	// TODO: checka om det finns typ 'path.tail' i Go
	if ok {
		p.path = path
	}
}

func mainPeople() {

	matrix := [][]int {
		{0,0,0,0,0,0,0},
		{0,0,0,0,0,0,0},
		{0,0,0,1,0,0,0},
		{0,0,0,1,0,0,0},
		{0,0,0,1,0,0,0},
		{0,0,0,0,0,0,0}, 
		{0,0,0,0,0,0,0}}
	testmap := TileConvert(matrix)

	start1 := &testmap[0][0]
	start2 := &testmap[0][6]
	var p1 = makePerson(start1)
	var p2 = makePerson(start2)

	goal := &testmap[6][3]

	plan1, _ := getPath(&testmap, start1, goal)
	plan2, _ := getPath(&testmap, start2, goal)

	p1.plan = plan1
	p2.plan = plan2

	fmt.Println(len(plan1))
	
	 for !p1.safe || !p2.safe{
		 if !p1.safe { p1.followPlan() }
		 if !p2.safe { p2.followPlan() }

		 fmt.Println("p1:", p1.path[len(p1.path) - 1])
		 fmt.Println("p2:", p2.path[len(p2.path) - 1])
		 fmt.Println("- - - - - - -")
	 }

	//Note: worksisch, TODO: create default tile for 'safe' people?

	
	
	/*	fmt.Println(p1)
	t1 := Tile{1,1}
	t2 := Tile{2,2}
	
	moveTo(&p1, t1)
	fmt.Println(p1)
	p1.plan = append(p1.plan, t2)
	fmt.Println(p1)

	followPlan(&p1)
	fmt.Println(p1)

	kill(&p1)
	fmt.Println(p1)/



/*	
	fmt.Println(getNeighbours(m, t1))
	
	path := getpath(m, t1, t4)
	fmt.Println(path)
	fmt.Println("\nThey shall all burn!");
	*/

}

