package main

//toppen av stupröret, yttersta lagret av löken

//time := int
//time = 0

var time *(int)
//*time = 0

type Event struct {
	tileMap [][]tile
	timeStamp int
}

func TimeInit(){
}

//func get ppl movement from map

func timeStamp(inMap [][]tile, currentTime int) Event{
	var ev Event
	ev.tileMap = inMap
	ev.timeStamp = currentTime
	return ev
}

//this function will in the future preferably grab just the changes on the map, rather than it as a whole
//outEvent := Event
//outEvent = timeStamp(tileMap, time)
//grabs data from map
//returns an Event

func Tick(){
	//should run as goroutine
	*time++
}

func MapInit(peopleList [][]int, newMap [][]int) [][]tile{
	//gets map data from GM and inits the map
	var currentMap ([][]tile)
	currentMap = TileConvert(newMap)

	//peopleArray := [len(peopleList)]Person{}

	//PeopleInit(peopleArray)
	//Populate()
	//gets ppl data from GM and calls the initpppl function in map 
	return currentMap
}

/*
func fetch(inMap [][]tile){
}
*/


func Discretize() { //Vemvare?
	// should be a goroutine

	// save timeStamp((fetch(*tileMap)), *time) to mem
}
