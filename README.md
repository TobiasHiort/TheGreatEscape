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
- 

#### People
- A\* movement algorithm
	- Heuristics - how scary is the fire vs how nice is the door?
- checks locks of adjacent tiles when moving
	- can lock one tile for booking
	-	can lock one tile for occupation
- recalcs movement as soon as either FIRE or PEOPLE are in the way. 
- HEALTH: implement A\* heuristics for people to maximize health

#### Rooms
- separate instances
- int smoke

#### GUI
- Up arrow is upload right now, replace with button later
- Remove move player with right arrow, but keep the mechanics
- Determine if PNG is to large or small
- Represent player as a drawn circle instead of a png? Better scaling
- Figure out how to make click/hover buttons and how to view different tabs
- Create test matrix run file with a player moving, and figure out how to step through it graphically/run it in some speed
- Think about layering (player movement, fire, smoke)

####
- Trick A\* to think that heat tiles are slower to pass. This can be an implementation of path priority /hhuehueristics

## USED CODE AND LINKS
[CODE] [A* pathfinding algorithm](http://code.activestate.com/recipes/578919-python-a-pathfinding-with-binary-heap/)  
[WIKI] [Von Neumann neighborhood](https://en.wikipedia.org/wiki/Von_Neumann_neighborhood)
