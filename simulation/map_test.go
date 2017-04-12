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

func TestTileConvert(t* testing.T) {
    //0 = floor, 1 = wall, 2 = door, 3 = outofbounds
    testMatrix := [][]int{
    {0, 1},	
    {2, 3}}

    amap := TileConvert(testMatrix) 
    actualtile := (amap[0][0]) //floor
    //x = row 0, y = column 0 

    if actualtile.wall != false {
      t.Errorf("Expected a floor, but this is a wall")
    } 

    if actualtile.door != false {
      t.Errorf("Expected a floor, but this is a door")
    } 

    if actualtile.outOfBounds != false {
      t.Errorf("Expected a floor, but this is a outofbounds")    
    }

    actualtile = amap[0][1] //wall
 
    if actualtile.wall != true {
      t.Errorf("Expected a wall, but this is not a wall")
    } 

    if actualtile.door != false {
      t.Errorf("Expected a wall, but this is a door")
    } 

    if actualtile.outOfBounds != false {
      t.Errorf("Expected a wall, but this is a outofbounds")  
    }

    if actualtile.outOfBounds != false && actualtile.door != false && actualtile.wall == false {
      t.Errorf("Expected a wall, but this is a floor")  
    }

    actualtile = amap[1][0] //door

    if actualtile.door != true {
      t.Errorf("Expected a door, but this is not a door")
    }  
    if actualtile.wall != false {
      t.Errorf("Expected a door, but this is a wall")
    } 
    if actualtile.outOfBounds != false {
      t.Errorf("Expected a door, but this is outofbounds")
    } 
    if actualtile.outOfBounds != false && actualtile.wall != false && actualtile.door == false {
      t.Errorf("Expected a door, but this is a floor")  
    }

    actualtile = amap[1][1] //out of bounds

    if actualtile.outOfBounds != true {
      t.Errorf("Expected outofbounds, but this is not a outofbounds")
    }  
    if actualtile.wall != false {
      t.Errorf("Expected a outofbounds, but this is a wall")
    } 
    if actualtile.door != false {
      t.Errorf("Expected a outofbounds, but this is door")
    } 
    if actualtile.door != false && actualtile.wall != false && actualtile.outOfBounds == false {
      t.Errorf("Expected outofbounds, but this is a floor")  
    }
 
  }




