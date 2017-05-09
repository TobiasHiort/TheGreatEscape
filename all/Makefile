map: map.go
	go build map.go people.go pathfinder.go queue.go print.go
	./map

test: map.go map_test.go
	go test -v map.go people.go map_test.go pathfinder.go queue.go

path: map.go queue.go pathfinder.go
	go run map.go queue.go pathfinder.go

tests: queue.go queue_test.go map.go pathfinder.go pathfinder_test.go
	go test -v queue.go pathfinder.go map.go people.go print.go queue_test.go pathfinder_test.go people_test.go

clear:
	rm map

people:
	go build people.go

all: map 
	go build gameMaster.go

gameMaster: clock.go gameMaster.go map.go
	go build clock.go gameMaster.go map.go

gui: 
	go build gotest.go map.go queue.go pathfinder.go people.go print.go 	

