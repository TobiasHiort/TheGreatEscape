package main

import "testing"

func makeTestMap(xSize, ySize int) [][]tile{
	testMatrix := [][]int{}

	for x := 0; x < xSize; x++ {
		row := []int{}
		for y := 0; y < ySize; y++ {
			row = append(row, y)
		}
		testMatrix = append(testMatrix, row)
	}
		
	return TileConvert(testMatrix)
}
/*
func TestValidTile() {

}*/

func TestCompactPath(t* testing.T) {
//	xSize := 10
//	ySize := 10

//	testMap := makeTestMap(xSize, ySize)

	
}
