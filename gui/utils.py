#!/usr/bin/python3
#from tkinter import *

# remove some of these?
import pygame
import sys
import os
import numpy
import math
#import wx
import time
import subprocess
import doctest # read from txt, read docs
import random

import tkinter as tk
from tkinter import filedialog

import matplotlib
matplotlib.use("Agg")
import matplotlib.backends.backend_agg as agg
import pylab
import matplotlib.pylab as plt
plt.rcParams["font.family"] = "Roboto"
plt.rcParams["font.weight"] = "medium"

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
                2 : COLOR_WHITE,        # door
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
                raise ValueError('Invalid RGBA value(s) in map. ' + '(x:' + str(column+1) + ', y:' + str(row+1) + '), wrong RGBA: ' +  str(mapRGBA[column, row]))

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

def buildMiniMap(path, mapSurface): # ~duplicate^
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
    playerSurface.fill((0, 0, 0, 0)) # remove last frame. Also black magic for AA gfxdraw blitting of players.

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
    # fireSurface.fill(COLOR_KEY) # remove last frame. Not needed?
    fireSurface.fill((0, 0, 0, 0))
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
                pygame.draw.rect(fireSurface, COLOR_YELLOW + (150,),
                                 (math.floor(0.5 * (sw - w * t + 2 * t * fire_pos[idx][0])),
                                    math.floor((sh - p)/2 - (h * t)/2 + t * fire_pos[idx][1]),
                                 tilesize, tilesize))
            if fire_pos[idx][2] == 2:
                pygame.draw.rect(fireSurface, COLOR_RED_PNG + (150,),
                                 (math.floor(0.5 * (sw - w * t + 2 * t * fire_pos[idx][0])),
                                    math.floor((sh - p)/2 - (h * t)/2 + t * fire_pos[idx][1]),
                                 tilesize, tilesize))
            if fire_pos[idx][2] == 3:
                pygame.draw.rect(fireSurface, COLOR_RED + (150,),
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
        placeText(rmenuSurface, minutes, 'digital-7-mono.ttf', 45, COLOR_YELLOW, 8, 249-17)
        placeText(rmenuSurface, seconds, 'digital-7-mono.ttf', 45, COLOR_YELLOW, 71, 249-17)
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

def createSurface(x, y):
    """Description.

    More...
    """
    #if alpha:
    #    surface = pygame.Surface((x, y), SRCALPHA)
        #surface = surface.convert_alpha() #?
    #elif not alpha:
    surface = pygame.Surface((x, y))
    #else:
    #    raise ValueError('Argument alpha must be Bool')
    surface = surface.convert_alpha() #?
    #surface.fill(COLOR_KEY) # COLOR_BACKGROUND?
    #surface.set_colorkey(COLOR_KEY)
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

def rawPlot2():
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

    l1, = ax.plot([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14], [1, 2, 4, 8, 15, 17, 18, 22, 23, 23, 24, 24, 25, 25], label = 'label1', linestyle = '--')
    l2, = ax.plot([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14], [1, 2, 4, 8, 15, 17, 18, 22, 23, 23, 24, 24, 25, 25][::-1], label = 'label2')
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

    l1.set_color((162/255, 19/255, 24/255))
    l2.set_color((0/255, 166/255, 56/255))

    plt.xlabel('X label', fontname = "Roboto", fontweight = 'medium', fontsize = 11)
    plt.ylabel('Y label', fontname = "Roboto", fontweight = 'medium', fontsize = 11)
    plt.title("Title", fontname = "Roboto", fontweight = 'medium', fontsize = 16)

    ax.spines['right'].set_visible(False)
    ax.spines['top'].set_visible(False)

    plt.legend(bbox_to_anchor=(1.00, 1), loc=1, borderaxespad=0.)

    return fig

def rawPlot3():
    """Description.

    More...
    """
    def f(t):
        return numpy.exp(-t) * numpy.cos(2*numpy.pi*-t)

    plot_x = 150
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
    labels = 'Dead', 'Survived'
    sizes = [30, 70]
    colors = [(162/255, 19/255, 24/255), (0/255, 166/255, 56/255)]
    explode = [0.1, 0]

    # Plot
    patches, texts, autotexts = plt.pie(sizes, labels=labels, explode=explode, colors=colors, autopct='%1.0f%%', shadow=True, startangle=45, labeldistance=1.25) # pctdistance=1.1
    texts[0].set_fontsize(9)
    texts[1].set_fontsize(9)

    plt.axis('equal')

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
    filename_pos = file_path.rfind('/')+1 # position for filename
    active_map_path_tmp = file_path[filename_pos:]
    return active_map_path_tmp
