map: map.go
	go build map.go people.go pathfinder.go queue.go
	./map

test: map.go map_test.go
	go test -v map.go people.go map_test.go pathfinder.go queue.go

path: map.go queue.go pathfinder.go
	go run map.go queue.go pathfinder.go

tests: queue.go queue_test.go map.go map_test.go pathfinder.go pathfinder_test.go
	go test -v ''

people:
	go build people.go

all: map 
	go build gameMaster.go

gameMaster: clock.go gameMaster.go map.go
	go build clock.go gameMaster.go map.go
