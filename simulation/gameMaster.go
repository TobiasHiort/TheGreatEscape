package main

//toppen av stupröret, yttersta lagret av löken

type personEvent struct {
	newPos *tile
	newAliveState *bool
	currentPerson *Person
}

type updateEvent {
	//each person has a spot in this array. if they dont move, the coords are blank
	//if they die, maybe a -1,-1???
	personMoveList		*([][]int)
	//sends a map of the coords that have fire that has increased
	fireIncreaseList	*([][]int)
	//timeStamp int
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
