# ![run_left](http://i.imgur.com/qT9yxGX.png) The Great Escape ![run_left](http://i.imgur.com/ttqg197.png)
The Great Escape is a building evacuation simulator developed in [Go](https://golang.org/) (back end) and [Python 3.x](https://www.python.org/downloads/) (front end) with support for custom [PNG](https://en.wikipedia.org/wiki/Portable_Network_Graphics) maps (see [2.1](https://github.com/TobiasHiort/TheGreatEscape#21-custom-maps)).




## 1. Getting Started
Clone the repository and read section [1](https://github.com/TobiasHiort/TheGreatEscape#1-getting-started) and [2](https://github.com/TobiasHiort/TheGreatEscape#2-deployment) below.
### 1.1. Prerequisites
The application is written in [Python 3.x](https://www.python.org/downloads/) and [Go](https://golang.org/):
* [Python 3.x](https://www.python.org/downloads/)
    * [Tkinter](https://wiki.python.org/moin/TkInter)
    * [PyGame](https://www.pygame.org/wiki/GettingStarted#PygameInstallation)
    * [NumPy](https://www.scipy.org/scipylib/building/index.html#building)
    * [inflect](https://pypi.python.org/pypi/inflect)
    * [matplotlib](https://matplotlib.org/users/installing.html)
    * [Pillow](https://pillow.readthedocs.io/en/latest/installation.html)
    * [psutil](https://pypi.python.org/pypi/psutil)
    * [colorama](https://pypi.python.org/pypi/colorama)
    * [termcolor](https://pypi.python.org/pypi/termcolor)
    * [SciPy](https://www.scipy.org/install.html)
* [Go](https://golang.org/)
### 1.2. Installing
#### 1.2.1.  Linux
##### Instructions for settings up a development/user environment:
```
$ sudo apt-get update
```
```
$ sudo apt-get install python3
```
Install [Tkinter](https://wiki.python.org/moin/TkInter) if needed:
```
$ sudo apt-get install python3-tk
```
Install [pip3](https://pypi.python.org/pypi/pip) for managing software packages within [Python 3.x](https://www.python.org/downloads/):
```
$ sudo apt-get install python3-pip
```
Then install [PyGame](https://www.pygame.org/wiki/GettingStarted#PygameInstallation):
```
$ sudo pip3 install pygame
```
and [NumPy](https://www.scipy.org/scipylib/building/index.html#building):
```
$ sudo pip3 install numpy
```
and [inflect](https://pypi.python.org/pypi/inflect):
```
$ sudo pip3 install inflect
```
and [matplotlib](https://matplotlib.org/users/installing.html):
```
$ sudo pip3 install matplotlib
```
and [psutil](https://pypi.python.org/pypi/psutil):
```
$ sudo pip3 install psutil
```
and [colorama](https://pypi.python.org/pypi/colorama):
```
$ sudo pip3 install colorama
```
and [termcolor](https://pypi.python.org/pypi/termcolor):
```
$ sudo pip3 install termcolor
```
and [SciPy](https://www.scipy.org/install.html):
```
$ sudo pip3 install scipy
```
Install Go
```
$ sudo apt-get install golang-go
```
Install testing in go
```
$ go get github.com/stretchr/testify/assert
```

#### 1.2.2. Windows and macOS
As long as [Python 3.x](https://www.python.org/downloads/), [Go](https://golang.org/) and [pip3](https://pypi.python.org/pypi/pip) are installed, use the [pip3](https://pypi.python.org/pypi/pip) commands in [1.2.1](https://github.com/TobiasHiort/TheGreatEscape#121--linux).

Windows users will need to look at [Unofficial Windows Binaries for Python Extension Packages](http://www.lfd.uci.edu/~gohlke/pythonlibs/#scipy) for [SciPy](https://www.scipy.org/install.html) and [Numpy+MKL](http://www.lfd.uci.edu/~gohlke/pythonlibs/#numpy).

macOS users will need to install (!TODO):
```
pip3 install -U --pre -f https://wxpython.org/Phoenix/snapshot-builds/ wxPython_Phoenix
```
if this is unsuccessful, you are pretty much out of luck for choosing maps through the file dialog.


## 2. Deployment
### 2.1 Custom Maps
A few sample maps are included, but it is possible to add custom maps in `/gui/maps` and load them in the Python program. Maps must be in [PNG](https://en.wikipedia.org/wiki/Portable_Network_Graphics) format and only consist of the following pixels with 100% opacity.
- ![#ffffff](https://placehold.it/15/ffffff/000000?text=+) `rgb(255, 255, 255)` &nbsp;&nbsp;- Floor and inner door.
- ![#000000](https://placehold.it/15/000000/000000?text=+) `rgb(0, 0, 0)` &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; - Wall.
- ![#00ff00](https://placehold.it/15/00ff00/000000?text=+) `rgb(0, 255, 0)` &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; - Exit door.
- ![#7f7f7f](https://placehold.it/15/7f7f7f/000000?text=+) `rgb(127, 127, 127)` &nbsp;&nbsp;- Out of bounds (outdoors).

Each pixel represents an area of 0.5m × 0.5m. Make sure that every room has an inner door and that there exists a valid path from every room to at least one exit door.

### 2.2. Run Program
Run the program from `/TheGreatEscape` with:
```
$ make run
```
### 2.3. Run Tests
Run the backend tests from `/TheGreatEscape` with:
```
$ make test
```

## 3. Created With
* [GitHub](https://github.com/) - Version control repository
* [Slack](https://slack.com/) - Team collaboration
* [Trello](https://trello.com/) - Scrum management
* [Photoshop](http://www.adobe.com/products/photoshop.html) - GUI development

## 4. Authors
* **Tobias Hiort**
* **Jenny Olsson**
* **Linn Löfquist**
* **Robin Larsson**
* **Elsa Slättegård**
* **Sinae Lee**
* **Axel Hallsenius**

## 6. License
This project is licensed under the MIT License. See the [LICENSE.md](https://github.com/TobiasHiort/TheGreatEscape/blob/master/LICENSE.md) file for details.

## 7. Repository Directory Tree (!TODO)
```
The Great Escape
│   .gitignore
|   thegreatescape.py (?)
│   LICENSE.md
|   README.md
|
└─── maps
|   |   ...
|
└─── gui
│   │   gui.py
│   │   utils.py
│   │
│   └─── fonts
│   |   │   ...
│   |
│   └─── gui
|   |   |   ...
|   |
│   └─── unit_tests
|       |   ...
|
└─── src
|   │   fire.go
|   |   gameMaster.go (camelcase?)
|   |   gotest.go (replace?)
|   |   main.go
|   |   map.go
|   |   pathfinder.go
|   |   people.go
|   |   print.go (needed later?)
|   |   queue.go
|   |
|   └─── unit_tests
|       |   ...
|
└─── tmp
    |   ...
```
