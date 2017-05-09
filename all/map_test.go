package main

import (
	"fmt"
	"testing"
  "github.com/stretchr/testify/assert"
)

func TestSizeOfTileConvert(t *testing.T) {

	testMatrix := [][]int{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{1, 1, 0, 1, 1},
		{0, 0, 0, 3, 3},
		{2, 0, 0, 3, 3}}

	expectedsize := 5
	expectedsize2 := 5
	amap := TileConvert(testMatrix)
	actualsizeamap := len(amap[0])
	actualsizeamap2 := len(amap)

	if expectedsize != actualsizeamap {
		t.Errorf("Expected %d, but got %d", expectedsize, actualsizeamap)
	}

	if expectedsize != actualsizeamap2 {
		t.Errorf("Expected %d, but got %d", expectedsize2, actualsizeamap)
	}

	testMatrix2 := [][]int{
		{0, 0},
		{0, 0},
		{1, 1}}

	expectedinnersize := 2
	expectedoutersize := 3
	amap2 := TileConvert(testMatrix2)
	actualinnersize := len(amap2[0])
	actualoutersize := len(amap2)

	if expectedinnersize != actualinnersize {
		t.Errorf("Expected %d, but got %d", expectedinnersize, actualinnersize)
	}

	if expectedoutersize != actualoutersize {
		t.Errorf("Expected %d, but got %d", expectedoutersize, actualoutersize)
	}
}

func TestTileConvert(t *testing.T) {
	//0 = floor, 1 = wall, 2 = door, 3 = outofbounds
	testMatrix := [][]int{
		{0, 1},
		{2, 3}}

	amap := TileConvert(testMatrix)
	actualtile := (amap[0][0]) //floor
	//x = row 0, y = column 0

  //test, expected, actual, errormessage
  assert.Equal(t, false, actualtile.wall, "Expected a floor but this is a wall")
  assert.Equal(t, false, actualtile.door, "Expected a floor but this is a door")
  assert.Equal(t, false, actualtile.outOfBounds, "Expected a floor, but this is a outofbounds")  
  
  actualtile = amap[0][1] //wall

  assert.Equal(t, true, actualtile.wall, "Expected a wall, but this is not a wall" )
  assert.Equal(t, false, actualtile.door, "Expected a wall, but this is a door")
  assert.Equal(t, false, actualtile.outOfBounds, "Expected a wall, but this is a outofbounds")  
  if actualtile.outOfBounds != false && actualtile.door != false && actualtile.wall == false {
		t.Errorf("Expected a wall, but this is a floor")
	}

	actualtile = amap[1][0] //door
  assert.Equal(t, true, actualtile.door, "Expected a door, but this is not a door" )
  assert.Equal(t, false, actualtile.wall, "Expected a door, but this is a wall")
  assert.Equal(t, false, actualtile.outOfBounds, "Expected a door, but this is a outofbounds")  
	if actualtile.outOfBounds != false && actualtile.wall != false && actualtile.door == false {
		t.Errorf("Expected a door, but this is a floor")
	}

	actualtile = amap[1][1] //out of bounds
  assert.Equal(t, true, actualtile.outOfBounds, "Expected a outofbounds, but this is not a outofbounds" )
  assert.Equal(t, false, actualtile.wall, "Expected a outofbounds, but this is a wall")
  assert.Equal(t, false, actualtile.door, "Expected a outofbounds, but this is a door")  	
  if actualtile.door != false && actualtile.wall != false && actualtile.outOfBounds == false {
		t.Errorf("Expected outofbounds, but this is a floor")
	}
}

func TestNeighbouringtiles(t *testing.T) {

	//0 = floor, 1 = wall, 2 = door, 3 = outofbounds
	testMatrix := [][]int{
		{0, 1},
		{2, 3}}

	amap := TileConvert(testMatrix)
	firsttile := amap[0][0] //0
	actualNorth := firsttile.neighborNorth
	actualEast := firsttile.neighborEast
	actualSouth := firsttile.neighborSouth
	actualWest := firsttile.neighborWest

	expectedEast := &amap[0][1]
	expectedSouth := &amap[1][0]
  
  //test, expected, actual, errormessage
  assert.Equal(t, ((*tile)(nil)), actualNorth, "Neighbor to the north is wrong")
  assert.Equal(t, expectedEast, actualEast, "Neigbor to the east is wrong")
  assert.Equal(t, expectedSouth, actualSouth, "Neigbor to the south is wrong")
  assert.Equal(t, ((*tile)(nil)), actualWest, "Neighbor to the west is wrong")

  testMatrix2 := [][]int{
		{0, 1},
		{2, 3},
		{0, 0}}

	amap2 := TileConvert(testMatrix2)
	lasttile := amap2[2][1] //0 last tile
	actualNorth = lasttile.neighborNorth
	actualEast = lasttile.neighborEast
	actualSouth = lasttile.neighborSouth
	actualWest = lasttile.neighborWest

	expectedNorth2 := &amap2[1][1]
	expectedWest2 := &amap2[2][0]

  assert.Equal(t, expectedNorth2, actualNorth, "Neighbor to the north is wrong")
  assert.Equal(t, ((*tile)(nil)), actualEast, "Neigbor to the east is wrong")
  assert.Equal(t, ((*tile)(nil)), actualSouth, "Neigbor to the south is wrong")
  assert.Equal(t, expectedWest2, actualWest, "Neighbor to the west is wrong")

	testMatrix3 := [][]int{
		{0, 1, 0},
		{2, 3, 0},
		{0, 0, 0}}

	amap3 := TileConvert(testMatrix3)
	lasttile = amap3[1][1] //3, middle tile
	actualNorth = lasttile.neighborNorth
	actualEast = lasttile.neighborEast
	actualSouth = lasttile.neighborSouth
	actualWest = lasttile.neighborWest

	expectedNorth3 := &amap3[0][1]
	expectedEast3 := &amap3[1][2]
	expectedSouth3 := &amap3[2][1]
	expectedWest3 := &amap3[1][0]

  //test, expected, actual, errormessage
  assert.Equal(t, expectedNorth3, actualNorth, "Neighbor to the north is wrong")
  assert.Equal(t, expectedEast3, actualEast, "Neigbor to the east is wrong")
  assert.Equal(t, expectedSouth3, actualSouth, "Neigbor to the south is wrong")
  assert.Equal(t, expectedWest3, actualWest, "Neighbor to the west is wrong")
}

func TestAMovement(t *testing.T) {

  testMatrix := [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
    {0, 0, 2, 0}}

    /*
    Initial people placement
    {0, X, 0, 0 },
		{0, 0, 0, X },
		{X, 0, 0, 0 },
    {0, 0, 2, 0}}
    */

  testmap := TileConvert(testMatrix)
  peoplecordinates := [][]int{
    {0,1},
    {2,0},
    {1,3}}

    peoplearray:= PeopleInit(testmap, peoplecordinates)

   //these test needs affection when we make people into threads 

   fmt.Println("Inital people placement")
   printTileMap(testmap)
   assert.Equal(t, peoplearray[0], testmap[0][1].occupied, "People1 is in the wrong tile")
   assert.Equal(t, peoplearray[1], testmap[2][0].occupied, "People2 is in the wrong tile")
   assert.Equal(t, peoplearray[2], testmap[1][3].occupied, "People3 is in the wrong tile")
   fmt.Println("\n")

   fmt.Println("One step")
   Run(testmap, peoplearray)
   printTileMap(testmap)
   assert.Equal(t, peoplearray[0], testmap[0][2].occupied, "People1 is in the wrong tile")
   assert.Equal(t, peoplearray[1], testmap[2][1].occupied, "People2 is in the wrong tile")
   assert.Equal(t, peoplearray[2], testmap[1][2].occupied, "People3 is in the wrong tile")
   fmt.Println("\n")

   fmt.Println("Two steps")
   Run(testmap, peoplearray)
   printTileMap(testmap)
   assert.Equal(t, peoplearray[0], testmap[0][2].occupied, "People1 is in the wrong tile")
   assert.Equal(t, peoplearray[1], testmap[2][2].occupied, "People2 is in the wrong tile")
   assert.Equal(t, peoplearray[2], testmap[1][2].occupied, "People3 is in the wrong tile")
   fmt.Println("\n")

   fmt.Println("Three steps")
   Run(testmap, peoplearray)
   printTileMap(testmap)
   assert.Equal(t, peoplearray[0], testmap[0][2].occupied, "People1 is in the wrong tile")
   assert.Equal(t, peoplearray[1], testmap[3][2].occupied, "People2 is in the wrong tile")
   assert.Equal(t, peoplearray[2], testmap[2][2].occupied, "People3 is in the wrong tile")
   fmt.Println("\n")

   fmt.Println("Four steps")
   Run(testmap, peoplearray)
   printTileMap(testmap)
   assert.Equal(t, peoplearray[0], testmap[1][2].occupied, "People1 is in the wrong tile")
   assert.NotEqual(t, peoplearray[1], testmap[3][2].occupied, "People2 is still in the map")
   assert.Equal(t, peoplearray[2], testmap[3][2].occupied, "People3 is in the wrong tile")
   fmt.Println("\n")

   fmt.Println("Five steps")
   Run(testmap, peoplearray)
   printTileMap(testmap)
   assert.Equal(t, peoplearray[0], testmap[2][2].occupied, "People1 is in the wrong tile")
   assert.NotEqual(t, peoplearray[2], testmap[3][2].occupied, "People3 is still in the map")
   fmt.Println("\n")

   fmt.Println("Six steps")
   Run(testmap, peoplearray)
   printTileMap(testmap)
   assert.Equal(t, peoplearray[0], testmap[3][2].occupied, "People1 is in the wrong tile")
   fmt.Println("\n")

   fmt.Println("Seven steps")
   Run(testmap, peoplearray)
   printTileMap(testmap)
   assert.NotEqual(t, peoplearray[0], testmap[3][2].occupied, "People1 is still in the map")
   fmt.Println("\n")
}


func TestFireSpread(t* testing.T) {

  testMatrix := [][]int{
    {0, 1, 0},
    {2, 3, 0},
    {0, 0, 0}}

    amap := TileConvert(testMatrix);

    SetFire(&amap[1][1]);

    for i := 1; i < 100; i++{
      FireSpread(amap)
    }
    assert.Equal(t, 3, amap[0][1].fireLevel, "The fire has not spread properly")
    assert.Equal(t, 3, amap[1][2].fireLevel, "The fire has not spread properly")
    assert.Equal(t, 3, amap[2][1].fireLevel, "The fire has not spread properly")
    assert.Equal(t, 3, amap[1][0].fireLevel, "The fire has not spread properly")

  }


  func printTile(thisTile tile) {
	if thisTile.wall {
		fmt.Print("[Wall(")
	} else if thisTile.door {
		fmt.Print("[Door(")
	} else if thisTile.outOfBounds {
		fmt.Print("[Out(")
	} else {
		fmt.Print("[Floor(")
	}
  if thisTile.occupied != nil {
	fmt.Print("X")
}
	//fmt.Print(" Heat: ")
	//fmt.Print(thisTile.heat)
	fmt.Print(")] ")
}

func printTileMap(inMap [][]tile) {
	mapXSize := len(inMap)
	mapYSize := len(inMap[0])

	for x := 0; x < mapXSize; x++ {
		for y := 0; y < mapYSize; y++ {
			printTile(inMap[x][y])
		}
		fmt.Print("\n")
	}
}

func printNeighbors(atile tile) {
	if atile.neighborNorth != nil {
		fmt.Print("North: ")
		printTile(*(atile.neighborNorth))
		fmt.Print("\n")
	} else {
		fmt.Print("North: nil\n")
	}
	if atile.neighborWest != nil {
		fmt.Print("West: ")
		printTile(*(atile.neighborWest))
		fmt.Print("\n")
	} else {
		fmt.Print("West: nil\n")
	}
	if atile.neighborEast != nil {
		fmt.Print("East: ")
		printTile(*(atile.neighborEast))
		fmt.Print("\n")
	} else {
		fmt.Print("East: nil\n")
	}
	if atile.neighborSouth != nil {
		fmt.Print("South: ")
		printTile(*(atile.neighborSouth))
		fmt.Print("\n")
	} else {
		fmt.Print("South: nil\n")
	}
}

/*
func main() {
	testMatrix := [][]int{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{1, 1, 0, 1, 1},
		{0, 0, 0, 3, 3}}

		amap := TileConvert(testMatrix)
		printTileMap(amap)
		fmt.Print("\n")
		printNeighbors(amap[0][0])

		//fire testing
		SetFire(&(amap[2][2]))
		printTileMap(amap)

		for i := 0; i < 100; i++{
			FireSpread(amap)
			fmt.Println("\n")
			printTileMap(amap)
		}
	}
*/
