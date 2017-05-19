#!/usr/bin/python3

import pygame
import sys
import os
import numpy
import math
import time
import tkinter as tk # replace
import subprocess
import copy
import json
#import wx
import signal
import colorsys
import inflect
#import psutil

import _thread

from tkinter import filedialog

import matplotlib
matplotlib.use("Agg")
import matplotlib.backends.backend_agg as agg
import pylab
import matplotlib.pylab as plt

from sys import getsizeof
from subprocess import Popen, PIPE
from utils import *
from pygame.locals import *
from PIL import Image
from pygame import gfxdraw # use later, AA

def restart_program():
    python = sys.executable
    os.execl(python, python, * sys.argv)

# file dialog init
#app = wx.App()
#frame_wx = wx.Frame(None, -1, 'win.py')

inflect = inflect.engine()

# init game
pygame.init()

# set window icon and program name
icon = pygame.image.load(os.path.join('gui', 'window_icon.png'))
pygame.display.set_icon(icon)
pygame.display.set_caption(GAME_NAME)

# game time
TIMER1000 = USEREVENT + 1
TIMER100 = USEREVENT + 2
TIMER10 = USEREVENT + 3
pygame.time.set_timer(TIMER1000, 1000)
pygame.time.set_timer(TIMER100, 100)
pygame.time.set_timer(TIMER10, 10)
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
map_error_length = 0
map_error_visible = 5
text_rect = None
error_text_spacing = 34
error_text_x = 745
error_text_y = 242

player_pos = [] # might use this as indicator to not populate instead of players_movement?
players_movement = []
fire_pos = []
fire_movement = []
smoke_pos = []		
smoke_movement = []

init_fires = 1

fire_percent = 0
survived = 0

surface_toggle = [True, True, True]

opacity = 0
opacity2 = 0

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

# create gradients, can not do it in utils...
#COLOR_FIRE_GRADIENT = linearInterpolationGradient([(255, 255, 0), (107, 14, 14)], 3)

COLOR_PLAYER_GRADIENT = interpolateTuple((46, 127, 0),(66, 181, 0), 100) # 2 steps == len 3
COLOR_FIRE_GRADIENT = interpolateTuple((253, 207, 88), (128, 9, 9), 100) # 2 steps == len 3
COLOR_SMOKE_GRADIENT = interpolateTuple((254, 254, 254), (100, 100, 100), 100) # 2 steps == len 3


# create the display surface, the overall main screen size that will be rendered
displaySurface = pygame.display.set_mode((GAME_RES)) # FULLSCREEN, DOUBLEBUF?
displaySurface.fill(COLOR_BACKGROUND) # no color leaks in the beginning
displaySurface.set_colorkey(COLOR_KEY, pygame.RLEACCEL) # RLEACCEL unstable?

# create surfaces
mapSurface = createSurface(907, 713-PADDING_MAP)
minimapSurface = createSurface(495, 344)
playerSurface = createSurface(907, 713-PADDING_MAP)
fireSurface = createSurface(907, 713-PADDING_MAP)
fireSurface.fill((0,0,0,0))
smokeSurface = createSurface(907, 713-PADDING_MAP)
smokeSurface.fill((0, 0, 0, 0))
rmenuSurface = createSurface(115, 723)
statisticsSurface = createSurface(1024, 713)
settingsSurface = createSurface(1024, 713)

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

TIMER_BACKGROUND = loadImage('gui', 'timer.png')

DIVIDER_LONG = loadImage('gui', 'divider_long.png')
DIVIDER_SHORT = loadImage('gui', 'divider_short.png')

BUTTON_SCALE = loadImage('gui', 'scale.png')
BUTTON_SCALE_PLUS = loadImage('gui', 'scale_plus.png')
BUTTON_SCALE_MINUS = loadImage('gui', 'scale_minus.png')

BUTTON_TIME_SPEED = loadImage('gui', 'time_speed.png')

BUTTON_INF = loadImage('gui', 'inf.png')
BUTTON_PEOPLE = loadImage('gui', 'people.png')
BUTTON_FIRE = loadImage('gui', 'fire.png')
BUTTON_SMOKE = loadImage('gui', 'smoke.png')

ERROR_BG = loadImage('gui', 'error_bg.png')

file_opt = fileDialogInit()

# game loop
while True:
    # event logic
    for event in pygame.event.get():
        if event.type == QUIT:
            pygame.quit()
            sys.exit()
        # time events
        TIMER10_bool = False # remove? used before for fading, could be used for other stuff
        if event.type == TIMER10:
            TIMER10_bool = True
        if event.type == TIMER100: # 100ms per movement (or frame), meaning top speed of ((0.5*(1000/100))/1)*3.6 = 18 km/h
            # FETCH PIPE HERE TO VAR, check if same (should never be, then Go sends data to slow)
            # ADD PIPE LOGIC HERE, FETCH EACH TIME FRAME
                # time check/update, only if new timeframe. Save/update time as "?counter_seconds" (handle float but render int)
                # people movement
                # fire movement
                # smoke movement
                # update values if not same...
                # save everything so that we can render backwards when paused. (see K_f and K_g)
                # render below

            # just demonstrate time movement, but read from pipe changes above. Sort out the if clauses
            #if False: # deactivate
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

                        #for player in range(len(player_pos)):
                        #    player_pos[player] = players_movement[player][current_frame]   # change this to pipe var later.
                                                                                            # handle empty or let Go fill it
                                                                                            # with dummy (same pos) data?
                                                                                            # like movement[player][dir at frame] == "up", go up
                        # pause if last frame, problem with go's pipe?
                        if current_frame == len(players_movement[0]) - 1:
                            paused = True

                        playerSurface, survived = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT)
                        fireSurface = drawFire(fireSurface, fire_pos, tilesize, mapheight, mapwidth, COLOR_FIRE_GRADIENT)
                        smokeSurface = drawSmoke(smokeSurface, smoke_pos, tilesize, mapheight, mapwidth, COLOR_SMOKE_GRADIENT)
                        
                        fire_percent = len(fire_pos) / (current_map_sqm * 4) # shady?
        if event.type == TIMER1000: # just specific for clock animation, 10*100ms below instead?
            counter_seconds += 1
        # keyboard events, later move to to mouse click event
        elif event.type == KEYDOWN:
            if event.key == K_r:
                restart_program()
            elif event.key == K_1: # depopulate
                    pygame.quit()
                    sys.exit()
            if active_tab_bools[0] and active_map_path is not None: # do not add time/pos if no map
                # these two need to read from _saved_ pipe movement, cant go back otherwise. and only possible when paused
                # add 'not' for not populated, time runs anyhow for these
                if event.key == K_g and paused and players_movement != []: # forwards player movement from players_movement, move later to timed game event
                        if current_frame < len(players_movement[0]) - 1: # no (more) movement tuples
                            current_frame += 1
                            current_time_float += 0.1
                            for player in range(len(player_pos)):
                                player_pos[player] = players_movement[player][current_frame]
                            #for tile in range(len(fire_pos)):
                                #fire_pos[tile] = fire_movement[tile][current_frame]
                            #if len(fire_movement) > current_frame:
                            fire_pos = fire_movement[current_frame]
                            smoke_pos = smoke_movement[current_frame]
                            playerSurface, survived = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT)
                            fireSurface = drawFire(fireSurface, fire_pos, tilesize, mapheight, mapwidth, COLOR_FIRE_GRADIENT)
                            smokeSurface = drawSmoke(smokeSurface, smoke_pos, tilesize, mapheight, mapwidth, COLOR_SMOKE_GRADIENT)
                            #displaySurface.blit(fireSurface, (0, 55))

                            fire_percent = len(fire_pos) / (current_map_sqm * 4) # shady?


                elif event.key == K_f and paused and player_pos != []: # backwards player movement from players_movement, move later to timed game event
                        if current_frame > 0: # no (more) movement tuples
                            current_frame -= 1
                            current_time_float -= 0.1
                            for player in range(len(player_pos)):
                                player_pos[player] = players_movement[player][current_frame]

                            fire_pos = fire_movement[current_frame]
                            playerSurface, survived = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT)
                            fireSurface = drawFire(fireSurface, fire_pos, tilesize, mapheight, mapwidth, COLOR_FIRE_GRADIENT)
                            smoke_pos = smoke_movement[current_frame]
                            smokeSurface = drawSmoke(smokeSurface, smoke_pos, tilesize, mapheight, mapwidth, COLOR_SMOKE_GRADIENT)

                            fire_percent = len(fire_pos) / (current_map_sqm * 4) # shady?
                elif event.key == K_2 and paused:
                    current_frame = 0
                    current_time_float = 0.0
                    playerSurface, survived = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT)
                elif event.key == K_m and paused and current_frame == 0:
                    _thread.start_new_thread(goThread, (mapMatrix, player_pos, players_movement, fire_pos, fire_movement, smoke_pos, smoke_movement) )

                elif event.key == K_s and paused and player_pos != []:                                  
                    if players_movement != [] and current_frame < len(players_movement[0]) - 1:  # do not start time frame clock if not pupulated.
                                                                                                 # problems if we have no people?
                                                                                                 # shaky logic with current frame, can otherwise
                                                                                                 # run/unpause at last frame
                        paused = False
                elif event.key == K_p: # for use with cursorHitBox
                    paused = True # if paused == True -> False?
                elif event.key == K_o: # for use with cursorHitBox
                    if pop_percent < 0.9:
                        pop_percent *= 1.25
                elif event.key == K_l: # for use with cursorHitBox
                    if pop_percent > 0.1:
                        pop_percent *= 0.8
                    
                ## testing fire
                elif event.key == K_PLUS:
                    init_fires = init_fires + 1
                elif event.key == K_MINUS:
                    if init_fires > 1:
                        init_fires = init_fires - 1                        
                ## testing fire

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
                    #if players_movement != [] and current_frame == 0: # warning, cannot run sim without people due to this.
                                                                       # shitty handling for no respawn (current_frame)?,
                                                                       # if respawn is needed, remove current_frame
                        player_pos, player_count, fire_pos = populateMap(mapMatrix, pop_percent, init_fires)
                        playerSurface, survived = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT)
                       # fireSurface = drawFire(fireSurface, fire_pos, tilesize, mapheight, mapwidth, COLOR_FIRE_GRADIENT)
                    else:
                        print('Depop first')
                elif event.key == K_z: # depopulate
                    _, current_frame, current_time_float, paused, player_pos, players_movement, player_count, fire_movement, fire_pos, survived, fire_percent, smoke_pos, smoke_movement = resetState()
                    playerSurface, survived = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT)
                    fireSurface = drawFire(fireSurface, fire_pos, tilesize, mapheight, mapwidth, COLOR_FIRE_GRADIENT)
                    smokeSurface = drawFire(smokeSurface, smoke_pos, tilesize, mapheight, mapwidth, COLOR_SMOKE_GRADIENT)
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
                    #openFileDialog = wx.FileDialog(frame_wx, "Open", "", "", "PNG Maps (*.png)|*.png", wx.FD_OPEN | wx.FD_FILE_MUST_EXIST)
                    #openFileDialog.ShowModal()
                    #active_map_path_tmp = openFileDialog.GetPath()
                    #openFileDialog.Destroy()
                    active_map_path_tmp = fileDialogPath()
                    #active_map_path_tmp = 'Map07.png'
                    if active_map_path_tmp != "": #and active_map_path != "/":
                        active_map_path_error = active_map_path_tmp
                        active_map_path = active_map_path_tmp # (2/2)fixed bug for exiting folder window, not sure why tmp is needed
                        # reset state.
                        player_scale, current_frame, current_time_float, paused, player_pos, players_movement, player_count, fire_movement, fire_pos, survived, fire_percent, smoke_pos, smoke_movement = resetState()
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
                            coord_x, coord_y, radius_scale = calcScaling(PADDING_MAP, tilesize, mapheight, mapwidth)

                            # compute sqm/exits
                            current_map_sqm = mapSqm(mapMatrix)
                            current_map_exits = mapExits(mapMatrix)

                            player_pos, player_count, fire_pos = populateMap(mapMatrix, pop_percent, init_fires)
                            players_movement = []

                            playerSurface, _ = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT)
                # upload button map error
                if cursorBoxHit(mouse_x, mouse_y, 450, 574, 335+250, 459+250, active_tab_bools[0]) and map_error:
                    #openFileDialog = wx.FileDialog(frame_wx, "Open", "", "", "PNG Maps (*.png)|*.png", wx.FD_OPEN | wx.FD_FILE_MUST_EXIST)
                    #openFileDialog.ShowModal()
                    #active_map_path_tmp = openFileDialog.GetPath()
                    #openFileDialog.Destroy()
                    active_map_path_tmp = fileDialogPath()
                    if active_map_path_tmp == "":
                        active_map_path_tmp = active_map_path
                    #active_map_path_tmp = 'Map07.png'
                    elif active_map_path_tmp != "": #and active_map_path != "/":
                        active_map_path_error = active_map_path_tmp
                        active_map_path = active_map_path_tmp # (2/2)fixed bug for exiting folder window, not sure why tmp is needed
                        # reset state.
                        player_scale, current_frame, current_time_float, paused, player_pos, players_movement, player_count, fire_movement, fire_pos, survived, fire_percent, smoke_pos, smoke_movement = resetState()
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
                            coord_x, coord_y, radius_scale = calcScaling(PADDING_MAP, tilesize, mapheight, mapwidth)

                            # compute sqm/exits
                            current_map_sqm = mapSqm(mapMatrix)
                            current_map_exits = mapExits(mapMatrix)

                            player_pos, player_count, fire_pos = populateMap(mapMatrix, pop_percent, init_fires)
                            players_movement = []

                            playerSurface, _ = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT)
                # upload button routine rmenu
                if cursorBoxHit(mouse_x, mouse_y, 937, 999, 685, 747, active_tab_bools[0]) and active_map_path is not None:
                    #openFileDialog = wx.FileDialog(frame_wx, "Open", "", "", "PNG Maps (*.png)|*.png", wx.FD_OPEN | wx.FD_FILE_MUST_EXIST)
                    #openFileDialog.ShowModal()
                    #active_map_path_tmp = openFileDialog.GetPath()
                    #openFileDialog.Destroy()
                    active_map_path_tmp = fileDialogPath()
                    #active_map_path_tmp = 'MapError.png'
                    if active_map_path_tmp != "": #and active_map_path != "/":
                        active_map_path_error = active_map_path_tmp
                        active_map_path = active_map_path_tmp # (2/2)fixed bug for exiting folder window, not sure why tmp is needed
                        # reset state.
                        player_scale, current_frame, current_time_float, paused, player_pos, players_movement, player_count, fire_movement, fire_pos, survived, fire_percent, smoke_pos, smoke_movement = resetState()
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
                            coord_x, coord_y, radius_scale = calcScaling(PADDING_MAP, tilesize, mapheight, mapwidth)

                            # compute sqm/exits
                            current_map_sqm = mapSqm(mapMatrix)
                            current_map_exits = mapExits(mapMatrix)

                            player_pos, player_count, fire_pos = populateMap(mapMatrix, pop_percent, init_fires)
                            players_movement = []

                            playerSurface, _ = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT)
                # scale plus/minus
                if cursorBoxHit(mouse_x, mouse_y, 918, 932, 364-23, 378-23, active_tab_bools[0]) and active_map_path is not None:
                    if player_scale > 0.5: # crashes if negative radius, keep it > zero
                        player_scale *= 0.8
                        playerSurface, _ = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT)
                if cursorBoxHit(mouse_x, mouse_y, 965, 979, 364-23, 378-23, active_tab_bools[0]) and active_map_path is not None:
                    if player_scale < 5: # not to big?
                        player_scale *= 1.25
                        playerSurface, _ = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale, COLOR_PLAYER_GRADIENT)

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
                mapSurface.blit(ERROR_BG, (260, 120))
                placeCenterText(mapSurface, pathToName(active_map_path_error) + '.png has ' + inflect.number_to_words(len(map_error)) + ' invalid pixels!', 'Roboto-Light.ttf', 22, COLOR_BLACK, 1024, 150)

                for idx in range(map_error_visible):
                    if idx < len(map_error):
                        if map_error[idx][2][3] != 255 and (map_error[idx][2][0:3] != COLOR_WHITE and map_error[idx][2][0:3] != COLOR_BLACK and map_error[idx][2][0:3] != COLOR_RED_PNG and map_error[idx][2][0:3] != COLOR_KEY):
                            _, _, text_y = placeCenterText(mapSurface, '(' + str(map_error[idx][0]) + ', ' + str(map_error[idx][1]) + ')', 'Roboto-Light.ttf', 22, COLOR_BLACK, error_text_x, error_text_y + idx * error_text_spacing) # coord
                            placeCenterText(mapSurface, str(map_error[idx][2][0:3]), 'Roboto-Light.ttf', 22, COLOR_BLACK, error_text_x + 333-15, error_text_y + idx * error_text_spacing) # rgb
                            placeCenterText(mapSurface, str("{:.0%}".format(map_error[idx][2][3]/255)), 'Roboto-Light.ttf', 22, COLOR_BLACK, error_text_x + 580-15, error_text_y + idx * error_text_spacing) # opacity
                            pygame.gfxdraw.box(mapSurface, pygame.Rect(423, text_y-16, 33, 33), map_error[idx][2][0:3])
                        elif map_error[idx][2][3] != 255:
                            _, _, text_y = placeCenterText(mapSurface, '(' + str(map_error[idx][0]) + ', ' + str(map_error[idx][1]) + ')', 'Roboto-Light.ttf', 22, COLOR_BLACK, error_text_x, error_text_y + idx * error_text_spacing) # coord
                            placeCenterText(mapSurface, str("{:.0%}".format(map_error[idx][2][3]/255)), 'Roboto-Light.ttf', 22, COLOR_BLACK, error_text_x + 580-15, error_text_y + idx * error_text_spacing) # opacity
                            pygame.gfxdraw.box(mapSurface, pygame.Rect(423, text_y-16, 33, 33), map_error[idx][2][0:3])
                        elif map_error[idx][2][0:3] != COLOR_WHITE and map_error[idx][2][0:3] and COLOR_BLACK and map_error[idx][2][0:3] != COLOR_RED_PNG and map_error[idx][2][0:3] != COLOR_KEY:
                            _, _, text_y = placeCenterText(mapSurface, '(' + str(map_error[idx][0]) + ', ' + str(map_error[idx][1]) + ')', 'Roboto-Light.ttf', 22, COLOR_BLACK, error_text_x, error_text_y + idx * error_text_spacing) # coord
                            placeCenterText(mapSurface, str(map_error[idx][2][0:3]), 'Roboto-Light.ttf', 22, COLOR_BLACK, error_text_x + 333-15, error_text_y + idx * error_text_spacing) # rgb
                            pygame.gfxdraw.box(mapSurface, pygame.Rect(423, text_y-16, 33, 33), map_error[idx][2][0:3])
                        map_error_length += 1
                        if idx == map_error_visible - 1 and len(map_error) - map_error_length > 0:
                            placeText(mapSurface, 'and ' + inflect.number_to_words(len(map_error) - map_error_length) + ' more...', 'Roboto-Light.ttf', 18, COLOR_GREY2, 323, error_text_y + (idx + 1) * error_text_spacing - 15)
                            map_error_length = 0
                            break

                placeCenterText(mapSurface, 'Fix and reload the map or try to [Auto-Repair]', 'Roboto-Light.ttf', 18, COLOR_BLACK, 1024, 430)
            displaySurface.blit(mapSurface, (0, 55)) # empty here
            displaySurface.blit(rmenuSurface, (909, 45)) # important
        # chosen map
        else: # warning, move most of this out of the render logic to events/semi-static surfaces
            rmenuSurface.blit(MENU_RIGHT, (0, 0))

            if current_frame == 0:
                if counter_seconds % 2 == 0: # even
                    rmenuSurface.blit(TIMER_BACKGROUND, (2, 228+1))
                    placeTextAlpha(rmenuSurface, '--', 'digital-7-mono.ttf', 45, COLOR_YELLOW, 71, 249-17+1)
                    placeTextAlpha(rmenuSurface, '--', 'digital-7-mono.ttf', 45, COLOR_YELLOW, 8, 249-17+1)
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

            placeCenterText(rmenuSurface, pathToName(active_map_path_tmp), 'Roboto-Regular.ttf', 20, COLOR_BLACK, 116, 19)

            placeText(rmenuSurface, str(format((round(current_map_sqm)), ',d')).replace(',', ' '), 'Roboto-Regular.ttf', 17, COLOR_BLACK, 31, 37)
            placeText(rmenuSurface, str(round(mapwidth*0.5)) + '×' + str(round(mapheight*0.5)), 'Roboto-Regular.ttf', 17, COLOR_BLACK, 31, 57)
            placeText(rmenuSurface, str(current_map_exits), 'Roboto-Regular.ttf', 17, COLOR_BLACK, 31, 77)

            # inf/people/fire/smoke. Move. Hover/click logic
            rmenuSurface.blit(BUTTON_INF, (13, 111))
            rmenuSurface.blit(BUTTON_PEOPLE, (10+3, 140))
            rmenuSurface.blit(BUTTON_FIRE, (32+3+3+2+4, 140))
            rmenuSurface.blit(BUTTON_SMOKE, (54+3+3+4+4+2+2+3, 140))

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
            placeTextAlpha(rmenuSurface, "1x", 'Roboto-Medium.ttf', 13, COLOR_BLACK, 88+3, 319-23)

            # player scale
            rmenuSurface.blit(BUTTON_SCALE, (49-21, 313-23))
            rmenuSurface.blit(BUTTON_SCALE_MINUS, (34-3-21, 319-23))
            rmenuSurface.blit(BUTTON_SCALE_PLUS, (75+3-21, 319-23))

            # rmenu statistics
            placeCenterText(rmenuSurface, "Total", 'Roboto-Regular.ttf', 13, COLOR_GREY2, 116, 338)
            placeCenterText(rmenuSurface, str(format(player_count, ',d')).replace(',', ' '), 'Roboto-Regular.ttf', 22, COLOR_BLACK, 116, 359)
            placeCenterText(rmenuSurface, "Left", 'Roboto-Regular.ttf', 13, COLOR_GREY2, 116, 379)
            placeCenterText(rmenuSurface, str(format(player_count - survived, ',d')).replace(',', ' '), 'Roboto-Regular.ttf', 22, COLOR_BLACK, 116, 400)
            placeCenterText(rmenuSurface, "Survivors", 'Roboto-Regular.ttf', 13, COLOR_GREY2, 116, 439-3)
            placeCenterText(rmenuSurface, str(format(survived, ',d')).replace(',', ' '), 'Roboto-Regular.ttf', 22, COLOR_BLACK, 116, 460-3)
            placeCenterText(rmenuSurface, "Dead", 'Roboto-Regular.ttf', 13, COLOR_GREY2, 116, 483-3)
            placeCenterText(rmenuSurface, str(format(0, ',d')).replace(',', ' '), 'Roboto-Regular.ttf', 22, COLOR_BLACK, 116, 504-3)

            placeCenterText(rmenuSurface, "Fire", 'Roboto-Regular.ttf', 13, COLOR_GREY2, 116, 539-3)
            placeCenterText(rmenuSurface, "{:.0%}".format(fire_percent), 'Roboto-Regular.ttf', 22, COLOR_BLACK, 116, 560-3)
            placeCenterText(rmenuSurface, "Smoke", 'Roboto-Regular.ttf', 13, COLOR_GREY2, 116, 583-3)
            placeCenterText(rmenuSurface, "{:.0%}".format(0.00), 'Roboto-Regular.ttf', 22, COLOR_BLACK, 116, 604-3)

            if current_frame > 0:
                rmenuSurface.blit(TIMER_BACKGROUND, (2, 228+1))
                setClock(rmenuSurface, math.floor(current_time_float))

            # important blit order
            displaySurface.blit(rmenuSurface, (909, 45))
            displaySurface.blit(mapSurface, (0, 55))

#            if surface_toggle[0]:
#                displaySurface.blit(playerSurface, (0, 55))
            if surface_toggle[2]:
                displaySurface.blit(smokeSurface, (0, 55))        
            if surface_toggle[1]:
                displaySurface.blit(fireSurface, (0, 55))
            if surface_toggle[0]:
                displaySurface.blit(playerSurface, (0, 55))


                
    elif active_tab_bools[1]: # settings tab
        # no chosen map
        if active_map_path == None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            settingsSurface.fill(COLOR_BACKGROUND)
            placeText(settingsSurface, "Choose map first [Settings], id01", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 200, 300)
            minimapSurface.fill(COLOR_BACKGROUND) # wierd1
            displaySurface.blit(minimapSurface, (517, 60)) # wierd2
        # map chosen
        else:

            settingsSurface.fill(COLOR_BACKGROUND)

            settingsSurface.blit(BG_SETTINGS, (6, 1))
            placeCenterText(settingsSurface, pathToName(active_map_path), 'Roboto-Regular.ttf', 26, COLOR_BLACK, 530, 30)

            if player_pos != []:
                placeText(settingsSurface, "Populated sim, but paused, id02", 'Roboto-Regular.ttf', 14, COLOR_BLACK, 100, 300)
            paused = True
            placeText(settingsSurface, "Placeholder settingsSurface, id03", 'Roboto-Regular.ttf', 14, COLOR_BLACK, 100, 200)

            minimapSurface.fill(COLOR_WHITE)
            minimapSurface, _, _, _, _ = buildMiniMap(active_map_path, minimapSurface)

        displaySurface.blit(settingsSurface, (0, 55))
        displaySurface.blit(MENU_FADE, (0, 45))
        displaySurface.blit(minimapSurface, (517, 60))

    elif active_tab_bools[2]: # statistics tab
        # no chosen map
        if active_map_path == None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            statisticsSurface.fill(COLOR_BACKGROUND)
            placeText(statisticsSurface, "Choose map first [Stats], id04", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 100, 300)
            displaySurface.blit(statisticsSurface, (0, 55))
            minimapSurface.fill(COLOR_BACKGROUND) # wierd1
            displaySurface.blit(minimapSurface, (517, 60)) # wierd2
        # map chosen
        else:
            statisticsSurface.fill(COLOR_BACKGROUND)

            statisticsSurface.blit(BG_STATISTICS, (6, 1))
            placeCenterText(statisticsSurface, pathToName(active_map_path), 'Roboto-Regular.ttf', 26, COLOR_BLACK, 530, 30)

            if not plot_rendered:
                raw_data = rawPlotRender(rawPlot())
                raw_data2 = rawPlotRender(rawPlot2())
                raw_data3 = rawPlotRender(rawPlot3())
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
                placeText(statisticsSurface, "Populated sim, but paused, id05", 'Roboto-Regular.ttf', 14, COLOR_BLACK, 100, 200)
            paused = True
            placeText(statisticsSurface, "Placeholder statisticsSurface, id06", 'Roboto-Regular.ttf', 14, COLOR_BLACK, 100, 270)

        displaySurface.blit(statisticsSurface, (0, 55))
        displaySurface.blit(MENU_FADE, (0, 45))
        displaySurface.blit(minimapSurface, (517, 60))
    else:
        raise NameError('No active tab')

    # fps calc, remove later
    curr_time = time.time() # so now we have time after processing
    diff = curr_time - prev_time # frame took this much time to process and render
    delay = max(1.0/target_fps - diff, 0) # if we finished early, wait the remaining time to desired fps, else wait 0 ms
    time.sleep(delay)
    fps = 1.0/(delay + diff) # fps is based on total time ("processing" diff time + "wasted" delay time)
    prev_time = curr_time
    #pygame.display.set_caption("{0}: {1:.2f}".format(GAME_NAME, fps))

    # debugger/. remove later, bad fps
    displaySurface.blit(MENU_BACKGROUND, (570, 0)) # bs1
    displaySurface.blit(MENU_FADE, (-120, 45)) # bs^2

    placeText(displaySurface, "DEBUGGER", 'Roboto-Regular.ttf', 11, COLOR_BLACK, 570, 0)
    placeText(displaySurface, "+mapwidth: " + str(mapwidth) + "til" + " (" + str(mapwidth*0.5)+ "m)", 'Roboto-Regular.ttf', 11, COLOR_BLACK, 570, 10)
    placeText(displaySurface, "+mapheight: " + str(mapheight) + "til" + " (" + str(mapheight*0.5)+ "m)", 'Roboto-Regular.ttf', 11, COLOR_BLACK, 570, 20)
    placeText(displaySurface, "+tab: " + str(active_tab_bools), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 570, 30)

    # (debugger) check out of bounds.
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
    placeText(displaySurface, "+pop_%: " + str(round(pop_percent, 2)), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 31)
    placeText(displaySurface, "+paused: " + str(paused), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 0)
    placeText(displaySurface, "+elapsed: " + str(counter_seconds), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 11)
    placeText(displaySurface, "+frame_float: " + str(round(current_time_float, 2)), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 21)

    placeText(displaySurface, "+p_scale: " + str(round(player_scale, 2)), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 0)
    placeText(displaySurface, "+populated: " + str(player_pos != []), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 11)
    placeText(displaySurface, "+fps: " + str(round(fps)), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 21)
    placeText(displaySurface, "+file: " + str(active_map_path), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 31)

    placeText(displaySurface, "+tilesize: " + str(tilesize), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 910, 0)
    placeText(displaySurface, "+mouse xy: " + str(mouse_x) + "," + str(mouse_y), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 910, 11)
    placeText(displaySurface, "+pipe_in: " + str(pipe_input), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 910, 21)
    # /debugger

    # update displaySurface
    pygame.display.flip() # .update(<surface_args>) instead?
