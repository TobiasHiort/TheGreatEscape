#!/usr/bin/python3

import pygame
import sys
import os
import numpy
import math
import time
import subprocess
import doctest # read from txt, read docs
import random
import scipy.spatial as sp
import tkinter as tk
import copy
import json
import inflect
inflect = inflect.engine()
import matplotlib.backends.backend_agg as agg
import pylab
import matplotlib.pylab as plt
import matplotlib
matplotlib.use("Agg")
plt.rcParams["font.family"] = "Roboto"
plt.rcParams["font.weight"] = "medium"

from sys import getsizeof
from pygame.locals import *
from PIL import Image
from pygame import gfxdraw # use later, AA

from colorama import Fore, Back, Style
from subprocess import Popen, PIPE
from tkinter import filedialog
from colorama import init
init(autoreset=True)

#doctest.testfile("unit_tests.txt") # doctest

# global constants
GAME_RES = (1024, 768)
GAME_NAME = 'The Great Escape'

PADDING_MAP = 10

COLOR_WHITE = (255, 255, 255)
COLOR_BLACK = (0, 0, 0)
COLOR_BLUE = (0, 111, 162)
COLOR_GREEN = (0, 166, 56)
COLOR_RED = (162, 19, 24)
COLOR_RED_PNG = (255, 0, 0)
COLOR_RED_DEAD = (216, 0, 1)
COLOR_YELLOW = (255, 238, 67)
COLOR_BACKGROUND = (245, 245, 245)
COLOR_KEY = (127, 127, 127)
COLOR_GREY1 = (226, 226, 226) # lighter
COLOR_GREY2 = (145, 145, 145) # darker, org: 221
COLOR_GREY3 = (120, 120, 120)
COLOR_WARNING = (242, 219, 16)

# variables
active_tab_bools = [True, False, False] # [sim, settings, stats]
active_map_path = None # do not start with any map

map_error_visible = 5
text_rect = None
error_text_spacing = 34
error_text_x = 745
error_text_y = 242

# dictionary for map matrix colors
colors = {
                0 : COLOR_WHITE,        # floor
                1 : COLOR_BLACK,        # wall
                2 : COLOR_WHITE,        # door
                3 : COLOR_BACKGROUND    # out of bounds
          }

# k-dim tree for nearest RGB
valid_colors = [
					COLOR_WHITE,
					COLOR_BLACK,
					COLOR_RED_PNG,
					COLOR_KEY
			   ]
color_tree = sp.KDTree(valid_colors)



def buildMap(path, mapSurface):
    """Returns mapSurface, mapMatrix, tilesize,
   mapwidth, mapheight after building map
   from png.

    More...
    """
    map_error = []

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
                #raise ValueError('Invalid RGBA value(s) in map. ' + '(x:' + str(column+1) + ', y:' + str(row+1) + '), wrong RGBA: ' +  str(mapRGBA[column, row]))                      raise ValueError('Invalid RGBA value(s) in map. ' + '(x:' + str(column+1) + ', y:' + str(row+1) + '), wrong RGBA: ' +  str(mapRGBA[column, row]))
                #placeText(mapSurface, 'Invalid RGBA value(s) in map. ' + '(x:' + str(column+1) + ', y:' + str(row+1) + '), wrong RGBA: ' +  str(mapRGBA[column, row]), 'Roboto-Regular.ttf', 11, COLOR_RED, 0, 0)
                #print('hejj')
                map_error.append([column + 1, row + 1, mapRGBA[column, row]])
                #print(map_error)
    if map_error == []:
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
    return mapSurface, mapMatrix, tilesize, mapwidth, mapheight, map_error

def buildMiniMap(path, mapSurface, result_matrix, COLOR_HEAT_GRADIENT, heatMap_bool): # ~duplicate^
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
        tilesize = math.floor((344)/mapImage.size[1])
    else:
        tilesize = math.floor((495)/mapImage.size[0])

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
            #elif mapRGBA[column, row] == COLOR_KEY + (255,): # warning: mapRGBA has [column, row]. RGBA
            #    mapMatrix[row][column] = 3
            #else:
            #    raise ValueError('Invalid RGB value(s) in map: ' + '(x: ' + str(column+1) + ', y: ' + str(row+1) + '), ' + 'wrong RGBA: ' +  str(mapRGBA[column, row]))

    # for formula
    t = tilesize
    sh = 344 # map surface height
    sw = 495 # map surface width
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
            if heatMap_bool :
                if result_matrix[row][column] != 0:
                    pygame.draw.rect(mapSurface, COLOR_HEAT_GRADIENT[result_matrix[row][column] - 1],
                                     (math.floor(0.5 * (sw - w * t + 2 * t * column)),
                                      math.floor((sh - p)/2 - (h * t)/2 + t * row),
                                      tilesize, tilesize))            
            
    return mapSurface, mapMatrix, tilesize, mapwidth, mapheight

def calcScalingCircle(PADDING_MAP, tilesize, mapheight, mapwidth, width, height):
    """Description.

    More...
    """
    coord_x = math.floor(0.5 * (width - mapwidth * tilesize)) + math.floor(tilesize / 2)
    coord_y = math.floor(0.5 * (-mapheight * tilesize + height - PADDING_MAP)) + math.floor(tilesize / 2)
    radius_scale = math.floor(tilesize/2)
    return coord_x, coord_y, radius_scale

# create drawPlayer2 for this AA solution. Only for playerSurface (circles). createSurface biggest change!
def drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT):
    """Description.

    More...
    """
    playerSurface.fill((0, 0, 0, 0)) # remove last frame. Also black magic for AA gfxdraw blitting of players.
    
    survived = 0 # reset for each draw
    dead = 0 # reset for each draw

    for player in range(len(player_pos)):
        if player_pos[player][2] <= 0:
            dead += 1
        if player_pos[player][0] == 0 and player_pos[player][1] == 0 and player_pos[player][2] > 0:
            survived += 1
            pygame.gfxdraw.aacircle(playerSurface,
                            coord_x,
                            coord_y,
                                    math.floor(radius_scale*player_scale), COLOR_BLACK + (0,)) # round()?

            pygame.gfxdraw.filled_circle(playerSurface,
                            coord_x,
                            coord_y,
                            math.floor(radius_scale*player_scale), COLOR_BLACK + (0,)) # round()?
        else:           
            if player_pos[player][2] < 1:
                pygame.gfxdraw.aacircle(playerSurface,
                                coord_x + tilesize * player_pos[player][0],
                                coord_y + tilesize * player_pos[player][1],
                                        math.floor(radius_scale*player_scale), COLOR_RED_DEAD) # round()?

                pygame.gfxdraw.filled_circle(playerSurface,
                                coord_x + tilesize * player_pos[player][0],
                                coord_y + tilesize * player_pos[player][1],
                                math.floor(radius_scale*player_scale), COLOR_RED_DEAD) # round()?
                
            # black magic
            elif player_pos[player][2] > 0:
                pygame.gfxdraw.aacircle(playerSurface,
                                        coord_x + tilesize * player_pos[player][0],
                                        coord_y + tilesize * player_pos[player][1],
                                        math.floor(radius_scale*player_scale), COLOR_PLAYER_GRADIENT[player_pos[player][2]]) # round()?
                
                pygame.gfxdraw.filled_circle(playerSurface,
                                             coord_x + tilesize * player_pos[player][0],
                                             coord_y + tilesize * player_pos[player][1],
                                             math.floor(radius_scale*player_scale), COLOR_PLAYER_GRADIENT[player_pos[player][2]]) # round()?
        
    return playerSurface, survived, dead

def calcScalingSquare(PADDING_MAP, tilesize, mapheight, mapwidth, width, height):
    """Description.

    More...
    """
    #coord_x = math.floor(0.5 * (907 - mapwidth * tilesize)) + math.floor(tilesize / 2)
    #coord_y = math.floor(0.5 * (-mapheight * tilesize + 713 - PADDING_MAP)) + math.floor(tilesize / 2)
    coord_x = math.floor(width - mapwidth * tilesize)
    coord_y = math.floor((height - PADDING_MAP)/2 - (mapheight * tilesize)/2)

    #math.floor(0.5 * (coord_x + 2 * t * fire_pos[idx][0]))
    #math.floor(coord_y + t * fire_pos[idx][1])

    #LOL = 2 * t * fire_pos[idx][0]
    #LOL = t * fire_pos[idx][1]

    return coord_x, coord_y

def drawFire(fireSurface, fire_pos, tilesize, coord_x, coord_y, COLOR_FIRE_GRADIENT, frame):
    """Description.

    More...
    """

    if frame == 0:
        return drawWarnings(fireSurface, fire_pos, tilesize, coord_x, coord_y)
    # fireSurface.fill(COLOR_KEY) # remove last frame. Not needed?
    else:
        fireSurface.fill((0, 0, 0, 0))
        
        # create the map with draw.rect on mapSurface
        for idx in range(len(fire_pos)):
            if fire_pos[idx][2] < 100:
                pygame.draw.rect(fireSurface, COLOR_FIRE_GRADIENT[fire_pos[idx][2]] + (180,),
                                 (math.floor(0.5 * (coord_x + 2 * tilesize * fire_pos[idx][0])),
                                  math.floor(coord_y + tilesize * fire_pos[idx][1]),
                                  tilesize, tilesize))
            elif fire_pos[idx][2] >= 99:
                #if fire_pos[idx][2] == 2:
                pygame.draw.rect(fireSurface, COLOR_FIRE_GRADIENT[99] + (180,),
                                 (math.floor(0.5 * (coord_x + 2 * tilesize * fire_pos[idx][0])),
                                  math.floor(coord_y + tilesize * fire_pos[idx][1]),
                                  tilesize, tilesize))
        return fireSurface

def drawSmoke(smokeSurface, smoke_pos, tilesize, coord_x, coord_y, COLOR_SMOKE_GRADIENT):
    """Description.

    More...
    """
    # fireSurface.fill(COLOR_KEY) # remove last frame. Not needed?
    smokeSurface.fill((0, 0, 0, 0))
    # for formula

    # create the map with draw.rect on mapSurface
    for idx in range(len(smoke_pos)):
        #if 20 < smoke_pos[idx][2] and smoke_pos[idx][2] <= 100:
        if smoke_pos[idx][2] <= 100:
            pygame.draw.rect(smokeSurface, COLOR_SMOKE_GRADIENT[smoke_pos[idx][2]] + (100,),
                             (math.floor(0.5 * (coord_x + 2 * tilesize * smoke_pos[idx][0])),
                                math.floor(coord_y + tilesize * smoke_pos[idx][1]),
                             tilesize, tilesize))
        elif smoke_pos[idx][2] >= 99:
            #if fire_pos[idx][2] == 2:
            pygame.draw.rect(smokeSurface, COLOR_SMOKE_GRADIENT[99] + (100,),
                             (math.floor(0.5 * (coord_x + 2 * tilesize * smoke_pos[idx][0])),
                                math.floor(coord_y + tilesize * smoke_pos[idx][1]),
                             tilesize, tilesize))
    return smokeSurface

def drawWarnings(fireSurface, fire_pos, tilesize, coord_x, coord_y):
    fireSurface.fill((0, 0, 0, 0))
    
    # JENNY
    #for fire in fire_pos:
    #    pygame.gfxdraw.filled_trigon(fireSurface,
    #                          coord_x + tilesize * (fire[0] - 1),
    #                          coord_y + tilesize * (fire[1] + 1),
    #                          coord_x + tilesize * (fire[0] + 1),
    #                          coord_y + tilesize * (fire[1] + 1),
    #                          coord_x + tilesize * (fire[0]),
    #                          coord_y + tilesize * (fire[1] - 1),
    #                                 COLOR_WARNING) # round()?        
    #    pygame.gfxdraw.trigon(fireSurface,
    #                          coord_x + tilesize * (fire[0] - 1),
    #                          coord_y + tilesize * (fire[1] + 1),
    #                          coord_x + tilesize * (fire[0] + 1),
    #                          coord_y + tilesize * (fire[1] + 1),
    #                          coord_x + tilesize * (fire[0]),
    #                          coord_y + tilesize * (fire[1] - 1),
    #                          COLOR_RED_DEAD) # round()?

    # WORKS FOR SQUARE INSTEAD
    #for idx in range(len(fire_pos)):
    #    pygame.draw.rect(fireSurface, COLOR_RED,
    #                     (math.floor(0.5 * (coord_x + 2 * tilesize * fire_pos[idx][0])),
    #                        math.floor(coord_y + tilesize * fire_pos[idx][1]),
    #                     tilesize, tilesize))
    #    #if fire_pos[idx][2] == 2:
    #    pygame.draw.rect(fireSurface, COLOR_RED,
    #                     (math.floor(0.5 * (coord_x + 2 * tilesize * fire_pos[idx][0])),
    #                        math.floor(coord_y + tilesize * fire_pos[idx][1]),
    #                     tilesize, tilesize))
    
    for idx in range(len(fire_pos)):
        pygame.gfxdraw.filled_trigon(fireSurface,
                                     math.floor(0.5 * (coord_x + 2 * tilesize * fire_pos[idx][0])),
                                     math.floor(coord_y + tilesize * fire_pos[idx][1]) + tilesize,
                                     math.floor(0.5 * (coord_x + 2 * tilesize * fire_pos[idx][0])) + tilesize,
                                     math.floor(coord_y + tilesize * fire_pos[idx][1]) + tilesize,
                                     math.floor(0.5 * (coord_x + 2 * tilesize * fire_pos[idx][0])) + math.floor(tilesize/2),
                                     math.floor(coord_y + tilesize * fire_pos[idx][1]), COLOR_RED_DEAD)

    return fireSurface

def placeText(surface, text, font, color, x, y):
    """Description.

    More...
    """
    #font = pygame.font.Font(font, size)
    text_tmp = font.render(text, True, color, COLOR_WHITE)
    text_rect = text_tmp.get_rect()
    surface.blit(text_tmp, (x, y))
    return text_rect, x, y

def placeTextAlpha(surface, text, font, color, x, y):
    """Description.

    More...
    """
    #font = pygame.font.Font(font, size)
    text_tmp = font.render(text, True, color)
    text_rect = text_tmp.get_rect()
    surface.blit(text_tmp, (x, y))
    return text_rect, x, y

def placeCenterText(surface, text, font, color, width, y):
    #font = pygame.font.Font(font, size)
    text_tmp = font.render(text, True, color, COLOR_WHITE)
    text_rect = text_tmp.get_rect(center = (width / 2, y))
    surface.blit(text_tmp, text_rect)
    return text_rect, width, y

def placeCenterTextAlpha(surface, text, font, color, width, y):
    #font = pygame.font.Font(font, size)
    text_tmp = font.render(text, True, color)
    text_rect = text_tmp.get_rect(center = (width / 2, y))
    surface.blit(text_tmp, text_rect)
    return text_rect, width, y

def placeClockText(rmenuSurface, font, minutes, seconds):
    """Description.

    More...
    """
    if len(minutes) == 2 and len(seconds) == 2:
        placeTextAlpha(rmenuSurface, minutes, font, COLOR_YELLOW, 8, 249-17+1)
        placeTextAlpha(rmenuSurface, seconds, font, COLOR_YELLOW, 71, 249-17+1)
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

def setClock(rmenuSurface, font, seconds):
    """Description.

    More...
    """
    minutes, seconds = timeToString(seconds)
    placeClockText(rmenuSurface, font, minutes, seconds)
    return rmenuSurface

def mapSqm(mapMatrix):
    """Description.

    More...
    """
    counter = 0.00
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
    players_movement = []
    player_count = 0
    fire_movement = []
    fire_pos = []
    smoke_pos = []
    smoke_movement = []
    survived = 0
    dead = 0
    fire_percent = 0
    smoke_percent = 0
    return player_scale, current_frame, current_time_float, paused, player_pos, players_movement, player_count, fire_movement, fire_pos, survived, fire_percent, smoke_pos, smoke_movement, smoke_percent, dead

def cursorBoxHit(mouse_x, mouse_y, x1, x2, y1, y2, tab):
    """Description.

    More...
    """
    if (mouse_x >= x1) and (mouse_x <= x2) and (mouse_y >= y1) and (mouse_y <= y2) and tab:
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

def createSurface(x, y, alpha):
    """Description.

    More...
    """
    if alpha:
        surface = pygame.Surface((x, y))
        surface = surface.convert_alpha()
        surface.fill((0, 0, 0, 0))
    elif not alpha:
        surface = pygame.Surface((x, y))
        surface = surface.convert_alpha()
    else:
        raise ValueError('Argument alpha must be Bool')
    return surface

#def populateMap(mapMatrix, pop_percent):
def populateMap(mapMatrix, pop_percent, init_fires):
    """Description.

    More...
    """
    if pop_percent < 0 or pop_percent > 1:
        raise ValueError('pop_percent must be positive and <= 1')

    floor_coords = []
    fire_coords = []

    counter = 0
    for row in range(len(mapMatrix)):
        for column in range(len(mapMatrix[0])):
            if mapMatrix[row][column] == 0: # floor
                floor_coords.append([column, row])
                counter += 1

    pop_remove = round(counter - (pop_percent * counter)) # how many to remove from floor_coords

    random.seed() # remove '5' later, using same seed rn. () = system clock
    # delete pop_remove number of players
    for i in range(pop_remove):
        rand_player = random.randint(0, counter - 1)
        #print("rand_player: " + str(rand_player))
        if i < init_fires:
            fire_coords.append(floor_coords[rand_player])
            fire_coords[i].append(0, )
        del floor_coords[rand_player] # delete player
        counter -= 1
        player_count = len(floor_coords)
    
    for idx in range(player_count):
        floor_coords[idx].append(100,)

    return floor_coords, player_count, fire_coords

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

def rawPlot():
    """Description.

    More...
    """
    def f(t):
        return numpy.exp(-t) * numpy.cos(2*numpy.pi*-t)

    plot_x = 495
    plot_y = 344
    fig = plt.figure(figsize=[plot_x * 0.01, plot_y * 0.01], # Inches.
                       dpi=100,        # 100 dots per inch, so the resulting buffer is 395x344 pixels
                       )

    fig.set_size_inches(plot_x * 0.01, plot_y * 0.01)

    ax = fig.gca()

    plt.xlabel('xlabel')
    plt.ylabel('ylabel')
    plt.title("Title")
    plt.gcf().subplots_adjust(bottom=0.15, top=0.90, left=0.14, right=0.95)


    #l1, = ax.plot([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14], [1, 2, 4, 8, 15, 17, 18, 22, 23, 23, 24, 24, 25, 25])
    #l1, = ax.plot(numpy.sin(numpy.linspace(0, 2 * numpy.pi)), 'r-o')
    t1 = numpy.arange(0.0, 5.0, 0.10)
    t2 = numpy.arange(0.0, 5.0, 0.02)
    #l1, = ax.plot(t1, f(t1), 'bo', t2, f(t2), 'k')

    plt.figure(1)
    p1 = plt.subplot(211)
    l1, = plt.plot(t1, f(t1), 'o')
    p2 = plt.subplot(212)
    l2, = plt.plot(t2, numpy.cos(2*numpy.pi*t2), 'r--')

    l1.set_color((162/255, 19/255, 24/255))
    l2.set_color((0/255, 166/255, 56/255))

    #plt.xlabel('xlabel')
    #plt.ylabel('ylabel')
    #plt.title("Title")

    p1.spines['right'].set_visible(False)
    p1.spines['top'].set_visible(False)

    p2.spines['right'].set_visible(False)
    p2.spines['top'].set_visible(False)
    return fig

def rawPlot2(escaped, died):
    """Description.

    More...
    """
    def f(t):
        return numpy.exp(-t) * numpy.cos(2*numpy.pi*-t)

    plot_x = 495
    plot_y = 344
    fig = plt.figure(figsize=[plot_x * 0.01, plot_y * 0.01], # Inches.
                       dpi=100,        # 100 dots per inch, so the resulting buffer is 495x344 pixels
                       )

    fig.set_size_inches(plot_x * 0.01, plot_y * 0.01)

    ax = fig.gca()
    plt.gcf().subplots_adjust(bottom=0.15, top=0.90, left=0.12, right=0.95)

    #
    escaped_counter = 0
    escaped_list = []
    time_list = []
    t = 0
    for e in escaped:
        escaped_counter += e
        t += 0.1
        escaped_list.append(escaped_counter)
        time_list.append(t)
        
    l1, = ax.plot(time_list, escaped_list, label = 'Escaped') #, linestyle = '--')
       
    died_counter = 0
    died_list = []
    for d in died:
        died_counter += d
        died_list.append(died_counter)
        
    l2, = ax.plot(time_list, died_list, label = 'Died')
    #
   # print(escaped_list)
    
  #  l1, = ax.plot([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14], [1, 2, 4, 8, 15, 17, 18, 22, 23, 23, 24, 24, 25, 25], label = 'label1', linestyle = '--')
  #  l2, = ax.plot([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14], [1, 2, 4, 8, 15, 17, 18, 22, 23, 23, 24, 24, 25, 25][::-1], label = 'label2')
    #l1, = ax.plot(numpy.sin(numpy.linspace(0, 2 * numpy.pi)), 'r-o')
    #t1 = numpy.arange(0.0, 5.0, 0.10)
    #t2 = numpy.arange(0.0, 5.0, 0.02)
    #l1, = ax.plot(t1, f(t1), 'bo', t2, f(t2), 'k')
    #l1, = ax.plot(t1, f(t1), 'bo')

    #plt.figure(1)
    #p1 = plt.subplot(211)
    #l1, = plt.plot(t1, f(t1), 'o')
    #plt.subplot(212)
    #l2, = plt.plot(t2, numpy.cos(2*numpy.pi*t2), 'r--')

    l1.set_color((0/255, 166/255, 56/255))
    l2.set_color((162/255, 19/255, 24/255))

#    plt.xlabel('X label', fontname = "Roboto", fontweight = 'medium', fontsize = 11)
#    plt.ylabel('Y label', fontname = "Roboto", fontweight = 'medium', fontsize = 11)
#    plt.title("Title", fontname = "Roboto", fontweight = 'medium', fontsize = 16)
    plt.xlabel('Time [s]', fontname = "Roboto", fontweight = 'medium', fontsize = 11)
    plt.ylabel('People', fontname = "Roboto", fontweight = 'medium', fontsize = 11)
    plt.title("Escapes and deaths over time", fontname = "Roboto", fontweight = 'medium', fontsize = 16)

    ax.spines['right'].set_visible(False)
    ax.spines['top'].set_visible(False)

    plt.legend(loc=2, borderaxespad=0.5)

    return fig

def rawPlot3(stats):
    """Description.

    More...
    """
    def f(t):
        return numpy.exp(-t) * numpy.cos(2*numpy.pi*-t)

    plot_x = 200
    plot_y = 120
    fig = plt.figure(figsize=[plot_x * 0.01, plot_y * 0.01], # inches
                       dpi=100,        # 100 dots per inch, so the resulting buffer is 150x120 pixels
                       )

    fig.set_size_inches(plot_x * 0.01, plot_y * 0.01)

    ax = fig.gca()
    plt.gcf().subplots_adjust(bottom=0.15, top=0.90, left=0.12, right=0.95)

    #l1, = ax.plot([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14], [1, 2, 4, 8, 15, 17, 18, 22, 23, 23, 24, 24, 25, 25])
    #l2, = ax.plot([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14], [1, 2, 4, 8, 15, 17, 18, 22, 23, 23, 24, 24, 25, 25][::-1])
    #l1, = ax.plot(numpy.sin(numpy.linspace(0, 2 * numpy.pi)), 'r-o')
    #t1 = numpy.arange(0.0, 5.0, 0.10)
    #t2 = numpy.arange(0.0, 5.0, 0.02)
    #l1, = ax.plot(t1, f(t1), 'bo', t2, f(t2), 'k')
    #l1, = ax.plot(t1, f(t1), 'bo')

    #plt.figure(1)
    #p1 = plt.subplot(211)
    #l1, = plt.plot(t1, f(t1), 'o')
    #plt.subplot(212)
    #l2, = plt.plot(t2, numpy.cos(2*numpy.pi*t2), 'r--')

    #labels = 'Frogs', 'Hogs', 'Dogs', 'Logs'
    #sizes = [15, 30, 45, 10]

    #explode = (0, 0.1, 0, 0)  # only "explode" the 2nd slice (i.e. 'Hogs')
    #fig1, ax = plt.subplots()

    #ax.pie(sizes, labels=labels, autopct='%1.0f%%', shadow=True, startangle=90) # explode=explode
    #ax.axis('equal')  # Equal aspect ratio ensures that pie is drawn as a circle.

    # Data to plot
    labels = 'Dead', 'Harmed', 'Unharmed'
   # sizes = [30, 70]
    colors = [(162/255, 19/255, 24/255), (0/255, 100/255, 56/255), (0/255, 166/255, 56/255)]
    #explode = [0.1, 0]

    # Plot
   
   # lst = [stats[1], stats[0] - stats[2], stats[2]]
    lst = [stats[1], stats[2], stats[0] - stats[2]]
    patches, texts, autotexts = plt.pie(lst, labels=labels, colors=colors, autopct='%1.0f%%', shadow=True, startangle=45, labeldistance=1.25) # pctdistance=1.1
    texts[0].set_fontsize(9)
    texts[1].set_fontsize(9)
    texts[2].set_fontsize(9)

    plt.axis('equal')

    return fig

def rawPlot3b(stats):
    """Description.

    More...
    """
    def f(t):
        return numpy.exp(-t) * numpy.cos(2*numpy.pi*-t)

    plot_x = 200
    plot_y = 120
    fig = plt.figure(figsize=[plot_x * 0.01, plot_y * 0.01], # inches
                       dpi=100,        # 100 dots per inch, so the resulting buffer is 150x120 pixels
                       )

    fig.set_size_inches(plot_x * 0.01, plot_y * 0.01)

    ax = fig.gca()
    plt.gcf().subplots_adjust(bottom=0.15, top=0.90, left=0.12, right=0.95)

    # Data to plot
    labels = 'Smoke', 'Fire'
   # sizes = [30, 70]
    colors = [(124/255, 124/255, 124/255), (162/255, 19/255, 24/255)] 
    #explode = [0.1, 0]

    # Plot
    patches, texts, autotexts = plt.pie(stats, labels=labels, colors=colors, autopct='%1.0f%%', shadow=True, startangle=45, labeldistance=1.25) # pctdistance=1.1
    texts[0].set_fontsize(9)
    texts[1].set_fontsize(9)

    plt.axis('equal')

    return fig

def rawPlot4(smoke, fire, scale):  # movementlists
    """Description.

    More...
    """
    def f(t):
        return numpy.exp(-t) * numpy.cos(2*numpy.pi*-t)

    plot_x = 495
    plot_y = 344
    fig = plt.figure(figsize=[plot_x * 0.01, plot_y * 0.01], # Inches.
                       dpi=100,        # 100 dots per inch, so the resulting buffer is 495x344 pixels
                       )

    fig.set_size_inches(plot_x * 0.01, plot_y * 0.01)

    ax = fig.gca()
    plt.gcf().subplots_adjust(bottom=0.15, top=0.90, left=0.12, right=0.95)

    #
    #smoke_counter = 0
    smoke_list = []
    time_list = []
    t = 0
    for s in smoke:
     #   smoke_counter += s
        t += 0.1
        smoke_list.append(98*(len(s)/scale))
        time_list.append(t)
 #       
    l1, = ax.plot(time_list, smoke_list, label = 'Smoke') #, linestyle = '--')
       
    #fire_counter = 0
    fire_list = []
    for f in fire:
      #  fire_counter += d
        fire_list.append(98*(len(f)/scale))
        
    l2, = ax.plot(time_list, fire_list, label = 'Fire')

    l1.set_color((120/255, 120/255, 120/255))
    l2.set_color((162/255, 19/255, 24/255))

#    plt.xlabel('X label', fontname = "Roboto", fontweight = 'medium', fontsize = 11)
#    plt.ylabel('Y label', fontname = "Roboto", fontweight = 'medium', fontsize = 11)
#    plt.title("Title", fontname = "Roboto", fontweight = 'medium', fontsize = 16)
    plt.xlabel('Time [s]', fontname = "Roboto", fontweight = 'medium', fontsize = 11)
    plt.ylabel('% covered', fontname = "Roboto", fontweight = 'medium', fontsize = 11)
    plt.title("Area covered in smoke and fire", fontname = "Roboto", fontweight = 'medium', fontsize = 16)

    ax.spines['right'].set_visible(False)
    ax.spines['top'].set_visible(False)

    plt.legend(loc=2, borderaxespad=0.5)

    return fig


def tablePlot(stats):
    plot_x = 300
    plot_y = 320
    #    fig = table(cellText=clust_data,colLabels=collabel,loc='center')
    fig = plt.figure(figsize=[plot_x * 0.01, plot_y * 0.01], # inches
                     dpi = 100,        # 100 dots per inch, so the resulting buffer is 150x120 pixels
                     )
 
#    plt.figure()
    ax = fig.gca(frame_on = False)

    # remove x/y-axis
    ax.xaxis.set_visible(False)
    ax.yaxis.set_visible(False)

    # remove padding around graph
    # TODO
    size = [0.7] #[plot_x/100]
    
   # y=[1,2,3,4,5,4,3,2,1,1,1,1,1,1,1,1]
    #plt.plot([10,10,14,14,10],[2,4,4,2,2],'r')
    col_labels = ['Av. escape values']
    row_labels = ['Time','Health', 'Moved']
#    table_vals=[[11,12,13],[21,22,23],[31,32,33]]
    table_vals = [[str(round(stats[1][0], 2)) + ' [s]'], [str(stats[2][0]) + ' [%]'], [str(round(stats[4][0]/2, 2)) + ' [m]']]
  #  table_vals = [['1'], ['2'], ['hej3']]
    # the rectangle is where I want to place the table
    exit_table = plt.table(cellText = table_vals,
                           colWidths = size,                         
                           cellLoc = 'center',
                           rowLabels = row_labels,
                           colLabels = col_labels,
                           loc = 'center')
    
    #plt.text(10,10,'Table Title',size = 70)
    col_labels = ['Exit values']
    row_labels = ['Alive','Dead','Injured']
    table_vals = [[stats[0][0]], [stats[0][1]], [stats[0][2]]]
    # the rectangle is where I want to place the table
    the_table = plt.table(cellText = table_vals,
                          colWidths = size,
                          cellLoc = 'center',
                          rowLabels = row_labels,
                          colLabels = col_labels,
                          loc = 'upper center')
    
    col_labels = ['# of deaths from']
    row_labels = ['Smoke','Fire']
    table_vals=[[stats[3][0]], [stats[3][1]]] 
    # the rectangle is where I want to place the table
    the_table = plt.table(cellText = table_vals,
                          colWidths = size,
                          cellLoc = 'center',
                          rowLabels = row_labels,
                          colLabels = col_labels,
                          loc = 'lower center')
    
    #plt.plot(y)
    #plt.show()
    return fig

def rawPlotRender(fig):
    """Description.

    More...
    """
    canvas = agg.FigureCanvasAgg(fig)
    canvas.draw()
    renderer = canvas.get_renderer()
    raw_data = renderer.tostring_rgb()
    return raw_data

def fileDialogInit():

    """Description.

    More...

    """

    # for opening map file in tkinter

    file_opt = options = {}
    options['defaultextension'] = '.png'
    options['filetypes'] = [('PNG Map Files', '.png')]
    options['initialdir'] = os.getcwd() + '\maps'
    options['initialfile'] = 'mapXX.png'
    options['title'] = 'Select Map'
    return file_opt

def fileDialogPath():
    """Description.

    More...

    """

    file_opt = {}
    root = tk.Tk()
    root.withdraw()
    file_path = filedialog.askopenfilename(**file_opt)
    filename_pos = file_path.rfind('\\')+1 # position for filename
    active_map_path_tmp = file_path[filename_pos:]
    return active_map_path_tmp

def interpolateTuple(startcolor, goalcolor, steps):
    """Description.

    More...

    """
    R = startcolor[0]
    G = startcolor[1]
    B = startcolor[2]

    targetR = goalcolor[0]
    targetG = goalcolor[1]
    targetB = goalcolor[2]

    DiffR = targetR - R
    DiffG = targetG - G
    DiffB = targetB - B

    gradient_list = []

    for i in range(0, steps + 1):
        new_R = R + (DiffR * i / steps)
        new_G = G + (DiffG * i / steps)
        new_B = B + (DiffB * i / steps)

        gradient_list.append((new_R, new_G, new_B))
    return gradient_list

def pathToName(path):
    return path[path.rfind('/') + 1:-4]
  
def goThread(mapMatrix, player_pos, players_movement, fire_pos, fire_movement, smoke_pos, smoke_movement, child_pid):
    """Description.

    More...

    """
    print('\n')
    print(Fore.WHITE + Back.BLUE + Style.BRIGHT + ' '*9 + 'NEW SIM' + ' '*9)
    
    # export json map matrix
    map_matrixInt = copy.deepcopy(mapMatrix).astype(int)
    map_jsons = json.dumps(map_matrixInt.tolist())
    tofile = open('../tmp/mapfile.txt', 'w+')
    tofile.write(map_jsons)
    tofile.close()
    print(Fore.WHITE + Back.GREEN + Style.DIM + 'wrote ' + Back.GREEN + Style.BRIGHT + 'mapfile.txt' + ' '*8)

    # export json people position list
    player_pos_str = json.dumps(player_pos)#_tmp)
    tofile3 = open('../tmp/playerfile.txt', 'w+')
    tofile3.write(player_pos_str)
    tofile3.close()
    print(Fore.WHITE + Back.GREEN + Style.DIM + 'wrote ' + Back.GREEN + Style.BRIGHT + 'playerfile.txt' + ' '*5)

    fire_pos_str = json.dumps(fire_pos)#_tmp)
    tofile3 = open('../tmp/firefile.txt', 'w+')
    tofile3.write(fire_pos_str)
    tofile3.close()
    print(Fore.WHITE + Back.GREEN + Style.DIM + 'wrote ' + Back.GREEN + Style.BRIGHT + 'firefile.txt' + ' '*7)

    # spawn Go subprocess
    child = Popen('../src/main', stdin=subprocess.PIPE, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True)
    pid_json = json.dumps(child.pid)
    tofile7 = open('../tmp/pid.txt', 'w+')
    tofile7.write(pid_json)
    tofile7.close()
    print(Fore.WHITE + Back.GREEN + Style.DIM + 'wrote ' + Back.GREEN + Style.BRIGHT + 'pid.txt' + ' '*12)

    child.stdout.flush()
    child.stdin.flush()
    #print('\n')
    print(Fore.WHITE + Back.CYAN + Style.BRIGHT + 'go subprocess started' + ' '*4)

    # first people
    json_ppl_bytes = child.stdout.readline().rstrip('\n')
    player_pos = json.loads(json_ppl_bytes)
    for pos in player_pos:
        players_movement.append([pos])
        
    json_ppl = json.loads(json_ppl_bytes)
    
    # first fire
    fromgo_json_fire = child.stdout.readline().rstrip('\n')
    fire_pos = json.loads(fromgo_json_fire)
    fire_movement.append(fire_pos)
    #for pos in fire_pos:
    #    fire_movement.append([pos])
    json_fire = json.loads(fromgo_json_fire)

    # first smoke
    fromgo_json_smoke = child.stdout.readline().rstrip('\n')
    smoke_pos = json.loads(fromgo_json_smoke)
    smoke_movement.append(smoke_pos)
    #for pos in smoke_pos:
    #    smoke_movement.append([pos])
    json_smoke = json.loads(fromgo_json_smoke)

    #print('\n')
    print(Fore.WHITE + Back.YELLOW + Style.BRIGHT + 'calculating simulation...')
    go_time_pre = time.clock()
    while len(json_ppl_bytes) > 0: #fromgo_json != []:
        json_ppl = json.loads(json_ppl_bytes)
        json_fire = json.loads(fromgo_json_fire)
        json_smoke = json.loads(fromgo_json_smoke)
        
        json_ppl_bytes = child.stdout.readline().rstrip('\n')
        for i in range(len(json_ppl)):
            players_movement[i].append(json_ppl[i])
            
        fromgo_json_fire = child.stdout.readline().rstrip('\n')
        fire_movement.append(json_fire)

        fromgo_json_smoke = child.stdout.readline().rstrip('\n')
        smoke_movement.append(json_smoke)
    #print('\n')
    print(Fore.WHITE + Back.MAGENTA + Style.BRIGHT + 'go subprocess done and' + ' '*3)
    white_space_clock = 10 - len(str(roundSig(time.clock() - go_time_pre)))
    print(Fore.WHITE + Back.MAGENTA + Style.BRIGHT + 'terminated in ' + Back.MAGENTA + Style.BRIGHT + Fore.CYAN + str(roundSig(time.clock() - go_time_pre)) + 's' + ' '*white_space_clock)
    #print('\nGo subprocess done and \nterminated in ' + str(roundSig(time.clock() - go_time_pre)) + "s")

    os.remove('../tmp/mapfile.txt')
    #print('\n')
    print(Fore.WHITE + Back.RED + Style.DIM + 'removed ' + Back.RED + Style.BRIGHT + 'mapfile.txt' + ' '*6)
    os.remove('../tmp/playerfile.txt')
    print(Fore.WHITE + Back.RED + Style.DIM + 'removed ' + Back.RED + Style.BRIGHT + 'playerfile.txt' + ' '*3)
    os.remove('../tmp/firefile.txt')
    print(Fore.WHITE + Back.RED + Style.DIM + 'removed ' + Back.RED + Style.BRIGHT + 'firefile.txt' + ' '*5)

    with open('../tmp/pid.txt', 'a') as out:
        out.write(json.dumps(0))
    print(Fore.WHITE + Back.RED + Style.DIM + 'reset ' + Back.RED + Style.BRIGHT + '  pid.txt' + ' '*10)
    
    print(Fore.WHITE + Back.BLUE + Style.BRIGHT + ' '*11 + 'END' + ' '*11)

    child.stdout.flush()
    child.stdin.flush()
    
def colorSurface(surface, rgb):
    """Description.

    More...

    """
    arr = pygame.surfarray.pixels3d(surface)
    arr[:,:,0] = rgb[0]
    arr[:,:,1] = rgb[1]
    arr[:,:,2] = rgb[2]

def rotateCenter(image, angle):
    """Description.

    More...

    """
    orig_rect = image.get_rect()
    rot_image = pygame.transform.rotate(image, angle)
    rot_rect = orig_rect.copy()
    rot_rect.center = rot_image.get_rect().center
    rot_image = rot_image.subsurface(rot_rect).copy()
    return rot_image

def repairMap(path):
    img = Image.open(os.path.join('maps', pathToName(path) + '.png'))
    img_data = img.load()
    for y in range(img.size[1]):
        for x in range(img.size[0]):
            if img_data[x, y][3] != 255:
                img_data[x, y] = (img_data[x, y][0], img_data[x, y][1], img_data[x, y][2] , 255)
            if img_data[x, y][0:3] != COLOR_BLACK or img_data[x, y][0:3] != COLOR_WHITE or img_data[x, y][0:3] != COLOR_RED_PNG or img_data[x, y][0:3] != COLOR_KEY:
            	_, result = color_tree.query(img_data[x, y][0:3])
            	img_data[x, y] = valid_colors[result] + (255,)
    img.save(os.path.join('maps', pathToName(path)) + '.png', "PNG")
    return_path = str(pathToName)
    return return_path

def roundSig(x, sig=2):
	return round(x, sig - int(math.floor(math.log10(abs(x)))) - 1)

def showErrorPage(mapSurface, ERROR_BG, FONT_ROBOTOLIGHT_18, FONT_ROBOTOLIGHT_22, active_map_path_error, map_error):
    def showCoord(mapSurface, map_error, idx, error_text_x, error_text_y, error_text_spacing):
        _, _, text_y = placeCenterText(mapSurface, '(' + str(map_error[idx][0]) + ', ' + str(map_error[idx][1]) + ')', FONT_ROBOTOLIGHT_22, COLOR_BLACK, error_text_x, error_text_y + idx * error_text_spacing) # coord
        return text_y
    
    def showRGB(mapSurface, map_error, idx, error_text_x, error_text_y, error_text_spacing):
        placeCenterText(mapSurface, str(map_error[idx][2][0:3]), FONT_ROBOTOLIGHT_22, COLOR_BLACK, error_text_x + 333-15, error_text_y + idx * error_text_spacing) # rgb
    
    def showOpacity(mapSurface, map_error, idx, error_text_x, error_text_y, error_text_spacing):
        placeCenterText(mapSurface, str("{:.0%}".format(map_error[idx][2][3]/255)), FONT_ROBOTOLIGHT_22, COLOR_BLACK, error_text_x + 580-15, error_text_y + idx * error_text_spacing) # opacity

    mapSurface.blit(ERROR_BG, (260, 120))
    placeCenterText(mapSurface, pathToName(active_map_path_error) + '.png has ' + inflect.number_to_words(len(map_error)) + ' invalid pixels!', FONT_ROBOTOLIGHT_22, COLOR_BLACK, 1024, 158)

    map_error_length = 0
    for idx in range(map_error_visible):
        if idx < len(map_error):
            if map_error[idx][2][3] != 255 and (map_error[idx][2][0:3] != COLOR_WHITE and map_error[idx][2][0:3] != COLOR_BLACK and map_error[idx][2][0:3] != COLOR_RED_PNG and map_error[idx][2][0:3] != COLOR_KEY):
                #_, _, text_y = placeCenterText(mapSurface, '(' + str(map_error[idx][0]) + ', ' + str(map_error[idx][1]) + ')', 'Roboto-Light.ttf', 22, COLOR_BLACK, error_text_x, error_text_y + idx * error_text_spacing) # coord
                text_y = showCoord(mapSurface, map_error, idx, error_text_x, error_text_y, error_text_spacing)
                showRGB(mapSurface, map_error, idx, error_text_x, error_text_y, error_text_spacing)
                showOpacity(mapSurface, map_error, idx, error_text_x, error_text_y, error_text_spacing)
                pygame.gfxdraw.box(mapSurface, pygame.Rect(423, text_y-16, 33, 33), map_error[idx][2][0:3])
            elif map_error[idx][2][3] != 255:
                text_y = showCoord(mapSurface, map_error, idx, error_text_x, error_text_y, error_text_spacing)
                showOpacity(mapSurface, map_error, idx, error_text_x, error_text_y, error_text_spacing)
                pygame.gfxdraw.box(mapSurface, pygame.Rect(423, text_y-16, 33, 33), map_error[idx][2][0:3])
            elif map_error[idx][2][0:3] != COLOR_WHITE and map_error[idx][2][0:3] and COLOR_BLACK and map_error[idx][2][0:3] != COLOR_RED_PNG and map_error[idx][2][0:3] != COLOR_KEY:
                #_, _, text_y = placeCenterText(mapSurface, '(' + str(map_error[idx][0]) + ', ' + str(map_error[idx][1]) + ')', 'Roboto-Light.ttf', 22, COLOR_BLACK, error_text_x, error_text_y + idx * error_text_spacing) # coord
                text_y = showCoord(mapSurface, map_error, idx, error_text_x, error_text_y, error_text_spacing)
                showRGB(mapSurface, map_error, idx, error_text_x, error_text_y, error_text_spacing)
                pygame.gfxdraw.box(mapSurface, pygame.Rect(423, text_y-16, 33, 33), map_error[idx][2][0:3])
            map_error_length += 1
            if idx == map_error_visible - 1 and len(map_error) - map_error_length > 0:
                placeText(mapSurface, 'and ' + inflect.number_to_words(len(map_error) - map_error_length) + ' more...', FONT_ROBOTOLIGHT_18, COLOR_GREY2, 323, error_text_y + (idx + 1) * error_text_spacing - 15)
                map_error_length = 0
                break
    placeText(mapSurface, 'Repair and reload the map or try to', FONT_ROBOTOLIGHT_18, COLOR_BLACK, 322, 427)
    return mapSurface

def showDebugger(displaySurface, MENU_BACKGROUND, MENU_FADE, FONT_ROBOTOREGULAR_11, mapwidth, mapheight, active_tab_bools, pop_percent, paused, counter_seconds, current_time_float, player_pos, fps, active_map_path, tilesize, mouse_x, mouse_y, pipe_input, player_scale):
    displaySurface.blit(MENU_BACKGROUND, (570, 0)) # bs1
    displaySurface.blit(MENU_FADE, (-120, 45)) # bs^2

    placeText(displaySurface, "DEBUGGER", FONT_ROBOTOREGULAR_11, COLOR_BLACK, 570, 0)
    placeText(displaySurface, "+mapwidth: " + str(mapwidth) + "til" + " (" + str(mapwidth*0.5)+ "m)", FONT_ROBOTOREGULAR_11, COLOR_BLACK, 570, 10)
    placeText(displaySurface, "+mapheight: " + str(mapheight) + "til" + " (" + str(mapheight*0.5)+ "m)", FONT_ROBOTOREGULAR_11, COLOR_BLACK, 570, 20)
    placeText(displaySurface, "+tab: " + str(active_tab_bools), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 570, 30)

    # check out of bounds.
    # crashes the fuck out if there are players outside mapMatrix's bounds,
    # as long as Go provides correct data this should not happen
    
    #p_oob = None
    #p_oob_id = []
    #if player_pos != []:
    #    for player in range(len(players_movement)):
    #        if mapMatrix[player_pos[player][1]][player_pos[player][0]] == 1 or mapMatrix[player_pos[player][1]][player_pos[player][0]] == 3:
    #            p_oob_id.append(player)
    #if p_oob_id == []:
    #    p_oob = False
    #else:
    #    p_oob = True

    #placeText(displaySurface, "+p_pos: " + str(player_pos), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 31)
    #placeText(displaySurface, "+p_oob: " + str(p_oob), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 31)
    #placeText(displaySurface, "+oob_id: " + str(p_oob_id), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 43)
    
    placeText(displaySurface, "+pop_%: " + str(round(pop_percent, 2)), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 710, 31)
    placeText(displaySurface, "+paused: " + str(paused), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 710, 0)
    placeText(displaySurface, "+elapsed: " + str(counter_seconds), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 710, 11)
    placeText(displaySurface, "+frame_float: " + str(round(current_time_float, 2)), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 710, 21)

    placeText(displaySurface, "+p_scale: " + str(round(player_scale, 2)), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 810, 0)
    placeText(displaySurface, "+populated: " + str(player_pos != []), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 810, 11)
    placeText(displaySurface, "+fps: " + str(round(fps)), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 810, 21)
    placeText(displaySurface, "+file: " + str(active_map_path), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 810, 31)

    placeText(displaySurface, "+tilesize: " + str(tilesize), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 910, 0)
    placeText(displaySurface, "+mouse xy: " + str(mouse_x) + "," + str(mouse_y), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 910, 11)
    placeText(displaySurface, "+pipe_in: " + str(pipe_input), FONT_ROBOTOREGULAR_11, COLOR_BLACK, 910, 21)

def calcFPS(prev_time, target_fps, caption_bool):
    # fps calc, remove later
    curr_time = time.time() # so now we have time after processing
    diff = curr_time - prev_time # frame took this much time to process and render
    delay = max(1.0/target_fps - diff, 0) # if we finished early, wait the remaining time to desired fps, else wait 0 ms
    time.sleep(delay)
    fps = 1.0/(delay + diff) # fps is based on total time ("processing" diff time + "wasted" delay time)
    prev_time = curr_time
    if caption_bool:
        pygame.display.set_caption("{0}: {1:.2f}".format(GAME_NAME, fps))
    return prev_time, fps


def heatMap(player_movement, mapMatrix):

    result_matrix = copy.deepcopy(mapMatrix).astype(int).tolist()
   # map_matrix = map_matrixInt.tolist()
    
   # result_matrix = copy.deepcopy(map_matrix)
    for row in range(len(result_matrix)):
        for col in range(len(result_matrix[0])):
            result_matrix[row][col] = 0 #.append([0,])
            
    for player in player_movement:
        for frame in player:
            result_matrix[frame[1]][frame[0]] += 1

    result_matrix[0][0] = 0
    
    heat_max = 1 #max([sublist[-1] for sublist in result_matrix])
    for row in range(len(result_matrix)):
        for col in range(len(result_matrix[0])):
            if result_matrix[row][col] > heat_max:
                heat_max = result_matrix[row][col]

    
    #max(map(lambda x: x[-1], result_matrix))

    for row in range(len(result_matrix)):
        for col in range(len(result_matrix[0])):
            result_matrix[row][col] =  round(100*(result_matrix[row][col]/heat_max))
            
    return result_matrix

def findMapCoord(mouse_x, mouse_y, mapheight, mapwidth, t, tab):
    
    sw = 495
    sh = 344

    h = mapheight
    w = mapwidth
    
    def findX(column):
        return math.floor(0.5 * (sw - w * t + 2 * t * column))
    
    def findY(row):
        return  math.floor((sh - PADDING_MAP)/2 - (h * t)/2 + t * row)
  #  print([math.floor(((mouse_x - 517) - findX(0))/t), math.floor(((mouse_y - 60) - findY(0))/t), 0])
    return [math.floor(((mouse_x - 517) - findX(0))/t), math.floor(((mouse_y - 60) - findY(0))/t), 0]

   # print(findX(0))
    #print(findX(1))
    #print(getTile(mouse_x, mouse_y))

   
def printShortKeys():
    def printSK(s):
        ln = 29 - len(s)
        print(Fore.WHITE + Back.MAGENTA + Style.BRIGHT + s + ' '*ln)

    print(Fore.WHITE + Back.BLUE + Style.BRIGHT + ' '*10 + 'SHORTKEYS' + ' '*10) # 10->9 change desc. below?
    printSK('r: restart program')
    print(Fore.WHITE + Back.YELLOW + Style.BRIGHT + ' ' *8 + 'In Simulation'+ ' '*8)
    printSK('o: increase population')
    printSK('l: decrease population')
    printSK('+: increase number of fires')
    printSK('-: decrease number of fires')
    printSK('a: repopulate map')
    printSK('z: depopulate map')
    printSK('m: initiate simulation')
    printSK('s: run simulation')
    printSK('p: pause simulation')
    printSK('g: step simulation forwards')
    printSK('f: step simulation backwards')
    printSK('2: reset simulation')
    print(Fore.WHITE + Back.YELLOW + Style.BRIGHT + ' '*9 + 'In Settings'+ ' '*9)
    printSK('q: switch place ppl/fire')
    print(Fore.WHITE + Back.BLUE + Style.BRIGHT + ' '*10 + '         ' + ' '*10)
