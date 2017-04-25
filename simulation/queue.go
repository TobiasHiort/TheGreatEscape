package main

import "sort"
//import "fmt"

type tileCost struct {
	tile tile
	cost float32
	//cost should always be positive! I think..
	//also possibly change cost to double?-
	//depending on later implementation choices
}

type queue []tileCost

func (slice queue) Len() int {
	return len(slice)
}

func (slice queue) Less(i, j int) bool {
	return slice[i].cost < slice[j].cost
}

func (slice queue) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortQueue(q *queue) {
	sort.Sort(q)
}


func (q *queue) Add(t tile, c float32) {
	s := tileCost{}
	s.tile = t
	s.cost = c
	if q.inQueue(t) {	
		q.Update(t, c)
	} else {
		*q = append(*q, s)  //how to write--?
	}
	sortQueue(q)
}

func (q *queue) AddTC(tc tileCost) {
	q.Add(tc.tile, tc.cost)
}

func (q *queue) Remove(t tile) {
//	fmt.Println("check")
	if (len(*q) == 0) {return }
	i := 0
	k := true
	for k && i < len(*q) {
		if (*q)[i].tile == t {
 			k = false 
		} else if i < (*q).Len()  {
			i++			
		}
		
	}
	if i < (*q).Len() {
		*q = append((*q)[:i], (*q)[i+1:]...)
	/*	fmt.Println("here?")
		fmt.Println(i)
		*q = append((*q)[:i-1], (*q)[i:]...)*/
	}	
}

func (q *queue) Update(t tile, c float32) {
	if q.inQueue(t){
		q.Remove(t)
		q.Add(t, c)
	}
}

func (q *queue) Pop() tileCost {
	if q.Len() == 0 {return tileCost{}} //tileCost{nil, nil}}
	t := (*q)[0]//q[q.Len()-1]
	q.Remove(t.tile)
	return t	
}

func (q queue) inQueue(t tile) bool {	
	for i := 0; i < len(q); i++ {
		if q[i].tile == t {
			return true
		}
	}
	return false
}

func (q queue) costOf(tile tile) float32 {
	for _, t := range q {
		if t.tile == tile {return t.cost}
	}
	return -1   // negative cost! null value wont work..
}

