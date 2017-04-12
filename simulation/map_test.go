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


