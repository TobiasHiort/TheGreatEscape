
from tkinter import *

# remove some of these?
import pygame
import sys
import os
import numpy
import math
import time
import subprocess
import doctest # read from txt, read docs
import random

from sys import getsizeof

from pygame.locals import *


from PIL import Image
from pygame import gfxdraw # use later, AA

#doctest.testfile("unit_tests.txt") # doctest

# global constants
GAME_RES = (1024, 768)
GAME_NAME = 'The Great Escape'

PADDING_MAP = 10

active_tab_bools = [True, False, False] # [sim, settings, stats]
active_map_path = None # do not start with any map

COLOR_WHITE = (255, 255, 255)
COLOR_BLACK = (0, 0, 0)
COLOR_BLUE = (0, 111, 162)
COLOR_GREEN = (0, 166, 56)
COLOR_RED = (162, 19, 24)
COLOR_RED_PNG = (255, 0, 0)
COLOR_YELLOW = (255, 238, 67)
COLOR_BACKGROUND = (245, 245, 245)
COLOR_KEY = (127, 127, 127)
COLOR_GREY1 = (226, 226, 226) # lighter
COLOR_GREY2 = (145, 145, 145) # darker, org: 221

# dictionary for map matrix to color
colors = {
                0 : COLOR_WHITE,        # floor
                1 : COLOR_BLACK,        # wall
                2 : COLOR_BLUE,        # door
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
            elif mapRGBA[column, row] == COLOR_RED_PNG + (255,): # warning: mapRGBA has [column, row]. RGBA
                mapMatrix[row][column] = 2 # expand for more than floor and wall...
            elif mapRGBA[column, row] == COLOR_KEY + (255,): # warning: mapRGBA has [column, row]. RGBA
                mapMatrix[row][column] = 3
            else:
                raise ValueError('Invalid RGB value(s) in map: ' + '(x: ' + str(column+1) + ', y: ' + str(row+1) + '), ' + 'wrong RGBA: ' +  str(mapRGBA[column, row]))

    # for formula
    t = tilesize
    sh = 713 # map surface height
    sw = 907 # map surface width
    p = PADDING_MAP
    h = mapheight
    w = mapwidth

    # create the map with draw.rect on mapSurface
    for row in range(mapheight):
        for column in range(mapwidth):
            # black magic
            pygame.draw.rect(mapSurface, colors[mapMatrix[row][column]],
                             (math.floor(0.5 * (sw - w * t + 2 * t * column)),
                                math.floor((sh - p)/2 - (h * t)/2 + t * row),
                             tilesize, tilesize))
    return mapSurface, mapMatrix, tilesize, mapwidth, mapheight

def calcScaling(PADDING_MAP, tilesize, mapheight, mapwidth):
    """Description.

    More...
    """
    coord_x = math.floor(0.5 * (907 - mapwidth * tilesize)) + math.floor(tilesize / 2)
    coord_y = math.floor(0.5 * (-mapheight * tilesize + 713 - PADDING_MAP)) + math.floor(tilesize / 2)
    radius_scale = math.floor(tilesize/2)
    return coord_x, coord_y, radius_scale

# create drawPlayer2 for this AA solution. Only for playerSurface (circles). createSurface biggest change!
def drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale):
    """Description.

    More...
    """
    playerSurface.fill((0, 0, 0, 0), None, BLEND_RGBA_MULT) # remove last frame. Also black magic for AA gfxdraw blitting of players.

    for player in range(len(player_pos)):
        # black magic
        pygame.gfxdraw.aacircle(playerSurface,
                              coord_x + tilesize * player_pos[player][0],
                               coord_y + tilesize * player_pos[player][1],
                           math.floor(radius_scale*player_scale), COLOR_GREEN) # round()?

        pygame.gfxdraw.filled_circle(playerSurface,
                              coord_x + tilesize * player_pos[player][0],
                               coord_y + tilesize * player_pos[player][1],
                           math.floor(radius_scale*player_scale), COLOR_GREEN) # round()?
        

    return playerSurface

def drawFire(fireSurface, fire_pos, tilesize, mapheight, mapwidth):
    """Description.

    More...
    """
    fireSurface.fill(COLOR_KEY) # remove last frame. Not needed?

    # for formula
    t = tilesize
    sh = 713 # map surface height
    sw = 907 # map surface width
    p = PADDING_MAP
    h = mapheight
    w = mapwidth

    # create the map with draw.rect on mapSurface
    for idx in range(len(fire_pos)):
            if fire_pos[idx][2] == 1:
                pygame.draw.rect(fireSurface, COLOR_YELLOW,
                                 (math.floor(0.5 * (sw - w * t + 2 * t * fire_pos[idx][0])),
                                    math.floor((sh - p)/2 - (h * t)/2 + t * fire_pos[idx][1]),
                                 tilesize, tilesize))
            if fire_pos[idx][2] == 2:
                pygame.draw.rect(fireSurface, COLOR_RED_PNG,
                                 (math.floor(0.5 * (sw - w * t + 2 * t * fire_pos[idx][0])),
                                    math.floor((sh - p)/2 - (h * t)/2 + t * fire_pos[idx][1]),
                                 tilesize, tilesize))
            if fire_pos[idx][2] == 3:
                pygame.draw.rect(fireSurface, COLOR_RED,
                                 (math.floor(0.5 * (sw - w * t + 2 * t * fire_pos[idx][0])),
                                    math.floor((sh - p)/2 - (h * t)/2 + t * fire_pos[idx][1]),
                                 tilesize, tilesize))
    return fireSurface

def placeText(surface, text, font, size, color, x, y):
    """Description.

    More...
    """
    font = pygame.font.Font(font, size)
    surface.blit(font.render(text, True, color), (x, y))

def placeCenterText(surface, text, font, size, color, width, y):
    font = pygame.font.Font(font, size)
    text_tmp = font.render(text, True, color)
    text_rect = text_tmp.get_rect(center = (width / 2, y))
    surface.blit(text_tmp, text_rect)

def placeClockText(rmenuSurface, minutes, seconds):
    """Description.

    More...
    """
    if len(minutes) == 2 and len(seconds) == 2:
        placeText(rmenuSurface, minutes, 'digital-7-mono.ttf', 45, COLOR_YELLOW, 8, 249)
        placeText(rmenuSurface, seconds, 'digital-7-mono.ttf', 45, COLOR_YELLOW, 71, 249)
    else:
        raise ValueError('Seconds and minutes must be of length 2')

def timeToString(seconds_input):
    """Description.

    More...
    """
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

def setClock(rmenuSurface, seconds):
    """Description.

    More...
    """
    minutes, seconds = timeToString(seconds)
    placeClockText(rmenuSurface, minutes, seconds)
    return rmenuSurface

def mapSqm(mapMatrix):
    """Description.

    More...
    """
    counter = 0.0
    for row in range(len(mapMatrix)):
        for column in range(len(mapMatrix[0])):
            if mapMatrix[row][column] == 0 or mapMatrix[row][column] == 2:
                counter += 0.25 # 0.5m^2
    return counter

def mapExits(mapMatrix):
    """Description.

    More...
    """
    counter = 0
    for row in range(len(mapMatrix)):
        for column in range(len(mapMatrix[0])):
            if mapMatrix[row][column] == 2:
                counter += 1
    return counter

def resetState():
    """Description.

    More...
    """
    player_scale = 1.0
    current_frame = 0 # for simulation clock, not system
    current_time_float = 0.0 # for simulation clock, not system
    paused = True
    player_pos = []
    player_count = 0
    return player_scale, current_frame, current_time_float, paused, player_pos, player_count

def cursorBoxHit(mouse_x, mouse_y, x1, x2, y1, y2, tab):
    """Description.

    More...
    """
    if (mouse_x > x1) and (mouse_x <= x2) and (mouse_y >= y1) and (mouse_y <= y2) and tab:
        return True
    else:
        return False

def loadImage(folder, file):
    """Description.

    More...
    """
    image = pygame.image.load(os.path.join(folder, file)).convert()
    return image

def loadImageAlpha(folder, file):
    """Description.

    More...
    """
    image = pygame.image.load(os.path.join(folder, file)).convert_alpha()
    return image

def createSurface(x, y):
    """Description.

    More...
    """
    surface = pygame.Surface((x, y), SRCALPHA)
    surface = surface.convert_alpha()
    surface.fill(COLOR_KEY) # COLOR_BACKGROUND?
    surface.set_colorkey(COLOR_KEY)
    return surface

def populateMap(mapMatrix, pop_percent):
    """Description.

    More...
    """
    if pop_percent < 0 or pop_percent > 1:
        raise ValueError('pop_percent must be positive and <= 1')

    floor_coords = []
    counter = 0
    for row in range(len(mapMatrix)):
        for column in range(len(mapMatrix[0])):
            if mapMatrix[row][column] == 0: # floor
                floor_coords.append([column, row])
                counter += 1

    pop_remove = round(counter - (pop_percent * counter)) # how many to remove from floor_coords

    #print("counter: " + str(counter))
    #print("pop_remove: " + str(pop_remove))

    random.seed() # remove '5' later, using same seed rn. () = system clock
    # delete pop_remove number of players
    for _ in range(pop_remove):
        rand_player = random.randint(0, counter - 1)
        #print("rand_player: " + str(rand_player))
        del floor_coords[rand_player] # delete player
        counter -= 1
        player_count = len(floor_coords)
    return floor_coords, player_count


def makeItr(byte_limit, str1):
    """Description.

    More...
    """
    itr = math.floor(len(str1) / byte_limit)
    return itr

def splitPipeData(byte_limit, str1):
    """Description.

    More...
    """
    if len(str1) < byte_limit:
        return str1
    else:
        if math.floor(len(str1) % byte_limit) == 0:
            #itr = math.floor(len(str1) / byte_limit)
            itr = makeItr(byte_limit, str1)
        else:
            #itr = math.floor(len(str1) / byte_limit) + 1
            itr = makeItr(byte_limit, str1)
        tmp_str = []
        idx = 0
        for _ in range(itr):
            tmp_str.append(str1[idx:idx+byte_limit])
            idx += byte_limit
        return tmp_str
