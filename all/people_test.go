package main

import "testing"
import "fmt"

func TestCurrentTile(t *testing.T) {
	matrix := [][]int {
		{0,0,0},
		{0,0,0},
		{0,0,2}}
	mappis := TileConvert(matrix)

	p1 := makePerson(&mappis[1][0])
	fmt.Println(p1.currentTile())
	p1.MovePerson(&mappis)
	fmt.Println(p1.currentTile())
}
