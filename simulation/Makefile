map: map.go
	go build map.go
	./map

test: map.go map_test.go
	go test -v map.go map_test.go

path: map.go queue.go pathfinder.go
	go run map.go queue.go pathfinder.go

qtest: queue.go queue_test.go
	go test -v ''
