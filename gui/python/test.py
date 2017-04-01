import pygame
import sys
import os
import numpy
import math
import time

import tkinter as tk # replace
from tkinter import filedialog # remove?

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

PLAYER_SCALE = 0.5



# read image to matrix
mapImage = Image.open(os.path.join('maps', 'map1.png'))
mapRGBA = mapImage.load()
mapMatrix = numpy.zeros((mapImage.size[1], mapImage.size[0])) # (rows, column)

# game dimensions
if mapImage.size[0] < mapImage.size[1]:
    TILESIZE = math.floor((713)/mapImage.size[1])
else:
    TILESIZE = math.floor((907)/mapImage.size[0])

MAPWIDTH = mapImage.size[0] # number of columns in matrix
MAPHEIGHT = mapImage.size[1] # number of rows in matrix

# player start coords, remove later
playerPos = [0, 0]

# start with simulation tab (sim, settings, stats)
active_tab = [True, False, False]

# dictionary for map matrix to color
colors = {
                0 : COLOR_BACKGROUND, # floor
                1 : COLOR_BLACK, # wall
                2 : COLOR_GREEN # door, etc...
          }

# create map matrix dependent on tile type
for row in range(mapImage.size[1]):
    for column in range(mapImage.size[0]):
        if mapRGBA[column, row] == (255, 255, 255, 255): # warning: mapRGBA has [column, row]
            mapMatrix[row][column] = 0
        elif mapRGBA[column, row] == (0, 0, 0, 255): # warning: mapRGBA has [column, row]
            mapMatrix[row][column] = 1 # expand for more than floor and wall...

# init game
pygame.init()

# fonts
FONT_ROBOTOMEDIUM24 = pygame.font.Font('Roboto-Medium.ttf', 48)

# clock from start
clock = pygame.time.Clock()

# set window icon
icon = pygame.image.load(os.path.join('gui', 'window_icon.png'))
pygame.display.set_icon(icon)

# create the display surface, the overall screen size that will be rendered
displaySurface = pygame.display.set_mode((GAME_RES)) # ,pygame.NOFRAME
displaySurface.fill(COLOR_BACKGROUND)
pygame.display.set_caption(GAME_NAME)

# map/simulation surface
mapSurface = pygame.Surface((907-0*PADDING_MAP, 713-PADDING_MAP))
mapSurface = mapSurface.convert()
mapSurface.fill(COLOR_BACKGROUND)

# statistics surface
statisticsSurface = pygame.Surface((907-0*PADDING_MAP, 713-PADDING_MAP)) # x=1024?
statisticsSurface = statisticsSurface.convert()
statisticsSurface.fill(COLOR_RED)

settingsSurface = pygame.Surface((907-0*PADDING_MAP, 713-PADDING_MAP)) # x=1024?
settingsSurface = statisticsSurface.convert()
settingsSurface.fill(COLOR_GREEN)

# load and blit menu components to main surface (possibly remove blit, blit then in the game loop?)
MENU_FADE = pygame.image.load(os.path.join('gui', 'menu_fade.png')).convert()
displaySurface.blit(MENU_FADE, (0, 45))

MENU_BACKGROUND = pygame.image.load(os.path.join('gui', 'menu_background.png')).convert()
displaySurface.blit(MENU_BACKGROUND, (0, 0))

MENU_RIGHT = pygame.image.load(os.path.join('gui', 'menu_right.png')).convert()
displaySurface.blit(MENU_RIGHT, (907, 45))

# load buttons
BUTTON_SIMULATION_BLANK = pygame.image.load(os.path.join('gui', 'simulation_blank.png')).convert_alpha()
BUTTON_SIMULATION_HOVER = pygame.image.load(os.path.join('gui', 'simulation_hover.png')).convert_alpha()
BUTTON_SIMULATION_ACTIVE = pygame.image.load(os.path.join('gui', 'simulation_active.png')).convert_alpha() # maybe remove convert_alpha? No alphac.

BUTTON_SETTINGS_ACTIVE = pygame.image.load(os.path.join('gui', 'settings_active.png')).convert_alpha()
BUTTON_SETTINGS_HOVER = pygame.image.load(os.path.join('gui', 'settings_hover.png')).convert_alpha()
BUTTON_SETTINGS_BLANK = pygame.image.load(os.path.join('gui', 'settings_blank.png')).convert_alpha()

BUTTON_STATISTICS_ACTIVE = pygame.image.load(os.path.join('gui', 'statistics_active.png')).convert_alpha()
BUTTON_STATISTICS_HOVER = pygame.image.load(os.path.join('gui', 'statistics_hover.png')).convert_alpha()
BUTTON_STATISTICS_BLANK = pygame.image.load(os.path.join('gui', 'statistics_blank.png')).convert_alpha()

# for opening map file in tkinter
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
        # keyboard events
        elif event.type == KEYDOWN:
            # change player pos which will be rendered in the next frame. Remove later, but save mechanics.
            # [0] = column, [1] = row
            if event.key == K_d:
                playerPos[0] += 1
            if event.key == K_a:
                playerPos[0] -= 1
            if event.key == K_x:
                playerPos[1] += 1
            if event.key == K_w:
                playerPos[1] -= 1
            if event.key == K_q:
                playerPos[0] -= 1
                playerPos[1] -= 1
            if event.key == K_e:
                playerPos[0] += 1
                playerPos[1] -= 1
            if event.key == K_z:
                playerPos[0] -= 1
                playerPos[1] += 1
            if event.key == K_c:
                playerPos[0] += 1
                playerPos[1] += 1
            elif event.key == K_u:
                root = tk.Tk()
                root.withdraw()
                file_path = filedialog.askopenfilename(**file_opt)
                filename_pos = file_path.rfind('/')+1 # position for filename
                print(file_path[filename_pos:]) # expand from here, probably need to create funcions for rendering before?
        # mouse motion events
        elif event.type == MOUSEMOTION:
            mouse_x, mouse_y = event.pos
            if (mouse_x >= 0) and (mouse_x <= 202) and (mouse_y >= 0) and (mouse_y <= 45) and not(active_tab[0]):
                displaySurface.blit(BUTTON_SIMULATION_HOVER, (0,0))
            elif active_tab[0] == True:
                displaySurface.blit(BUTTON_SIMULATION_ACTIVE, (0,0))
            else:
            	displaySurface.blit(BUTTON_SIMULATION_BLANK, (0, 0))
            
            if (mouse_x > 202) and (mouse_x <= 382) and (mouse_y >= 0) and (mouse_y <= 45) and not(active_tab[1]):
                displaySurface.blit(BUTTON_SETTINGS_HOVER, (202,0))
            elif active_tab[1] == True:
                displaySurface.blit(BUTTON_SETTINGS_ACTIVE, (202,0))
            else:
            	displaySurface.blit(BUTTON_SETTINGS_BLANK, (202, 0))
            
            if (mouse_x > 382) and (mouse_x <= 575) and (mouse_y >= 0) and (mouse_y <= 45) and not(active_tab[2]):
                displaySurface.blit(BUTTON_STATISTICS_HOVER, (382,0))
            elif active_tab[2] == True:
                displaySurface.blit(BUTTON_STATISTICS_ACTIVE, (382,0))
            else:
            	displaySurface.blit(BUTTON_STATISTICS_BLANK, (382, 0))
        # mouse button events
        elif event.type == MOUSEBUTTONDOWN:
            # left click
            if event.button == 1:
                mouse_x, mouse_y = event.pos
                if (mouse_x >= 0) and (mouse_x <= 202) and (mouse_y >= 0) and (mouse_y <= 45):
                    print('simulation tab')
                    displaySurface.blit(BUTTON_SIMULATION_ACTIVE, (0, 0))
                    displaySurface.blit(BUTTON_SETTINGS_BLANK, (202, 0))
                    displaySurface.blit(BUTTON_STATISTICS_BLANK, (382, 0))
                    active_tab = [True, False, False]
                if (mouse_x > 202) and (mouse_x <= 382) and (mouse_y >= 0) and (mouse_y <= 45):
                    print('settings tab')
                    displaySurface.blit(BUTTON_SIMULATION_BLANK, (0, 0))
                    displaySurface.blit(BUTTON_SETTINGS_ACTIVE, (202, 0))
                    displaySurface.blit(BUTTON_STATISTICS_BLANK, (382, 0))
                    active_tab = [False, True, False]
                if (mouse_x > 382) and (mouse_x <= 575) and (mouse_y >= 0) and (mouse_y <= 45):
                    print('statistics tab')
                    displaySurface.blit(BUTTON_SIMULATION_BLANK, (0, 0))
                    displaySurface.blit(BUTTON_SETTINGS_BLANK, (202, 0))
                    displaySurface.blit(BUTTON_STATISTICS_ACTIVE, (382, 0))
                    active_tab = [False, False, True]

    # create the map with draw.rect and the player and then blit them
    for row in range(MAPHEIGHT):
	    for column in range(MAPWIDTH):
	        pygame.draw.rect(mapSurface, colors[mapMatrix[row][column]], (column*TILESIZE+((907-2*PADDING_MAP)/(2))-((MAPWIDTH*TILESIZE)/2)+PADDING_MAP, (row*TILESIZE+((713-1*PADDING_MAP)/(2))-((MAPHEIGHT*TILESIZE)/2)), TILESIZE, TILESIZE)) 

    # draw player on simulation tab/mapsurface
    pygame.draw.circle(mapSurface, COLOR_GREEN, ((playerPos[0]*TILESIZE + math.floor(TILESIZE/2) + math.floor(((907-2*PADDING_MAP)/(2))-((MAPWIDTH*TILESIZE)/2)+PADDING_MAP)), playerPos[1]*TILESIZE+round(TILESIZE/2) + round(0*TILESIZE+((713-1*PADDING_MAP)/(2))-((MAPHEIGHT*TILESIZE)/2))), round((TILESIZE/2)*PLAYER_SCALE))

    # settings page
    settings_text = FONT_ROBOTOMEDIUM24.render("settingsSurface placeholder", 0, COLOR_BLACK)
    settingsSurface.blit(settings_text, (200, 300))

    # statistics page
    statistics_text = FONT_ROBOTOMEDIUM24.render("statisticsSurface placeholder", 2, COLOR_BLACK)
    statisticsSurface.blit(statistics_text, (200, 300))

    if active_tab[0]:
        displaySurface.blit(mapSurface, (0*PADDING_MAP, 55))
    elif active_tab[1]:
    	displaySurface.blit(settingsSurface, (0*PADDING_MAP, 55))
    elif active_tab[2]:
    	displaySurface.blit(statisticsSurface, (0*PADDING_MAP, 55))
    else:
    	print('Error: No active tab')

    # just the messy width and height from above... remove later
    #print((0*TILESIZE+((907-2*PADDING_MAP)/(2))-((MAPWIDTH*TILESIZE)/2)+PADDING_MAP)) # width
    #print((0*TILESIZE+((713-1*PADDING_MAP)/(2))-((MAPHEIGHT*TILESIZE)/2))) # height

    # update display if not quitting
    pygame.display.flip() # .update(<surface_args>) instead?

    # fps, remove later
    clock.tick(50)
    pygame.display.set_caption("fps: " + str(clock.get_fps()))