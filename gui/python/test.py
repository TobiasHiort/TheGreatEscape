import pygame
import sys
import os
import numpy
import math

import tkinter as tk
from tkinter import filedialog

from PIL import Image
from pygame import gfxdraw
from pygame.locals import *

# constants
GAME_RES = (1024, 768)
GAME_NAME = 'The Great Escape'

COLOR_WHITE = (255, 255, 255)
COLOR_BLACK = (0, 0, 0)
COLOR_GREEN = (0, 255, 0)
COLOR_RED = (255, 0, 0)
COLOR_BACKGROUND = (245, 245, 245)

PADDING_MAP = 10

# read image to matrix
mapImage = Image.open(os.path.join('maps\map1.png'))
mapRGBA = mapImage.load()
mapMatrix = numpy.zeros((mapImage.size[1], mapImage.size[0])) # (rows, column)

# game dimensions
if mapImage.size[0] < mapImage.size[1]:
    TILESIZE = math.floor((713)/mapImage.size[1])
else:
    TILESIZE = math.floor((907)/mapImage.size[0])


MAPWIDTH = mapImage.size[0] # number of columns in matrix
MAPHEIGHT = mapImage.size[1] # number of rows in matrix

# player start coords
playerPos = [2, 2] # remove later

# dictionary
colors = {
                0 : COLOR_BACKGROUND,
                1 : COLOR_BLACK,
                2 : COLOR_GREEN
          }

for row in range(mapImage.size[1]):
    for column in range(mapImage.size[0]):
        if mapRGBA[column, row] == (255, 255, 255, 255): # warning: mapRGBA has [column, row]
            mapMatrix[row][column] = 0
        elif mapRGBA[column, row] == (0, 0, 0, 255): # warning: mapRGBA has [column, row]
            mapMatrix[row][column] = 1

def addSprite(path, x, y):
    tmpSprite = pygame.sprite.Sprite()
    tmpSprite.image = pygame.image.load(os.path.join(path)).convert() # .convert?
    tmpSprite.rect = tmpSprite.image.get_rect()
    #tmpSprite.topleft = [0, 0]
    #displaySurface.blit(tmpSprite.image, tmpSprite.rect)
    displaySurface.blit(tmpSprite.image, (x, y))
    return tmpSprite


# init game
pygame.init()

icon = pygame.image.load("gui\window_icon.png")
pygame.display.set_icon(icon)

displaySurface = pygame.display.set_mode((GAME_RES)) # ,pygame.NOFRAME


#pygame.image.load(os.path.join('gui\window_icon.png')).convert()
#pygame.display.set_icon(pygame.image.load(os.path.join('gui\window_icon.png')).convert())

displaySurface.fill(COLOR_BACKGROUND)
pygame.display.set_caption(GAME_NAME)


myNewSurface = pygame.Surface((907-0*PADDING_MAP, 713-PADDING_MAP))
myNewSurface = myNewSurface.convert()
myNewSurface.fill(COLOR_BACKGROUND)
#myNewSurface.blit(icon, (0,0))

#displaySurface.blit(myNewSurface, (70, 70))

# fonts
FONT_ROBOTOMEDIUM = pygame.font.Font('Roboto-Medium.ttf', 18)

# import/display images
PLAYER = pygame.image.load(os.path.join('gui', 'player.png')).convert_alpha()

MENU_FADE = pygame.image.load(os.path.join('gui', 'menu_fade.png')).convert()
displaySurface.blit(MENU_FADE, (0, 45))

MENU_BACKGROUND = pygame.image.load(os.path.join('gui', 'menu_background.png')).convert()
displaySurface.blit(MENU_BACKGROUND, (0, 0))

MENU_RIGHT = pygame.image.load(os.path.join('gui', 'menu_right.png')).convert()
displaySurface.blit(MENU_RIGHT, (907, 45))

# layers ???
#layer1 = pygame.sprite.LayeredUpdates()
#BUTTON_SETTINGS_ACTIVE = addSprite('gui\settings_active.png', 203, 0)
#layer1.add(BUTTON_SETTINGS_ACTIVE)

# add button sprites
BUTTON_SIMULATION_ACTIVE = addSprite('gui\simulation_active.png', 0, 0)
#BUTTON_SIMULATION_BLANK = addSprite('gui\simulation_blank.png', 0, 0)
#BUTTON_SIMULATION_HOVER = addSprite('gui\simulation_hover.png', 0, 0)

#BUTTON_SETTINGS_ACTIVE = addSprite('gui\settings_active.png', 202, 0)
BUTTON_SETTINGS_BLANK = addSprite('gui\settings_blank.png', 202, 0)
#BUTTON_SETTINGS_HOVER = addSprite('gui\settings_hover.png', 202, 0)

#BUTTON_STATISTICS_ACTIVE = addSprite('gui\statistics_active.png', 382, 0)
BUTTON_STATISTICS_BLANK = addSprite('gui\statistics_blank.png', 382, 0)
#BUTTON_STATISTICS_HOVER = addSprite('gui\statistics_hover.png', 382, 0)

#displaySurface.blit(BUTTON_SETTINGS_ACTIVE, (0, 0))

# opening map file
file_opt = options = {}
options['defaultextension'] = '.png'
options['filetypes'] = [('PNG Map Files', '.png')]
options['initialdir'] = os.getcwd() + '\maps'
options['initialfile'] = 'mapXX.png'
#options['parent'] = root
options['title'] = 'Select Map'

# game loop
while True:
    for event in pygame.event.get():
        if event.type == QUIT:
            pygame.quit()
            sys.exit()
        elif event.type == KEYDOWN:
            if event.key == K_RIGHT:
                playerPos[0] += 1
            elif event.key == K_UP:
                root = tk.Tk()
                root.withdraw()
                file_path = filedialog.askopenfilename(**file_opt)
                filename_pos = file_path.rfind('/')+1 # position for filename
                print(file_path[filename_pos:]) # expand from here, probably need to create funcions for rendering before?

    #pygame.gfxdraw.aacircle(displaySurface, 500, 500, 30, COLOR_GREEN)
    #pygame.gfxdraw.filled_circle(displaySurface, 500, 500, 30, COLOR_GREEN)

    #pygame.draw.rect(displaySurface, COLOR_GREEN, (0, 0, 20, 20))
    #pygame.draw.circle(myNewSurface, COLOR_GREEN, (0, 0), 5)

    for row in range(MAPHEIGHT):
        for column in range(MAPWIDTH):
            pygame.draw.rect(myNewSurface, colors[mapMatrix[row][column]], (column*TILESIZE+((907-2*PADDING_MAP)/(2))-((MAPWIDTH*TILESIZE)/2)+PADDING_MAP, (row*TILESIZE+((713-1*PADDING_MAP)/(2))-((MAPHEIGHT*TILESIZE)/2)), TILESIZE, TILESIZE)) 
            myNewSurface.blit(PLAYER, (playerPos[0]*TILESIZE+15, playerPos[1]*TILESIZE+15))
            displaySurface.blit(myNewSurface, (0*PADDING_MAP, 55))

    # update display if not quitting
    pygame.display.flip() # .flip instead?




    # rows = 1


    # NOTES
    # - Up arrow is upload right now. Remove later for button press.
    # - Remove move player? However keep the mechanics.
    # - Determine if png is to large or small?
