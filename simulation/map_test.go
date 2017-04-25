package main

import "testing"

func TestSizeOfTileConvert(t* testing.T) {

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



func printTile(thisTile tile) {
	if thisTile.wall {
		fmt.Print("[vägg(")
	} else if thisTile.door {
		fmt.Print("[dörr(")
	} else if thisTile.outOfBounds {
		fmt.Print("[ute(")
	} else {
		fmt.Print("[golv(")
	}
  fmt.Print(thisTile.fireLevel)

  fmt.Print(" Heat: ")
  fmt.Print(thisTile.heat)
	fmt.Print(")] ")
}


func printTileMap(inMap [][]tile) {
	mapXSize := len(inMap)
	mapYSize := len(inMap[0])

	for x:= 0; x < mapXSize; x++{
		for y:= 0; y < mapYSize; y++{
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
