# The Great Escape/Simulation
## TODO:
### Tiles
- [ ] Every tile is an object. Every tile needs a few variables.
	- Int fire (1-10)
	- Bool wall
  - Bool door
  - Neighbouring tiles
		- Thinner walls? (optimization)
	- [ ] Lock
		- Reserved
		- Occupied
    

### Initialize map from 2D-array to 2D-array of tiles
  - set * Tile Neighbours for each Tile
  - set pointers to persons

### Tile changing functions
  - Create fire, Fire spreading
  - Occupied
  - Reserved
  
### Simulation
  - Save info of persons whearabots 


### Fire
- [ ] spreading algorithm
	- checks adjacent tiles for heat. If heat > 9 then burn
	- pushes up the heat on adjacent tiles
- pipes wherabouts to people

### Players
- A\* movement algorithm
	- Heuristics - how scary is the fire vs. how nice is the door?
- Checks locks of adjacent tiles when moving
	- Can lock one tile for booking
	- Can lock one tile for occupation
- Recalcs movement as soon as either FIRE or PEOPLE are in the way. 
- HEALTH: implement A\* heuristics for people to maximize health

### Rooms
- separate instances
- int smoke

### Pathfinding
- Trick A\* to think that heat tiles are slower to pass. This can be an implementation of path priority
