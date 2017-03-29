# THE GREAT ESCAPE's REPO
## TODO:
#### Tiles
- [ ] Every tile is an object. Every tile needs a few variables.
	- int heat (1-10)
	- bool wall
		- thinner walls?
	- [ ] Lock
		- Reserved
		- Occupied

#### Fire
- [ ] spreading algorithm
	- checks adjacent tiles for heat. If heat > 9 then burn
	- pushes up the heat on adjacent tiles
- pipes wherabouts to people

#### People
- A\* movement algorithm
	- Heuristics - how scary is the fire vs how nice is the door?
- checks locks of adjacent tiles when moving
	- can lock one tile for booking
	-	can lock one tile for occupation
- recalcs movement as soon as either FIRE or PEOPLE are in the way. 

#### Rooms
- separate instances




## USED CODE AND LINKS
[CODE] [A* pathfinding algorithm](http://code.activestate.com/recipes/578919-python-a-pathfinding-with-binary-heap/)  
[WIKI] [Von Neumann neighborhood](https://en.wikipedia.org/wiki/Von_Neumann_neighborhood)
