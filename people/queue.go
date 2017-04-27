package main

import "sort"
import "fmt"

type Thing struct {
	tile Tile
	cost int
}

type Queue []Thing

func (slice Queue) Len() int {
	return len(slice)
}

func (slice Queue) Less(i, j int) bool {
	return slice[i].cost < slice[j].cost
}

func (slice Queue) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortQueue(q *Queue) {
	sort.Sort(q)
}


func (q *Queue) Add(t Tile, c int) {
	s := Thing{}
	s.tile = t
	s.cost = c
	if q.inQueue(t) {
		q.Update(t, c)
	} else {
		*q = append(*q, s)  //how to write--?
	}
	sortQueue(q)
}

func (q *Queue) Remove(t Tile) {
	i := 0
	k := true
	for k {
		if (*q)[i].tile == t {
			k = false 
		} else {
			i++			
		}
	}
	if i <= (*q).Len() {
		*q = append((*q)[:i], (*q)[i+1:]...)
	}	
}

func (q Queue) Update(t Tile, c int) {
	q.Remove(t)
	q.Add(t, c)
}

func (q *Queue) Pop() Thing {
	fmt.Println("testa")
	t := (*q)[0]//q[q.Len()-1]
	fmt.Println("test")
	q.Remove(t.tile)
	return t	
}

func (q Queue) inQueue(t Tile) bool{
	for i := 0; i < len(q); i++ {
		if q[i].tile == t {
			return true
		}
	}
	return false
}


func (q Queue) costOf(tile Tile) int {
	for _, t := range q {
		if t.tile == tile {return t.cost}
	}
	return -1   // negative cost! null value wont work..
}


// todo: if-statemnt to check for invalid input
