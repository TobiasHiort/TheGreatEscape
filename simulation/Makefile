map: map.go
	go build map.go people.go pathfinder.go queue.go print.go
	./map

test: map.go map_test.go
	go test -v map.go people.go map_test.go pathfinder.go queue.go

path: map.go queue.go pathfinder.go
	go run map.go queue.go pathfinder.go

tests: queue.go queue_test.go map.go pathfinder.go pathfinder_test.go
	go test -v queue.go pathfinder.go map.go people.go print.go queue_test.go pathfinder_test.go

clear:
	rm map
