#!/usr/bin/python3

import pygame
import sys
import os
import numpy
import math
import time
import tkinter as tk
import subprocess
import copy
import json
import signal
import colorsys
import inflect
import psutil
import _thread
import colorama
import matplotlib
matplotlib.use("Agg")
import matplotlib.backends.backend_agg as agg
import matplotlib.pylab as plt
import pylab

from tkinter import filedialog
from termcolor import colored
from sys import getsizeof
from subprocess import Popen, PIPE
from utils import *
from pygame.locals import *
from PIL import Image
from pygame import gfxdraw # use later, AA
from termcolor import colored, cprint
from colorama import Fore, Back, Style
from colorama import init
init(autoreset=True)

def restartProgram():
    python = sys.executable
    os.execl(python, python, * sys.argv)

# init game
pygame.init()

# set window icon and program name
icon = pygame.image.load(os.path.join('gui', 'window_icon.png'))
pygame.display.set_icon(icon)
pygame.display.set_caption(GAME_NAME)

# game time
TIMER1000 = USEREVENT + 1
TIMER100 = USEREVENT + 2
TIMER50 = USEREVENT + 3
pygame.time.set_timer(TIMER1000, 1000)
pygame.time.set_timer(TIMER100, 100)
pygame.time.set_timer(TIMER50, 50)
target_fps = 60
prev_time = time.time() # for fps

# variables
counter_seconds = 0 # counter for TIMER1000
current_frame = 0 # which time frame for movement, int
current_time_float = 0.0 # float time for accurate time frame measurement, right now 0.1s per time frame.
paused = True
player_scale = 1.0
player_count = 0
pop_percent = 0.01 # init as this later?

map_error = []

player_pos = [] # might use this as indicator to not populate instead of players_movement?
players_movement = []
fire_pos = []
fire_movement = []
smoke_pos = []		
smoke_movement = []

init_fires = 1

fire_percent = 0
smoke_percent = 0
survived = 0
dead = 0

surface_toggle = [True, True, True]

go_running = False
child_pid = None

throbber_angle = 0

plot_rendered = False
plot_x = 495
plot_y = 344

# how much data is sent in each pipe
byte_limit = 5

# debugger var inits, not needed later
active_map_path_tmp = None
tilesize = None
mapwidth = 0
mapheight = 0
pipe_input = None
mapMatrix = []
mouse_x = 0
mouse_y = 0
current_map_sqm = 0
current_map_exits = 0

# create gradients, can not do this in utils...
COLOR_PLAYER_GRADIENT = interpolateTuple(( 36, 102,   0), ( 66, 181,   0), 100) # 2 steps == len 3
COLOR_FIRE_GRADIENT   = interpolateTuple((253, 207,  88), (170,   6,   6), 100) # 2 steps == len 3
COLOR_SMOKE_GRADIENT  = interpolateTuple((254, 254, 254), (100, 100, 100), 100) # 2 steps == len 3

# create the display surface, the overall main screen size that will be rendered
displaySurface = pygame.display.set_mode((GAME_RES)) # FULLSCREEN, DOUBLEBUF?
displaySurface.fill(COLOR_BACKGROUND) # no color leaks in the beginning
displaySurface.set_colorkey(COLOR_KEY, pygame.RLEACCEL) # RLEACCEL unstable?

# create surfaces
mapSurface = createSurface(907, 713-PADDING_MAP, False)
minimapSurface = createSurface(495, 344, False)
playerSurface = createSurface(907, 713-PADDING_MAP, False)
fireSurface = createSurface(907, 713-PADDING_MAP, True)
smokeSurface = createSurface(907, 713-PADDING_MAP, True)
rmenuSurface = createSurface(115, 723, False)
statisticsSurface = createSurface(1024, 713, False)
settingsSurface = createSurface(1024, 713, False)
throbberSurface = createSurface(29, 29, True)

# load and blit menu components to main surface (possibly remove blit, blit then in the game loop?)
MENU_FADE = loadImage('gui', 'menu_fade.png')
displaySurface.blit(MENU_FADE, (0, 45)) # blit in game_loop?

MENU_BACKGROUND = loadImage('gui', 'menu_background.png')
displaySurface.blit(MENU_BACKGROUND, (0, 0))

MENU_RIGHT = loadImage('gui', 'menu_right.png')

# load buttons in init state
BUTTON_SIMULATION_ACTIVE = loadImage('gui', 'simulation_active.png')
BUTTON_SIMULATION_BLANK = loadImage('gui', 'simulation_blank.png')
BUTTON_SIMULATION_HOVER = loadImage('gui', 'simulation_hover.png')

BG_SETTINGS = loadImage('gui', 'settings_bg.png')
BUTTON_SETTINGS_ACTIVE = loadImage('gui', 'settings_active.png')
BUTTON_SETTINGS_BLANK = loadImage('gui', 'settings_blank.png')
BUTTON_SETTINGS_HOVER = loadImage('gui', 'settings_hover.png')

BG_STATISTICS = loadImage('gui', 'statistics_bg.png')
BUTTON_STATISTICS_ACTIVE = loadImage('gui', 'statistics_active.png')
BUTTON_STATISTICS_BLANK = loadImage('gui', 'statistics_blank.png')
BUTTON_STATISTICS_HOVER = loadImage('gui', 'statistics_hover.png')

BUTTON_RUN_BLANK = loadImage('gui', 'run_blank.png')
BUTTON_RUN_HOVER = loadImage('gui', 'run_hover.png')

BUTTON_RUN2_RED = loadImage('gui', 'paused.png')
BUTTON_RUN2_GREEN = loadImage('gui', 'playing.png')
BUTTON_RUN2_BW = loadImageAlpha('gui', 'bw.png')
BUTTON_RUN2_FBW = loadImageAlpha('gui', 'fbw.png')
BUTTON_RUN2_FFW = loadImageAlpha('gui', 'ffw.png')
BUTTON_RUN2_FW = loadImageAlpha('gui', 'fw.png')
BUTTON_RUN2_PAUSE = loadImageAlpha('gui', 'pause.png')
BUTTON_RUN2_PLAY = loadImageAlpha('gui', 'play.png')

BUTTON_UPLOAD_SMALL = loadImage('gui', 'upload_small.png')
BUTTON_UPLOAD_SMALL0 = loadImage('gui', 'upload_small0.png')

BUTTON_UPLOAD_LARGE = loadImage('gui', 'upload_large.png')
BUTTON_UPLOAD_LARGE0 = loadImage('gui', 'upload_large0.png')

BUTTON_REPAIR_BLANK = loadImage('gui', 'repair_blank.png')
BUTTON_REPAIR_HOVER = loadImage('gui', 'repair_hover.png')

TIMER_BACKGROUND = loadImage('gui', 'timer.png')

DIVIDER_LONG = loadImage('gui', 'divider_long.png')
DIVIDER_SHORT = loadImage('gui', 'divider_short.png')

BUTTON_SCALE = loadImage('gui', 'scale.png')
BUTTON_SCALE_PLUS = loadImage('gui', 'scale_plus.png')
BUTTON_SCALE_MINUS = loadImage('gui', 'scale_minus.png')

BUTTON_TIME_SPEED = loadImage('gui', 'time_speed.png')

BUTTON_PEOPLE = loadImage('gui', 'people.png')
BUTTON_FIRE = loadImage('gui', 'fire.png')
BUTTON_SMOKE = loadImage('gui', 'smoke.png')

ERROR_BG = loadImage('gui', 'error_bg.png')

THROBBER = loadImageAlpha('gui', 'throbber.png')

# init fonts for performance
FONT_DIGITAL7MONO_45 = pygame.font.Font('digital-7-mono.ttf', 45)

FONT_ROBOTOREGULAR_13 = pygame.font.Font('Roboto-Regular.ttf', 13)
FONT_ROBOTOREGULAR_14 = pygame.font.Font('Roboto-Regular.ttf', 14)
FONT_ROBOTOREGULAR_17 = pygame.font.Font('Roboto-Regular.ttf', 17)
FONT_ROBOTOREGULAR_20 = pygame.font.Font('Roboto-Regular.ttf', 20)
FONT_ROBOTOREGULAR_22 = pygame.font.Font('Roboto-Regular.ttf', 22)
FONT_ROBOTOREGULAR_24 = pygame.font.Font('Roboto-Regular.ttf', 24)
FONT_ROBOTOREGULAR_26 = pygame.font.Font('Roboto-Regular.ttf', 26)

FONT_ROBOTOMEDIUM_13 = pygame.font.Font('Roboto-Medium.ttf', 13)

FONT_ROBOTOLIGHT_18 = pygame.font.Font('Roboto-Light.ttf', 18)
FONT_ROBOTOLIGHT_22 = pygame.font.Font('Roboto-Light.ttf', 22)

# file dialog
file_opt = fileDialogInit()

# game loop
while True:
    # event logic
    for event in pygame.event.get():
        if event.type == QUIT:
            pygame.quit()
            sys.exit()
        # time events
        if event.type == TIMER50 and active_map_path and go_running:
            # throbber
            throbberSurface.fill((0, 0, 0, 0))
            throbberSurface.blit(THROBBER, (0,0))
            throbberSurface = rotateCenter(throbberSurface, throbber_angle)
            colorSurface(throbberSurface, COLOR_GREY3)
            throbber_angle -= 30
        if event.type == TIMER100: # 100ms per movement (or frame), meaning top speed of ((0.5*(1000/100))/1)*3.6 = 18 km/h
            if players_movement != []: # warning, handling unpopulated map
                if active_map_path is not None or active_map_path != "": # "" = cancel on choosing map
                    if current_frame < len(players_movement[0]) - 1 and not paused: # no more movement coords and not paused
                        current_frame += 1
                        current_time_float += 0.1

                        for player in range(len(player_pos)):
                            player_pos[player] = players_movement[player][current_frame]

                        if len(fire_movement) > current_frame:
                            fire_pos = fire_movement[current_frame]
			
                        if len(smoke_movement) > current_frame:
                            smoke_pos = smoke_movement[current_frame]                            

                            
                        # pause if last frame, problem with go's pipe?
                        if current_frame == len(players_movement[0]) - 1:
                            paused = True

                        playerSurface, survived, dead = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)
                        fireSurface = drawFire(fireSurface, fire_pos, tilesize, coord_x_square, coord_y_square, COLOR_FIRE_GRADIENT)
                        smokeSurface = drawSmoke(smokeSurface, smoke_pos, tilesize, coord_x_square, coord_y_square, COLOR_SMOKE_GRADIENT)
                        
                        fire_percent = len(fire_pos) / (current_map_sqm * 4) # shady?
                        smoke_percent = len(smoke_pos) / (current_map_sqm * 4) # shady?
        if event.type == TIMER1000: # just specific for clock animation, 10*100ms below instead?
            counter_seconds += 1
        # keyboard events, later move to to mouse click event
        elif event.type == KEYDOWN:
            if event.key == K_r:
                restartProgram()
            elif event.key == K_1: # depopulate
                    pygame.quit()
                    sys.exit()
            if active_tab_bools[0] and active_map_path is not None: # do not add time/pos if no map
                if event.key == K_g and paused and players_movement != []: # forwards player movement from players_movement, move later to timed game event
                        if current_frame < len(players_movement[0]) - 1: # no (more) movement tuples
                            current_frame += 1
                            current_time_float += 0.1
                            
                            for player in range(len(player_pos)):
                                player_pos[player] = players_movement[player][current_frame]

                            fire_pos = fire_movement[current_frame]
                            smoke_pos = smoke_movement[current_frame]
                            playerSurface, survived, dead = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)
                            fireSurface = drawFire(fireSurface, fire_pos, tilesize, coord_x_square, coord_y_square, COLOR_FIRE_GRADIENT)
                            smokeSurface = drawSmoke(smokeSurface, smoke_pos, tilesize, coord_x_square, coord_y_square, COLOR_SMOKE_GRADIENT)

                            fire_percent = len(fire_pos) / (current_map_sqm * 4) # shady?
                            smoke_percent = len(smoke_pos) / (current_map_sqm * 4) # shady?


                elif event.key == K_f and paused and player_pos != []: # backwards player movement from players_movement, move later to timed game event
                        if current_frame > 0: # no (more) movement tuples
                            current_frame -= 1
                            current_time_float -= 0.1
                            for player in range(len(player_pos)):
                                player_pos[player] = players_movement[player][current_frame]

                            fire_pos = fire_movement[current_frame]
                            playerSurface, survived, dead = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)
                            fireSurface = drawFire(fireSurface, fire_pos, tilesize, coord_x_square, coord_y_square, COLOR_FIRE_GRADIENT)
                            smoke_pos = smoke_movement[current_frame]
                            smokeSurface = drawSmoke(smokeSurface, smoke_pos, tilesize, coord_x_square, coord_y_square, COLOR_SMOKE_GRADIENT)

                            fire_percent = len(fire_pos) / (current_map_sqm * 4) # shady?
                            smoke_percent = len(smoke_pos) / (current_map_sqm * 4) # shady?
                elif event.key == K_2 and paused:
                    current_frame = 0
                    current_time_float = 0.0
                    playerSurface, survived, dead = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)
                elif event.key == K_m and paused and current_frame == 0:
                    if not go_running:
                        _thread.start_new_thread(goThread, (mapMatrix, player_pos, players_movement, fire_pos, fire_movement, smoke_pos, smoke_movement, child_pid))
                        go_running = True
                        pygame.time.wait(20)
                elif event.key == K_s and paused and player_pos != []:                                  
                    if players_movement != [] and current_frame < len(players_movement[0]) - 1:  # do not start time frame clock if not pupulated.
                        paused = False
                elif event.key == K_p: # for use with cursorHitBox
                    paused = True # if paused == True -> False?
                elif event.key == K_o: # for use with cursorHitBox
                    if pop_percent < 0.8:
                        pop_percent *= 1.25
                elif event.key == K_l: # for use with cursorHitBox
                    if pop_percent > 0.1:
                        pop_percent *= 0.8
                elif event.key == K_PLUS:
                    init_fires += 1
                elif event.key == K_MINUS:
                    if init_fires > 1:
                        init_fires -= 1
                elif event.key == K_5: # for use with cursorHitBox
                    if surface_toggle[0]:
                        surface_toggle[0] = False
                    else:
                        surface_toggle[0] = True
                elif event.key == K_6: # for use with cursorHitBox
                    if surface_toggle[1]:
                        surface_toggle[1] = False
                    else:
                        surface_toggle[1] = True
                elif event.key == K_7: # for use with cursorHitBox
                    if surface_toggle[2]:
                        surface_toggle[2] = False
                    else:
                        surface_toggle[2] = True                   
                elif event.key == K_a: # populate, warning. use after randomizing init pos
                    paused = True
                    player_scale = 1.0
                    # function of this? maybe scrap for direction movement instead
                    player_count = len(player_pos)
                    if current_frame == 0:
                        player_pos, player_count, fire_pos = populateMap(mapMatrix, pop_percent, init_fires)
                        playerSurface, survived, dead = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)
                        fireSurface = drawWarnings(fireSurface, fire_pos, tilesize, coord_x_square, coord_y_square)
                        #fireSurface = drawFire(fireSurface, fire_pos, tilesize, mapheight, mapwidth, COLOR_FIRE_GRADIENT)
                    else:
                        print('Depop first')
                elif event.key == K_z: # depopulate
                    _, current_frame, current_time_float, paused, player_pos, players_movement, player_count, fire_movement, fire_pos, survived, fire_percent, smoke_pos, smoke_movement, smoke_percent, dead = resetState()
                    playerSurface, survived, dead = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)
                    fireSurface = drawFire(fireSurface, fire_pos, tilesize, coord_x_square, coord_y_square, COLOR_FIRE_GRADIENT)
                    smokeSurface = drawSmoke(smokeSurface, smoke_pos, tilesize, coord_x_square, coord_y_square, COLOR_SMOKE_GRADIENT)
        # mouse motion events (hovers), only for tab buttons on displaySurface. Blit in render logic for others.
        elif event.type == MOUSEMOTION:
            mouse_x, mouse_y = event.pos
            # simulation button
            if cursorBoxHit(mouse_x, mouse_y, 0, 202, 0, 45, not(active_tab_bools[0])):
                displaySurface.blit(BUTTON_SIMULATION_HOVER, (0,0))
            elif active_tab_bools[0]:
                displaySurface.blit(BUTTON_SIMULATION_ACTIVE, (0,0))
            else:
                displaySurface.blit(BUTTON_SIMULATION_BLANK, (0, 0))
            # settings button
            if cursorBoxHit(mouse_x, mouse_y, 202, 382, 0, 45, not(active_tab_bools[1])):
                displaySurface.blit(BUTTON_SETTINGS_HOVER, (202,0))
            elif active_tab_bools[1]:
                displaySurface.blit(BUTTON_SETTINGS_ACTIVE, (202,0))
            else:
                displaySurface.blit(BUTTON_SETTINGS_BLANK, (202, 0))
            # statistics button
            if cursorBoxHit(mouse_x, mouse_y, 383, 575, 0, 45, not(active_tab_bools[2])):
                displaySurface.blit(BUTTON_STATISTICS_HOVER, (382,0))
            elif active_tab_bools[2]:
                displaySurface.blit(BUTTON_STATISTICS_ACTIVE, (382,0))
            else:
                displaySurface.blit(BUTTON_STATISTICS_BLANK, (382, 0))
        # mouse button events (clicks)
        elif event.type == MOUSEBUTTONDOWN: # import as function?
            # left click
            if event.button == 1:
                mouse_x, mouse_y = event.pos
                # simulation button
                if cursorBoxHit(mouse_x, mouse_y, 0, 202, 0, 45, not(active_tab_bools[0])):
                    displaySurface.fill(COLOR_BACKGROUND)
                    displaySurface.blit(MENU_FADE, (0, 45))
                    displaySurface.blit(MENU_BACKGROUND, (0, 0))
                    displaySurface.blit(BUTTON_SETTINGS_BLANK, (202, 0))
                    displaySurface.blit(BUTTON_STATISTICS_BLANK, (382, 0))
                    displaySurface.blit(BUTTON_SIMULATION_ACTIVE, (0, 0))
                    active_tab_bools = [True, False, False]
                # settings button
                elif cursorBoxHit(mouse_x, mouse_y, 202, 382, 0, 45, not(active_tab_bools[1])):
                    displaySurface.blit(BUTTON_SIMULATION_BLANK, (0, 0))
                    displaySurface.blit(BUTTON_STATISTICS_BLANK, (382, 0))
                    displaySurface.blit(BUTTON_SETTINGS_ACTIVE, (202, 0))
                    active_tab_bools = [False, True, False]
                # statistics button
                elif cursorBoxHit(mouse_x, mouse_y, 383, 575, 0, 45, not(active_tab_bools[2])):
                    displaySurface.blit(BUTTON_SIMULATION_BLANK, (0, 0))
                    displaySurface.blit(BUTTON_SETTINGS_BLANK, (202, 0))
                    displaySurface.blit(BUTTON_STATISTICS_ACTIVE, (382, 0))
                    active_tab_bools = [False, False, True]
                # upload button routine startup
                if cursorBoxHit(mouse_x, mouse_y, 450, 574, 335, 459, active_tab_bools[0]) and active_map_path is None and not map_error:
                    active_map_path_tmp = fileDialogPath()
                    if active_map_path_tmp != "": #and active_map_path != "/":
                        active_map_path_error = active_map_path_tmp
                        active_map_path = active_map_path_tmp # (2/2)fixed bug for exiting folder window, not sure why tmp is needed
                        # reset state.
                        player_scale, current_frame, current_time_float, paused, player_pos, players_movement, player_count, fire_movement, fire_pos, survived, fire_percent, smoke_pos, smoke_movement, smoke_percent, dead = resetState()
                        # clear old map and players
                        mapSurface.fill(COLOR_BACKGROUND)
                        playerSurface.fill(COLOR_KEY)
                        # build new map
                        mapSurface, mapMatrix, tilesize, mapwidth, mapheight, map_error = buildMap(active_map_path, mapSurface)
                        if map_error != []:
                            active_map_path = None

                        if map_error == []: # dont draw players and calculate if error(s)
                            #mapSurface.set_alpha(0)
                            #opacity3 = 0
                            # precalc (better performance) for scaling formula
                            coord_x_circle, coord_y_circle, radius_scale = calcScalingCircle(PADDING_MAP, tilesize, mapheight, mapwidth)
                            coord_x_square, coord_y_square = calcScalingSquare(PADDING_MAP, tilesize, mapheight, mapwidth)

                            # compute sqm/exits
                            current_map_sqm = mapSqm(mapMatrix)
                            current_map_exits = mapExits(mapMatrix)

                            player_pos, player_count, fire_pos = populateMap(mapMatrix, pop_percent, init_fires)
                            players_movement = []

                            playerSurface, _, _ = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)
                # upload button map error
                if cursorBoxHit(mouse_x, mouse_y, 450, 574, 335+250, 459+250, active_tab_bools[0]) and map_error:
                    active_map_path_tmp = fileDialogPath()
                    if active_map_path_tmp != "": #and active_map_path != "/":
                        active_map_path_error = active_map_path_tmp
                        active_map_path = active_map_path_tmp # (2/2)fixed bug for exiting folder window, not sure why tmp is needed
                        # reset state.
                        player_scale, current_frame, current_time_float, paused, player_pos, players_movement, player_count, fire_movement, fire_pos, survived, fire_percent, smoke_pos, smoke_movement, smoke_percent, dead = resetState()
                        # clear old map and players
                        mapSurface.fill(COLOR_BACKGROUND)
                        playerSurface.fill(COLOR_KEY)
                        # build new map
                        mapSurface, mapMatrix, tilesize, mapwidth, mapheight, map_error = buildMap(active_map_path, mapSurface)
                        if map_error != []:
                            active_map_path = None
                        if map_error == []: # dont draw players and calculate if error(s)
                            #mapSurface.set_alpha(0)
                            #opacity3 = 0
                            # precalc (better performance) for scaling formula
                            coord_x_circle, coord_y_circle, radius_scale = calcScalingCircle(PADDING_MAP, tilesize, mapheight, mapwidth)
                            coord_x_square, coord_y_square = calcScalingSquare(PADDING_MAP, tilesize, mapheight, mapwidth)

                            # compute sqm/exits
                            current_map_sqm = mapSqm(mapMatrix)
                            current_map_exits = mapExits(mapMatrix)

                            player_pos, player_count, fire_pos = populateMap(mapMatrix, pop_percent, init_fires)
                            players_movement = []

                            playerSurface, _, _ = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)
                # upload button routine rmenu
                if cursorBoxHit(mouse_x, mouse_y, 937, 999, 685, 747, active_tab_bools[0]) and active_map_path is not None:
                    active_map_path_tmp = fileDialogPath()
                    if active_map_path_tmp != "": #and active_map_path != "/":
                        active_map_path_error = active_map_path_tmp
                        active_map_path = active_map_path_tmp # (2/2)fixed bug for exiting folder window, not sure why tmp is needed
                        # reset state.
                        player_scale, current_frame, current_time_float, paused, player_pos, players_movement, player_count, fire_movement, fire_pos, survived, fire_percent, smoke_pos, smoke_movement, smoke_percent, dead = resetState()
                        # clear old map and players
                        mapSurface.fill(COLOR_BACKGROUND)
                        playerSurface.fill(COLOR_KEY)
                        fireSurface.fill((0, 0, 0, 0))
                        smokeSurface.fill((0, 0, 0, 0))
                        # build new map
                        mapSurface, mapMatrix, tilesize, mapwidth, mapheight, map_error = buildMap(active_map_path, mapSurface)
                        if map_error != []:
                            active_map_path = None

                        elif map_error == []: # dont draw players and calculate if error(s)
                            # precalc (better performance) for scaling formula
                            coord_x_circle, coord_y_circle, radius_scale = calcScalingCircle(PADDING_MAP, tilesize, mapheight, mapwidth)
                            coord_x_square, coord_y_square = calcScalingSquare(PADDING_MAP, tilesize, mapheight, mapwidth)

                            # compute sqm/exits
                            current_map_sqm = mapSqm(mapMatrix)
                            current_map_exits = mapExits(mapMatrix)

                            player_pos, player_count, fire_pos = populateMap(mapMatrix, pop_percent, init_fires)
                            players_movement = []

                            playerSurface, _, _ = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)
                # repair button
                if cursorBoxHit(mouse_x, mouse_y, 602, 738, 477, 507, active_tab_bools[0]) and map_error:
                    repairMap(active_map_path_error)

                    active_map_path = active_map_path_error
                    active_map_path_tmp = active_map_path_error

                    player_scale, current_frame, current_time_float, paused, player_pos, players_movement, player_count, fire_movement, fire_pos, survived, fire_percent, smoke_pos, smoke_movement, smoke_percent, dead = resetState()
                    # clear old map and players
                    mapSurface.fill(COLOR_BACKGROUND)
                    playerSurface.fill(COLOR_KEY)
                    fireSurface.fill((0, 0, 0, 0))
                    smokeSurface.fill((0, 0, 0, 0))
                    # build new map
                    mapSurface, mapMatrix, tilesize, mapwidth, mapheight, map_error = buildMap(active_map_path, mapSurface)
                    if map_error != []:
                        active_map_path = None

                    elif map_error == []: # dont draw players and calculate if error(s)
                        # precalc (better performance) for scaling formula
                        coord_x_circle, coord_y_circle, radius_scale = calcScalingCircle(PADDING_MAP, tilesize, mapheight, mapwidth)
                        coord_x_square, coord_y_square = calcScalingSquare(PADDING_MAP, tilesize, mapheight, mapwidth)

                        # compute sqm/exits
                        current_map_sqm = mapSqm(mapMatrix)
                        current_map_exits = mapExits(mapMatrix)

                        player_pos, player_count, fire_pos = populateMap(mapMatrix, pop_percent, init_fires)
                        players_movement = []

                        playerSurface, _, _ = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)
                # scale plus/minus
                if cursorBoxHit(mouse_x, mouse_y, 918, 932, 364-23, 378-23, active_tab_bools[0]) and active_map_path is not None:
                    if player_scale > 0.5: # crashes if negative radius, keep it > zero
                        player_scale *= 0.8
                        playerSurface, _, _ = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)
                if cursorBoxHit(mouse_x, mouse_y, 965, 979, 364-23, 378-23, active_tab_bools[0]) and active_map_path is not None:
                    if player_scale < 5: # not to big?
                        player_scale *= 1.25
                        playerSurface, _, _ = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x_circle, coord_y_circle, radius_scale, COLOR_PLAYER_GRADIENT)

    # render logic
    if active_tab_bools[0]: # simulation tab
        # no chosen map
        if active_map_path is None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            #displaySurface.fill(COLOR_KEY) # hack
            mapSurface.fill(COLOR_BACKGROUND)
            rmenuSurface.fill(COLOR_BACKGROUND) # important
            rmenuSurface.blit(MENU_FADE, (0, 0)) # important

            # large upload button
            if not map_error:
                if cursorBoxHit(mouse_x, mouse_y, 450, 574, 335, 459, active_tab_bools[0]):
                    mapSurface.blit(BUTTON_UPLOAD_LARGE, (450, 280))
                else:
                    mapSurface.blit(BUTTON_UPLOAD_LARGE0, (450, 280))

            if map_error:
                if cursorBoxHit(mouse_x, mouse_y, 450, 574, 335+250, 459+250, active_tab_bools[0]):
                    mapSurface.blit(BUTTON_UPLOAD_LARGE, (450, 280+250))
                else:
                    mapSurface.blit(BUTTON_UPLOAD_LARGE0, (450, 280+250))
                
                # error page
                mapSurface = showErrorPage(mapSurface, ERROR_BG, FONT_ROBOTOLIGHT_18, FONT_ROBOTOLIGHT_22, active_map_path_error, map_error)

                # repair button hover/blank
                if cursorBoxHit(mouse_x, mouse_y, 602, 738, 477, 507, active_tab_bools[0]):
                    mapSurface.blit(BUTTON_REPAIR_HOVER, (600, 420))
                else:
                    mapSurface.blit(BUTTON_REPAIR_BLANK, (600, 420))

            displaySurface.blit(mapSurface, (0, 55)) # empty here
            displaySurface.blit(rmenuSurface, (909, 45)) # important
        # chosen map
        else: # warning, move most of this out of the render logic to events/semi-static surfaces. mousemotion events etc. must be here though
            rmenuSurface.blit(MENU_RIGHT, (0, 0))

            if current_frame == 0:
                if counter_seconds % 2 == 0: # even
                    rmenuSurface.blit(TIMER_BACKGROUND, (2, 228+1))
                    placeTextAlpha(rmenuSurface, '--', FONT_DIGITAL7MONO_45, COLOR_YELLOW, 71, 249-17+1)
                    placeTextAlpha(rmenuSurface, '--', FONT_DIGITAL7MONO_45, COLOR_YELLOW, 8, 249-17+1)
                else:
                    rmenuSurface.blit(TIMER_BACKGROUND, (2, 228+1))

            # dividers
            #rmenuSurface.blit(DIVIDER_SHORT, (23, y))
            #rmenuSurface.blit(DIVIDER_LONG, (5, y))
            rmenuSurface.blit(DIVIDER_LONG, (5, 33))
            rmenuSurface.blit(DIVIDER_LONG, (5, 102))
            rmenuSurface.blit(DIVIDER_LONG, (5, 322))
            rmenuSurface.blit(DIVIDER_SHORT, (23, 419))
            rmenuSurface.blit(DIVIDER_LONG, (5, 520))
            rmenuSurface.blit(DIVIDER_LONG, (5, 621))

            placeCenterText(rmenuSurface, pathToName(active_map_path_tmp), FONT_ROBOTOREGULAR_20, COLOR_BLACK, 116, 19)

            placeText(rmenuSurface, str(format((round(current_map_sqm)), ',d')).replace(',', ' '), FONT_ROBOTOREGULAR_17, COLOR_BLACK, 31, 37)
            placeText(rmenuSurface, str(round(mapwidth*0.5)) + '×' + str(round(mapheight*0.5)), FONT_ROBOTOREGULAR_17, COLOR_BLACK, 31, 57)
            placeText(rmenuSurface, str(current_map_exits), FONT_ROBOTOREGULAR_17, COLOR_BLACK, 31, 77)

            # inf/people/fire/smoke. Move. Hover/click logic
            rmenuSurface.blit(BUTTON_PEOPLE, (13-2, 111))
            rmenuSurface.blit(BUTTON_FIRE, (44, 111))
            rmenuSurface.blit(BUTTON_SMOKE, (75+2, 111))

            # run button hover/blank
            if current_frame == 0:
                if cursorBoxHit(mouse_x, mouse_y, 900, 1024, 236, 270, active_tab_bools[0]):
                    rmenuSurface.blit(BUTTON_RUN_HOVER, (2, 191))
                else:
                    rmenuSurface.blit(BUTTON_RUN_BLANK, (2, 191))

            elif current_frame > 0:
                rmenuSurface.blit(BUTTON_RUN2_GREEN, (2, 191))
                rmenuSurface.blit(BUTTON_RUN2_FBW, (2+8, 191+8))
                rmenuSurface.blit(BUTTON_RUN2_BW, (2+8+5+14*1, 191+8))
                rmenuSurface.blit(BUTTON_RUN2_PAUSE, (51, 191+8))
                rmenuSurface.blit(BUTTON_RUN2_FW, (113-8-5-14*2, 191+8))
                rmenuSurface.blit(BUTTON_RUN2_FFW, (113-8-14, 191+8))

            # upload button hover/blank
            if cursorBoxHit(mouse_x, mouse_y, 937, 999, 685, 747, active_tab_bools[0]):
                rmenuSurface.blit(BUTTON_UPLOAD_SMALL, (28, 640))
            else:
                rmenuSurface.blit(BUTTON_UPLOAD_SMALL0, (28, 640))


            # timer
            rmenuSurface.blit(BUTTON_TIME_SPEED, (82+3, 311-23))
            placeTextAlpha(rmenuSurface, "1x", FONT_ROBOTOMEDIUM_13, COLOR_BLACK, 88+3, 319-23)

            # player scale
            rmenuSurface.blit(BUTTON_SCALE, (49-21, 313-23))
            rmenuSurface.blit(BUTTON_SCALE_MINUS, (34-3-21, 319-23))
            rmenuSurface.blit(BUTTON_SCALE_PLUS, (75+3-21, 319-23))

            # rmenu statistics
            placeCenterText(rmenuSurface, "Total", FONT_ROBOTOREGULAR_13, COLOR_GREY2, 116, 338)
            placeCenterText(rmenuSurface, str(format(player_count, ',d')).replace(',', ' '), FONT_ROBOTOREGULAR_22, COLOR_BLACK, 116, 359)
            placeCenterText(rmenuSurface, "Left", FONT_ROBOTOREGULAR_13, COLOR_GREY2, 116, 379)
            placeCenterText(rmenuSurface, str(format(player_count - survived - dead, ',d')).replace(',', ' '), FONT_ROBOTOREGULAR_22, COLOR_BLACK, 116, 400)
            placeCenterText(rmenuSurface, "Survivors", FONT_ROBOTOREGULAR_13, COLOR_GREY2, 116, 439-3)
            placeCenterText(rmenuSurface, str(format(survived, ',d')).replace(',', ' '), FONT_ROBOTOREGULAR_22, COLOR_BLACK, 116, 460-3)
            placeCenterText(rmenuSurface, "Dead", FONT_ROBOTOREGULAR_13, COLOR_GREY2, 116, 483-3)
            placeCenterText(rmenuSurface, str(format(dead, ',d')).replace(',', ' '), FONT_ROBOTOREGULAR_22, COLOR_BLACK, 116, 504-3)

            placeCenterText(rmenuSurface, "Fire", FONT_ROBOTOREGULAR_13, COLOR_GREY2, 116, 539-3)
            placeCenterText(rmenuSurface, "{:.0%}".format(fire_percent*0.98), FONT_ROBOTOREGULAR_22, COLOR_BLACK, 116, 560-3) # *0.98 nasty fix
            placeCenterText(rmenuSurface, "Smoke", FONT_ROBOTOREGULAR_13, COLOR_GREY2, 116, 583-3)
            placeCenterText(rmenuSurface, "{:.0%}".format(smoke_percent*0.98), FONT_ROBOTOREGULAR_22, COLOR_BLACK, 116, 604-3) # *0.98 nasty fix

            if current_frame > 0:
                rmenuSurface.blit(TIMER_BACKGROUND, (2, 228+1))
                setClock(rmenuSurface, FONT_DIGITAL7MONO_45, math.floor(current_time_float))

            # important blit order
            displaySurface.blit(rmenuSurface, (909, 45))
            displaySurface.blit(mapSurface, (0, 55))
            if surface_toggle[2]:
                displaySurface.blit(smokeSurface, (0, 55))        
            if surface_toggle[1]:
                displaySurface.blit(fireSurface, (0, 55))
            if surface_toggle[0]:
                displaySurface.blit(playerSurface, (0, 55))
            displaySurface.blit(throbberSurface, (952, 200))

    elif active_tab_bools[1]: # settings tab
        # no chosen map
        if active_map_path == None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            settingsSurface.fill(COLOR_BACKGROUND)
            placeText(settingsSurface, "Choose map first [Settings], id01", FONT_ROBOTOREGULAR_24, COLOR_BLACK, 200, 300)
            minimapSurface.fill(COLOR_BACKGROUND) # wierd1
            displaySurface.blit(minimapSurface, (517, 60)) # wierd2
        # map chosen
        else:
            settingsSurface.fill(COLOR_BACKGROUND)

            settingsSurface.blit(BG_SETTINGS, (6, 1))
            placeCenterText(settingsSurface, pathToName(active_map_path), FONT_ROBOTOREGULAR_26, COLOR_BLACK, 530, 30)

            if player_pos != []:
                placeText(settingsSurface, "Populated sim, but paused, id02", FONT_ROBOTOREGULAR_14, COLOR_BLACK, 100, 300)
            paused = True
            placeText(settingsSurface, "Placeholder settingsSurface, id03", FONT_ROBOTOREGULAR_14, COLOR_BLACK, 100, 200)

            minimapSurface.fill(COLOR_WHITE)
            minimapSurface, _, _, _, _ = buildMiniMap(active_map_path, minimapSurface)

        displaySurface.blit(settingsSurface, (0, 55))
        displaySurface.blit(MENU_FADE, (0, 45))
        displaySurface.blit(minimapSurface, (517, 60))

    elif active_tab_bools[2]: # statistics tab
        # no chosen map
        if active_map_path == None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            statisticsSurface.fill(COLOR_BACKGROUND)
            placeText(statisticsSurface, "Choose map first [Stats], id04", FONT_ROBOTOREGULAR_24, COLOR_BLACK, 100, 300)
            displaySurface.blit(statisticsSurface, (0, 55))
            minimapSurface.fill(COLOR_BACKGROUND) # wierd1
            displaySurface.blit(minimapSurface, (517, 60)) # wierd2
        # map chosen
        else:
            statisticsSurface.fill(COLOR_BACKGROUND)

            statisticsSurface.blit(BG_STATISTICS, (6, 1))
            placeCenterText(statisticsSurface, pathToName(active_map_path), FONT_ROBOTOREGULAR_26, COLOR_BLACK, 530, 30)

            if plot_rendered:
                raw_data = rawPlotRender(rawPlot())
                raw_data2 = rawPlotRender(rawPlot2())
                raw_data3 = rawPlotRender(rawPlot3(json_stat_content[0]))
                plot_rendered = True
                
                # quadrant 1
                #surf = pygame.image.fromstring(raw_data, (plot_x, plot_y), "RGB")
                #statisticsSurface.blit(surf, (10, 5))
                
                # quadrant 2
                surf = pygame.image.fromstring(raw_data3, (150, 120), "RGB")
                statisticsSurface.blit(surf, (345, 60))
                
                surf = pygame.image.fromstring(raw_data3, (150, 120), "RGB")
                statisticsSurface.blit(surf, (345, 200))
                
                # quadrant 3
                surf = pygame.image.fromstring(raw_data2, (plot_x, plot_y), "RGB")
                statisticsSurface.blit(surf, (10, 361))

                # quadrant 4
                surf = pygame.image.fromstring(raw_data, (plot_x, plot_y), "RGB")
                statisticsSurface.blit(surf, (517, 361))
                
                minimapSurface.fill(COLOR_WHITE)
                minimapSurface, _, _, _, _ = buildMiniMap(active_map_path, minimapSurface)

            if player_pos != []:
                placeText(statisticsSurface, "Populated sim, but paused, id05", FONT_ROBOTOREGULAR_14, COLOR_BLACK, 100, 200)
            paused = True
            placeText(statisticsSurface, "Placeholder statisticsSurface, id06", FONT_ROBOTOREGULAR_14, COLOR_BLACK, 100, 270)

        displaySurface.blit(statisticsSurface, (0, 55))
        displaySurface.blit(MENU_FADE, (0, 45))
        displaySurface.blit(minimapSurface, (517, 60))
    else:
        raise NameError('No active tab')

    prev_time, fps = calcFPS(prev_time, target_fps, True)

    #showDebugger(displaySurface, MENU_BACKGROUND, MENU_FADE, FONT_ROBOTOREGULAR_11 mapwidth, mapheight, active_tab_bools, pop_percent, paused, counter_seconds, current_time_float, player_pos, fps, active_map_path, tilesize, mouse_x, mouse_y, pipe_input, player_scale)

    # move to timed event for performance?
    if go_running:
        json_pid = open('../src/pid.txt', 'r').read()
        json_pid_content = json.loads(json_pid)
        if not psutil.pid_exists(json_pid_content):
            go_running = False
            throbberSurface.fill((0, 0, 0, 0))
            
            json_stat = open('peopleStats.txt', 'r').read()
            json_stat_content = json.loads(json_stat)
            plot_rendered = True


    # update displaySurface
    pygame.display.flip() # .update(<surface_args>) instead?
