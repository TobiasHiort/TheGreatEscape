#!/usr/bin/python3

import pygame
import sys
import os
import numpy
import math
import time
import tkinter as tk # replace
import subprocess
import doctest # read from txt, read docs

from pygame.locals import *
from tkinter import filedialog # remove?
from PIL import Image
#from pygame import gfxdraw # use later, AA

#doctest.testfile("unit_tests.txt")

# global constants
GAME_RES = (1024, 768)
GAME_NAME = 'The Great Escape'

PADDING_MAP = 10

active_tab_bools = [True, False, False] # [sim, settings, stats]
active_map_path = None # do not start with any map

COLOR_WHITE = (255, 255, 255)
COLOR_BLACK = (0, 0, 0)
COLOR_GREEN = (0, 255, 0)
COLOR_RED = (255, 0, 0)
COLOR_YELLOW = (255, 238, 67)
COLOR_BACKGROUND = (245, 245, 245)
COLOR_KEY = (127, 127, 127)

# init vars
#player_scale = 1.0 # beginning player scale
#current_frame = 0 # for simulation clock, not system # remove?

# dictionary for map matrix to color
colors = {
                0 : COLOR_WHITE,        # floor
                1 : COLOR_BLACK,        # wall
                2 : COLOR_RED,        # door
                3 : COLOR_BACKGROUND    # out of bounds
          }

def buildMap(path, mapSurface):
    """Returns mapSurface, mapMatrix, tilesize,
   mapwidth, mapheight after building map
   from png.

    More...
    """

    # read image to matrix
    mapImage = Image.open(os.path.join('maps', path))
    mapRGBA = mapImage.load()
    mapMatrix = numpy.zeros((mapImage.size[1], mapImage.size[0])) # (rows, column)

    # game dimensions
    if mapImage.size[0] < mapImage.size[1]:
        tilesize = math.floor((713)/mapImage.size[1])
    else:
        tilesize = math.floor((907)/mapImage.size[0])

    mapwidth = mapImage.size[0] # number of columns in matrix
    mapheight = mapImage.size[1] # number of rows in matrix

    # create map matrix dependent on tile type
    for row in range(mapheight):
        for column in range(mapwidth):
            if mapRGBA[column, row] == COLOR_WHITE + (255,): # warning: mapRGBA has [column, row]. RGBA
                mapMatrix[row][column] = 0
            elif mapRGBA[column, row] == COLOR_BLACK + (255,): # warning: mapRGBA has [column, row]. RGBA
                mapMatrix[row][column] = 1 # expand for more than floor and wall...
            elif mapRGBA[column, row] == COLOR_RED + (255,): # warning: mapRGBA has [column, row]. RGBA
                mapMatrix[row][column] = 2 # expand for more than floor and wall...
            elif mapRGBA[column, row] == COLOR_KEY + (255,): # warning: mapRGBA has [column, row]. RGBA
                mapMatrix[row][column] = 3
                # expand for more than floor and wall...

    # create the map with draw.rect on mapSurface
    for row in range(mapheight):
        for column in range(mapwidth):
            # warning, fix formula
            pygame.draw.rect(mapSurface, colors[mapMatrix[row][column]], (math.floor(column*tilesize+((907-2*PADDING_MAP)/(2))-((mapwidth*tilesize)/2)+PADDING_MAP), math.floor(row*tilesize+((713-1*PADDING_MAP)/(2))-((mapheight*tilesize)/2)), tilesize, tilesize)) 
    return mapSurface, mapMatrix, tilesize, mapwidth, mapheight

def drawPlayer(playerSurface, player_pos, tilesize, mapheight, mapwidth, player_scale):
    playerSurface.fill(COLOR_KEY) # remove last frame
    # draw player on simulation tab/mapsurface, remove second later, create funcion instead. fix the formula...
    for player in range(len(player_pos)):
        pygame.draw.circle(playerSurface, COLOR_GREEN, ((player_pos[player][0]*tilesize + math.floor(tilesize/2) + math.floor(((907-2*PADDING_MAP)/(2))-((mapwidth*tilesize)/2)+PADDING_MAP)), player_pos[player][1]*tilesize+round(tilesize/2) + round(0*tilesize+((713-1*PADDING_MAP)/(2))-((mapheight*tilesize)/2))), round((tilesize/2)*player_scale))
    return playerSurface

def placeText(surface, text, font, size, color, x, y):
    font = pygame.font.Font(font, size)
    surface.blit(font.render(text, True, color), (x, y))

def placeClockText(rmenuSurface, minutes, seconds):
    if len(minutes) == 2 and len(seconds) == 2:
        placeText(rmenuSurface, minutes, 'digital-7-mono.ttf', 45, COLOR_YELLOW, 8, 164)
        placeText(rmenuSurface, seconds, 'digital-7-mono.ttf', 45, COLOR_YELLOW, 71, 164)
    else:
        raise ValueError('Seconds and minutes must be of length 2')

def timeToString(seconds_input):
    mmss = divmod(seconds_input, 60)
    if mmss[0] < 10:
        minutes = "0" + str(mmss[0])
    
    else:
        minutes = str(mmss[0])
    
    if mmss[1] < 10:
        seconds = "0" + str(mmss[1])
    else:
        seconds = str(mmss[1])

    if mmss[0] > 99:
        return ("99", "++")
    else:
        return minutes, seconds

def mapSqm(mapMatrix):
    counter = 0.0
    for row in range(len(mapMatrix)):
        for column in range(len(mapMatrix[0])):
            if mapMatrix[row][column] == 0 or mapMatrix[row][column] == 2:
                counter += 0.25
    return counter

def mapExits(mapMatrix):
    counter = 0
    for row in range(len(mapMatrix)):
        for column in range(len(mapMatrix[0])):
            if mapMatrix[row][column] == 2:
                counter += 1
    return counter

def fileDialogInit():
    # for opening map file in tkinter
    file_opt = options = {}
    options['defaultextension'] = '.png'
    options['filetypes'] = [('PNG Map Files', '.png')]
    options['initialdir'] = os.getcwd() + '\maps'
    options['initialfile'] = 'mapXX.png'
    options['title'] = 'Select Map'
    return file_opt

def fileDialogPath():
    file_opt = {}
    root = tk.Tk()
    root.withdraw()
    file_path = filedialog.askopenfilename(**file_opt)
    filename_pos = file_path.rfind('/')+1 # position for filename
    active_map_path_tmp = file_path[filename_pos:]
    return active_map_path_tmp

# !temporary resets for player positions
def resetState():
    player_scale = 1.0
    current_frame = 0 # for simulation clock, not system
    current_time_float = 0.0 # for simulation clock, not system
    paused = True
    player_pos = []
    return player_scale, current_frame, current_time_float, paused, player_pos

def cursorBoxHit(mouse_x, mouse_y, x1, x2, y1, y2, tab):
    if (mouse_x > x1) and (mouse_x <= x2) and (mouse_y >= y1) and (mouse_y <= y2) and tab:
        return True
    else:
        return False

def buildButton(button):
    blank = pygame.image.load(os.path.join('gui', button +'_blank.png')).convert()
    hover = pygame.image.load(os.path.join('gui', button +'_hover.png')).convert()
    active = pygame.image.load(os.path.join('gui', button +'_active.png')).convert()
    return active, blank, hover

def createSurface(x, y):
    surface = pygame.Surface((x, y))
    surface = surface.convert()
    surface.fill(COLOR_BACKGROUND)
    surface.set_colorkey(COLOR_KEY)
    return surface

def setClock(rmenuSurface, seconds):
    minutes, seconds = timeToString(seconds)
    placeClockText(rmenuSurface, minutes, seconds)
    return rmenuSurface
