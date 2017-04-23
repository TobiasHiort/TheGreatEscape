package main


import	"fmt"




func getpath(m []Tile, from Tile, to Tile) []Tile{

//	checked := []Tile{}
	//	checked = append(checked, from)
	checked := []Tile{}  // m[:0]
//	unchecked := []Tile{}
	unchecked := m[1:]
	cQueue := Queue{}


	//initialize the costqueue 
	cQueue.Add(from, 0)
	for i := 1; i < len(m); i++ {
		cQueue.Add(m[i], 10)   // 10~infinite		
	}

	// initiate parentof map
	var parentOf map[Tile]Tile
	parentOf = make(map[Tile]Tile)

	// essential loop
	for len(unchecked) != 0 {
		fmt.Println("hit")
		current := (&cQueue).Pop()
		fmt.Println(current)
		neighbours := getNeighbours(m, current.tile)
		for _, neighbor := range neighbours {
			if neighbor == to {
				return compactPath(parentOf, from, to)
			}

			
			if contains(checked, neighbor) {
				//nothing
			} else {
				cost := current.cost + 1  // 1 cost per step? define
				if cost < cQueue.costOf(neighbor) {
					fmt.Println("parentmap?")
					parentOf[neighbor] = current.tile
					fmt.Println("parentmap??")
					cQueue.Update(neighbor, cost)					
				} 
			}
			
		}
		checked = append(checked, current.tile) //checked.Append(current)
		remove(unchecked, current.tile)
	}

	return compactPath(parentOf, from, to)
}

func compactPath(parentOf map[Tile]Tile, from Tile, to Tile) []Tile {
	path := []Tile{}

	current := to

	for current != from {
		path = append(path, parentOf[current])//path.Append(parentOf[current])
	}

	return append(path, from)	
}

func contains(slice []Tile, tile Tile) bool {
	for _, t := range slice {
		if t == tile {
			return true
		}
	}
	return false
}

func remove(slice []Tile, tile Tile) {
	for i, t := range slice {
		if t == tile {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
}


//}


func getNeighbours(m []Tile, t Tile) []Tile {
	neighbours := []Tile{}

	for i := 0; i < len(m); i++ {//m.Len(); i++ {
		current := m[i]
		if current.X == t.X {
			if current.Y == t.Y + 1 || current.Y == t.Y - 1 {
				neighbours = append(neighbours, current)
			}		
		} else if current.Y == t.Y {
			if current.X == t.X + 1 || current.X == t.X - 1 {
				neighbours = append(neighbours, current)
			}
		}
	}

	return neighbours
}

