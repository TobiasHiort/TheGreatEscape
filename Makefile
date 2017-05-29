#Use this to compile and run The Great Escape
run:
	cd src && make all
	cd gui && python3 gui.py

#Automated testing
test:
	cd src && make test
	
#For use when program is interrupted, and not able to remove temporary files by itself
clean:
	rm -r -f gui/__pycache__/
	rm -f gui/mapStats.txt
	rm -f gui/peopleStats.txt
	rm -f gui/timeStats.txt
	rm -f src/main
	rm -f src/pid.txt
