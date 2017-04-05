# ![run_left](http://i.imgur.com/qT9yxGX.png) The Great Escape ![run_left](http://i.imgur.com/ttqg197.png)
The Great Escape is a building evacuation simulator written in Go (back-end) and Python (front end) with support for custom maps.

## 1. Getting Started
Clone the repository, and for now run the Python part and Go part separately.

### 1.1 Prerequisites
The program is written in Python (3.x) and Go, see links below:
* [Python 3.x](https://www.python.org/downloads/)
    * [Pygame](https://www.pygame.org/wiki/GettingStarted#PygameInstallation)
    * [Numpy](https://www.scipy.org/scipylib/building/index.html#building)
    * [Tkinter](https://wiki.python.org/moin/TkInter)
* [Go](https://golang.org/)

### 1.2 Installing

#### 1.2.1  Linux
Instructions for setting up a development/user environment.

```
$ sudo apt-get update
```
```
$ sudo apt-get install python3
```
Install pip for managing software packages within Python:
```
$ sudo apt-get install -y python3-pip
```
Install Pygame:
```
$ pip3 install pygame
```
Install Numpy:
```
$ pip3 install numpy
```
Install Tkinter:
```
$ apt-get install python3-tk
```
## 2. Deployment
### 2.1 Custom maps
Add custom maps in `/gui/python/maps` and load them in the Python program. Maps must be in `.PNG` format and consist of the following pixels without transparency:
- ![#ffffff](https://placehold.it/15/ffffff/000000?text=+) `rgb(255, 255, 255)` &nbsp;&nbsp;- Floor.
- ![#000000](https://placehold.it/15/000000/000000?text=+) `rgb(0, 0, 0)` &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Wall.
- ![#00ff00](https://placehold.it/15/00ff00/000000?text=+) `rgb(0, 255, 0)` &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- Door.
- ![#7f7f7f](https://placehold.it/15/7f7f7f/000000?text=+) `rgb(127, 127, 127)` &nbsp;&nbsp;- Out of bounds (outdoors).

### 2.2 Run program
Run Python front end from `/gui/python` with:
```
python3 test.py
```

## 3. Created With
* [GitHub](https://github.com/) - Version control repository
* [Slack](https://slack.com/) - Team collaboration
* [Trello](https://trello.com/) - Scrum management
* [Photoshop](http://www.adobe.com/products/photoshop.html) - GUI development

## 4. Authors
* **Tobias Hiort**

## 6. License
This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## 7. USED CODE AND LINKS (!update)
[CODE] [A* pathfinding algorithm](http://code.activestate.com/recipes/578919-python-a-pathfinding-with-binary-heap/)
