package tileConvert

type tile struct {
	neighborNorth *tile
	neighborEast  *tile
	neighborSouth *tile
	neighborWest  *tile

	//north/south
	longitude int
	//east/west
	latitude	int

	heat			int
	fireLevel int

	wall bool

	occupied bool
	personID int
}

/*
func FireSpread(thisTile tile){
	if thisTile.heat > 9 {
		thisTile.fireLevel = 1
	}
	if thisTile.heat > 19 {
		thisTile.fireLevel = 2
	}
	if thisTile.heat > 29 {
		thisTile.fireLevel = 3
	}

	tile.*neighborNorth.heat += fireLevel
	tile.*neighborEast.heat	+= fireLevel
	tile.*neighborWest.heat	+= fireLevel
	tile.*neighborSouth.heat += fireLevel

	/*
	tile.neighborNorth.heat += heat/10
	tile.neighborEast.heat += heat/10
	tile.neighborWest.heat += heat/10
	tile.neighborSouth.heat += heat/10
}
*/

func tileConvert(inMap [][]int) [][]tile{
	tileMap := [][]tile{}

	for i:= 0; i < len(inMap[0]); i++{
		for a:= 0; a < len(inMap); a++{
			//check element

			if [i][a]inMap == 1 {
				//placera vägg här
				//appenda till tileMap ( som är en slice )
			} else if [i][a]inMap == 2 {
				//placera dörr här
			} else if [i][a]inMap == 3 {
				//placera ute
			} else if [i][a]inMap == 0 {
				//placera vanligt golv
			}
		}
	}
	return tileMap
}
