map: map.go people
	go build map.go

test: map.go map_test.go
	go test -v map.go map_test.go

path: map.go queue.go pathfinder.go
	go run map.go queue.go pathfinder.go

qtest: queue.go queue_test.go
	go test -v ''

people:
	go build people.go

all: map 
	go build gameMaster.go

gameMaster: clock.go gameMaster.go map.go
	go build clock.go gameMaster.go map.go
