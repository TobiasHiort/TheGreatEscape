package main

import "fmt"

// map prints 
func PrintTileP(thisTile tile) {
	if thisTile.xCoord == 28 && thisTile.yCoord == 49 {
		fmt.Print("Ö")
	} else if thisTile.occupied != nil{
		fmt.Print("X")
	} else if thisTile.wall {
		fmt.Print("1")
	} else if thisTile.door {
		fmt.Print("2")
	} else if thisTile.outOfBounds {
		fmt.Print("3")
	} else if thisTile.heat > 0{
		fmt.Print("*")
	} else if thisTile.smoke > 0 {
		fmt.Print("'")
	} else {
		//fmt.Print("0")
		fmt.Print(" ")
	} 
}

func PrintTileMapP(inMap [][]tile) {
	mapXSize := len(inMap)
	mapYSize := len(inMap[0])

	for x:= 0; x < mapXSize; x++ {
		for y:= 0; y < mapYSize; y++{
			PrintTileP(inMap[x][y])
		}
		fmt.Print("\n")
	}
}

func PrintTile(atile tile) {
	fmt.Println(atile.xCoord, atile.yCoord)
	fmt.Print("\n")
}

func PrintNeighbors(atile tile) {
	if atile.neighborNorth != nil {
		fmt.Print("North: ")
		PrintTile(*(atile.neighborNorth))
	//	fmt.Print("\n")
	} else {
		fmt.Print("North: nil\n")
	}
	if atile.neighborWest != nil {
		fmt.Print("West: ")
		PrintTile(*(atile.neighborWest))
		//fmt.Print("\n")
	} else {
		fmt.Print("West: nil\n")
	}
	if atile.neighborEast != nil {
		fmt.Print("East: ")
		PrintTile(*(atile.neighborEast))
		//fmt.Print("\n")
	} else {
		fmt.Print("East: nil\n")
	}
	if atile.neighborSouth != nil {
		fmt.Print("South: ")
		PrintTile(*(atile.neighborSouth))
		//fmt.Print("\n")
	} else {
		fmt.Print("South: nil\n")
	}
	if atile.neighborNW != nil {
		fmt.Print("NW: ")
		PrintTile(*(atile.neighborNW))
	//	fmt.Print("\n")
	} else {
		fmt.Print("NW: nil\n")
	}
	if atile.neighborNE != nil {
		fmt.Print("NE: ")
		PrintTile(*(atile.neighborNE))
	//	fmt.Print("\n")
	} else {
		fmt.Print("NE: nil\n")
	}
	if atile.neighborSE != nil {
		fmt.Print("SE: ")
		PrintTile(*(atile.neighborSE))
		//fmt.Print("\n")
	} else {
		fmt.Print("SE: nil\n")
	}
	if atile.neighborSW != nil {
		fmt.Print("SW: ")
		PrintTile(*(atile.neighborSW))
		//fmt.Print("\n")
	} else {
		fmt.Print("SW: nil\n")
	}	
}

func printPath(path []*tile) {
	if path == nil {
		fmt.Println("No valid path exists")
	}
	for i, t := range path {
		if (t == nil) {
			fmt.Println("End")
		} else {fmt.Println(i, ":", t.xCoord, ",", t.yCoord)}
	}
}
