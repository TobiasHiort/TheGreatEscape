package main

import (
	"math"
	"testing"
)

func TestGetPath(t *testing.T) {

}

func TestStepCost(t *testing.T) {

	ti := makeNewTile(0, 0, 0)

	if stepCost(ti) != 1 {
		t.Errorf("Expected stepcost: 1, but got stepcost: %d", stepCost(ti))
	}

	for i := float32(0); i < 10; i++ { //TODO: om vi ändrar cost för heat så redigera testet!
		if stepCost(ti) != float32(i/5+1) {
			t.Errorf("Expected stepcost: %g, but got stepcost: %g", float32(i/5+1), stepCost(ti))
		}
		ti.heat += 1
	}
	SetFire(&ti)
	if !math.IsInf(float64(stepCost(ti)), 1) {
		t.Errorf("Expected stepcost: %g, but got stepcost: %g", float32(math.Inf(1)), stepCost(ti))
	}

	// empty tile = 1
	// heatlvl tile = 1 + heatlvl/5
	// fire tile = infinity
}

func TestGetNeighbors(t *testing.T) {
	matrix := [][]int{
		{0, 1, 0, 1, 0, 1, 0}, // no neighbors
		{1, 1, 1, 1, 1, 1, 1}, //
		{0, 0, 0, 0, 0, 0, 0},
		{0, 3, 3, 3, 3, 3, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0}}
	testmap := TileConvert(matrix)

	for i, list := range testmap {
		for j, ti := range list {
			neighbors := getNeighbors(&ti)
			if i == 0 && validTile(&ti) {
				if len(neighbors) != 0 {
					t.Errorf("Expected 0 neigbors, but got %d neighbors", len(neighbors))
				}
			} else if i == 2 {
				if len(neighbors) != 2 {
					t.Errorf("Expected 2 neigbors, but got %d neighbors", len(neighbors))
				}
			} else if (i == 5 || i == 6) && j > 0 && j < 6 {
				if len(neighbors) != 4 {
					t.Errorf("Expected 4 neigbors, but got %d neighbors", len(neighbors))
				}
			}
		}
	}
}

func TestValidTile(t *testing.T) {
	matrix := [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{3, 3, 3, 3}}
	testmap := TileConvert(matrix)

	for i, list := range testmap {
		for _, ti := range list {
			if i < 2 {
				if !validTile(&ti) {
					t.Errorf("Expected validtile, but got invalidtile")
				}
			} else {
				if validTile(&ti) {
					t.Errorf("Expected invalidvalidtile, but got validtile")
				}
			}
		}
	}
	if validTile(nil) {
		t.Errorf("Expected invalidvalidtile, but got validtile")
	}
}

/*
func TestCompactPath(t* testing.T) {
	size := 10
	maprentOf

	for i := 0; i < size; i++ {

	}
}*/

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
