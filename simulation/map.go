package tileConvert

type tile struct {
	xCoord int
	yCoord int

	/*
	neighborNorth *tile
	neighborEast  *tile
	neighborSouth *tile
	neighborWest  *tile
	*/

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
func makeNewTile(thisPoint int, x int, y int) tile{

	//några gemensamma grejer

	newTile := tile{x, y, 0, 0, false, false, 0}
		/*
		xCoord := x,
		YCoord := y,

		heat := 0,
		fireLevel := 0,

		wall := false,

		occupied := false
		//make to a ponter l8ter
		personID := 0
	}
	*/

	if thisPoint == 0 {
		//make normal floor
		//helt normalt flour

		//append to tilemap
	} else if thisPoint == 1 {
		//wall
	} else if thisPoint == 2 {
		//door
	} else if thisPoint == 3 {
		//out of bounds
	}

	return newTile
}

func tileConvert(inMap [][]int) [][]tile{
	const mapXSize := len(inMap[0])
	const mapYSize := len(inMap)

	tileMap := [mapXSize][mapYSize]tile{}

	for x:= 0; x < mapXSize; x++{
		for y:= 0; y < mapYSize; y++{
			thisPoint := inMap[x][y]
			tileMap[x][y] := makeNewTile(thisPoint, x, y)
		}
	}
	return tileMap
}
