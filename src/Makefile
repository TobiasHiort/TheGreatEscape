map: map.go
	go build map.go people.go pathfinder.go queue.go print.go
	./map

#test: map.go map_test.go
#	go test -v map.go people.go map_test.go pathfinder.go queue.go

con: concurrency_test.go
	go test -v concurrency_test.go map.go people.go pathfinder.go queue.go print.go

path: map.go queue.go pathfinder.go
	go run map.go queue.go pathfinder.go

test_q: queue.go queue_test.go map.go pathfinder.go pathfinder_test.go map_test.go print.go
	go test -v queue.go pathfinder.go map.go people.go print.go queue_test.go map_test.go pathfinder_test.go people_test.go

test_c: concurrency_test.go map_test_speed.go people.go queue.go print.go pathfinder.go 
	go test -v concurrency_test.go map_test_speed.go people.go queue.go print.go pathfinder.go  

clear:
	rm main

people:
	go build people.go

all: main.go gameMaster.go map.go people.go queue.go print.go pathfinder.go stats.go
	go build main.go gameMaster.go map.go people.go queue.go print.go pathfinder.go stats.go

gameMaster: clock.go gameMaster.go map.go
	go build clock.go gameMaster.go map.go

gui: gotest.go
	go build gotest.go map.go queue.go pathfinder.go people.go print.go 	

stats: stats.go
	go build stats.go map.go people.go pathfinder.go queue.go print.go


py: gui
	python3 gui.py
