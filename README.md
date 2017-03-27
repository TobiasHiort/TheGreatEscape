# THE GREAT ESCAPE's REPO

##### PROBLEMS
- Can start on wall. (Kinda solved @ row 102 (draw player section). Should be solved while spawning).
- Can finish with a door in a corner (or if an inner wall is adjacent to the outer wall), because diagnoal movement is allowed. Should be handled somehow.
- Geometrically speaking, diagonal movement is longer than xy-movement.
- Only one door is possible right now.
- Fire could right now be represented as walls. We might however want to add degrees of fire hazard, then we need to look into the A* algorithm to handle more than 1/0. Plus people can't die or stand in walls...

##### TODOS
- Somewhat inconsistent with [] vs. () for arrays/tuples etc.
- Walls should be thinner IMO, but then the pathfinding mechanism wont work as is.
- Untried for MxN matrix, but should work.
- Implement reading matrix from .PNGs via color coding. Not reasonable to manually facilitate a much larger matrix.
- Later on we need to count waiting on a spot as 1 (or more) time units. Current algorithm only counts steps.
- Crazy idea: put out small fires with fire extinguishers?
- We later need to solve the problem with two players trying to reach the same node at the same time unit. Random wait? Communication between the players? Rerouting?
- Separate runtime and graphics print timme


## USED CODE AND LINKS
[CODE] [A* pathfinding algorithm](http://code.activestate.com/recipes/578919-python-a-pathfinding-with-binary-heap/)

