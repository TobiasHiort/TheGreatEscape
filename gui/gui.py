#!/usr/bin/python3

import pygame
import sys
import os
import numpy
import math
import time
import tkinter as tk # replace
import subprocess

from utils import *
from pygame.locals import *
from tkinter import filedialog # remove?
from PIL import Image
#from pygame import gfxdraw # use later, AA

# horrible function, does not render correct without it, must be in main
def blitMenuButtons(displaySurface, button):
    if button == "simulation":
        # simulation button
        displaySurface.blit(BUTTON_SETTINGS_BLANK, (202, 0))
        displaySurface.blit(BUTTON_STATISTICS_BLANK, (382, 0))
        displaySurface.blit(BUTTON_SIMULATION_ACTIVE, (0, 0))
        active_tab_bools = [True, False, False]

    elif button == "settings":
        # settings button
        displaySurface.blit(BUTTON_SIMULATION_BLANK, (0, 0))
        displaySurface.blit(BUTTON_STATISTICS_BLANK, (382, 0))
        displaySurface.blit(BUTTON_SETTINGS_ACTIVE, (202, 0))
        active_tab_bools = [False, True, False]
    
    elif button == "statistics":
        # statistics button
        displaySurface.blit(BUTTON_SIMULATION_BLANK, (0, 0))
        displaySurface.blit(BUTTON_SETTINGS_BLANK, (202, 0))
        displaySurface.blit(BUTTON_STATISTICS_ACTIVE, (382, 0))
        active_tab_bools = [False, False, True]
    else:
        raise ValueError('No button named "' + button + '"')
    return displaySurface, active_tab_bools

# init game
pygame.init()

# set window icon and program name
icon = pygame.image.load(os.path.join('gui', 'window_icon.png'))
pygame.display.set_icon(icon)
pygame.display.set_caption(GAME_NAME)

# clock from init
clock = pygame.time.Clock()

# game event timer, 1 and 0.1 seconds, and vars
TIMER1000 = USEREVENT + 1
TIMER100 = USEREVENT + 2
pygame.time.set_timer(TIMER1000, 1000)
pygame.time.set_timer(TIMER100, 100) # if 0.1s is the smallest time unit, only add 0.1 seconds then and floor to 1
counter_seconds = 0 # counter for TIMER1000
current_frame = 0 # which time frame for movement, int
current_time_float = 0.0 # float time for accurate time frame measurement, right now 0.1s per time frame.
paused = True
player_scale = 1.0

# debugger inits, not needed later
active_map_path_tmp = None
tilesize = None
mapwidth = 0
mapheight = 0
pipe_input = None

player_pos = [] # might use this as indicator to not populate instead of players_movement?

# movement, warning: temp
#players_movement = []
                   # include start pos
players_movement = [[(0, 1), (0, 2), (0, 3), (0, 4), (0, 5), (0, 6), (0, 7),  (0, 8),  (0, 9),  (0, 10), (0, 11), (0, 12)],
                    [(1, 2), (1, 2), (1, 3), (1, 3), (1, 4), (1, 4), (1, 5),  (1, 5),  (1, 6),  (1,  6), (1,  7), (1,  7)],
                    [(2, 3), (2, 4), (2, 5), (2, 6), (2, 7), (2, 8), (2, 9),  (2, 10), (2, 11), (2, 12), (2, 13), (2, 14)],
                    [(3, 4), (3, 5), (3, 6), (3, 7), (3, 8), (3, 9), (3, 10), (3, 11), (3, 12), (3, 13), (3, 14), (3, 15)]] 

# create the display surface, the overall main screen size that will be rendered
displaySurface = pygame.display.set_mode((GAME_RES)) # ,pygame.NOFRAME
displaySurface.fill(COLOR_BACKGROUND) # no color leaks in the beginning
displaySurface.set_colorkey(COLOR_KEY, pygame.RLEACCEL)

# create surfaces
mapSurface = createSurface(907, 713-PADDING_MAP)
playerSurface = createSurface(907, 713-PADDING_MAP)
rmenuSurface = createSurface(115, 723)
statisticsSurface = createSurface(1024, 713)
settingsSurface = createSurface(1024, 713)

# load and blit menu components to main surface (possibly remove blit, blit then in the game loop?). Warning for blits.
MENU_FADE = pygame.image.load(os.path.join('gui', 'menu_fade.png')).convert()
displaySurface.blit(MENU_FADE, (0, 45))

MENU_BACKGROUND = pygame.image.load(os.path.join('gui', 'menu_background.png')).convert()
displaySurface.blit(MENU_BACKGROUND, (0, 0))

MENU_RIGHT = pygame.image.load(os.path.join('gui', 'menu_right.png')).convert()
rmenuSurface.blit(MENU_RIGHT, (0, 0)) # remove? blit game_loop

# load and blit buttons in init state
BUTTON_SIMULATION_ACTIVE, BUTTON_SIMULATION_BLANK, BUTTON_SIMULATION_HOVER = buildButton("simulation")
BUTTON_SETTINGS_ACTIVE, BUTTON_SETTINGS_BLANK, BUTTON_SETTINGS_HOVER = buildButton("settings")
BUTTON_STATISTICS_ACTIVE, BUTTON_STATISTICS_BLANK, BUTTON_STATISTICS_HOVER = buildButton("statistics")

# warning for blits, expand later with hover etc., do buildButton but for rmenu
BUTTON_RUN = pygame.image.load(os.path.join('gui', 'run.png')).convert()
rmenuSurface.blit(BUTTON_RUN, (2, 80)) # remove  later, blit in gameloop.
BUTTON_UPLOAD = pygame.image.load(os.path.join('gui', 'upload_small.png')).convert_alpha()
rmenuSurface.blit(BUTTON_UPLOAD, (28, 640)) # remove  later, blit in gameloop.
TIMER_BACKGROUND = pygame.image.load(os.path.join('gui', 'timer.png')).convert()
rmenuSurface.blit(TIMER_BACKGROUND, (2, 160)) # remove  later, blit in gameloop.

# get file dialog options
file_opt = fileDialogInit()

# game loop
while True:
    # event logic
    for event in pygame.event.get():
        if event.type == QUIT:
            pygame.quit()
            sys.exit()
        # user event (time)
        elif event.type == TIMER1000: # just specific for clock animation, 10*100ms below instead?
            # clock animation when no map loaded
            if active_map_path == None or active_map_path == "": # "" = cancel on choosing map
                if counter_seconds % 2 == 0: # even
                    placeText(rmenuSurface, '--', 'digital-7-mono.ttf', 45, COLOR_YELLOW, 71, 164)
                    placeText(rmenuSurface, '--', 'digital-7-mono.ttf', 45, COLOR_YELLOW, 8, 164)
                else:
                    rmenuSurface.blit(TIMER_BACKGROUND, (2, 160))
            counter_seconds += 1

        elif event.type == TIMER100: # 100ms per movement (or frame), meaning top speed of ((0.5*(1000/100))/1)*3.6 = 18 km/h
            # FETCH PIPE HERE TO VAR, check if same (should never be, then Go sends data to slow)
            # ADD PIPE LOGIC HERE, FETCH EACH TIME FRAME
                # time check/update, only if new timeframe. Save/update time as "?counter_seconds" (handle float but render int)
                # people movement
                # fire movement
                # smoke movement
                # update values if not same...
                # save everything so that we can render backwards when paused. (see K_f and K_g)
                # render below

            # just demonstrate time movement, but read from pipe changes above
            #if False: # deactivate
            if players_movement != []: # warning, handling unpopulated map
                if active_map_path != None or active_map_path != "": # "" = cancel on choosing map
                    if current_frame < len(players_movement[0])-1 and not paused: # no more movement tuples and not paused
                        current_frame += 1
                        current_time_float += 0.1
                        for player in range(len(player_pos)):
                            player_pos[player] = players_movement[player][current_frame] # change this to pipe var later.
                                                                                        # handle empty or let go fill it
                                                                                        # with dummy (same pos) data?
                                                                                        # like movement[player][dir at frame] == "up", go up

        # keyboard events (import as function?)
        elif event.type == KEYDOWN:
            if active_tab_bools[0] and active_map_path != None: # do not add time/pos if no map
                if event.key == K_o: # increase player scale, remove later
                    if player_scale < 3: # not to big?
                        player_scale *= 1.25
                elif event.key == K_l: # decrease player scale, remove later
                    if player_scale > 0.5: # crashes if negative radius, keep it > zero
                        player_scale *= 0.8

                # these two need to read from _saved_ pipe movement, cant go back otherwise. and only possible when paused
                # add 'not' for not populated, time runs anyhow for these
                elif event.key == K_g and paused and player_pos != []: # forwards player movement from player_movement, move later to timed game event
                        if current_frame < len(players_movement[0])-1: # no (more) movement tuples
                            current_frame += 1
                            current_time_float += 0.1
                            for player in range(len(player_pos)):
                                player_pos[player] = players_movement[player][current_frame]
                elif event.key == K_f and paused and player_pos != []: # backwards player movement from player_movement, move later to timed game event
                        if current_frame > 0: # no (more) movement tuples
                            current_frame -= 1
                            current_time_float -= 0.1
                            for player in range(len(player_pos)):
                                player_pos[player] = players_movement[player][current_frame]
                
                elif event.key == K_m and paused:
                    # read stdout through pipe TEST
                    #popen = subprocess.call('./hello') # just a call
                    popen = subprocess.Popen("./hello", stdout=subprocess.PIPE)
                    popen.wait()
                    pipe_input = popen.stdout.read()
                    print(str(pipe_input))

                elif event.key == K_s and paused and player_pos != []:
                    if players_movement != []:  # do not start time frame clock if not pupulated. problems if we have no people?
                        paused = False
                elif event.key == K_p: # for use with cursorHitBox
                    paused = True # if paused == True -> False?
                elif event.key == K_a: # populate, warning. use after randomizing init pos
                    paused = True
                    # function of this? maybe scrap for direction movement instead
                    if players_movement != [] and current_frame == 0: # warning, cannot run sim without people due to this.
                                                                      # shitty handling for no respawn (current_frame)?,
                                                                      # if respawn is needed, remove current_time_frame
                        player_pos = []
                        for player in range(len(players_movement)): # change players_movement to randomized init positions (this is before Go)
                            player_pos.append(players_movement[player][0]) # start pos for each player
                    else:
                        print("Handle no players, or already run")
        # mouse motion events (hovers), (import as function?)
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
                    displaySurface, active_tab_bools = blitMenuButtons(displaySurface, "simulation")
                # settings button
                elif cursorBoxHit(mouse_x, mouse_y, 202, 382, 0, 45, not(active_tab_bools[1])):
                    displaySurface, active_tab_bools = blitMenuButtons(displaySurface, "settings")
                # statistics button
                elif cursorBoxHit(mouse_x, mouse_y, 383, 575, 0, 45, not(active_tab_bools[2])):
                    displaySurface, active_tab_bools = blitMenuButtons(displaySurface, "statistics")
                
                # upload button routine
                elif cursorBoxHit(mouse_x, mouse_y, 937, 999, 685, 747, active_tab_bools[0]):
                    active_map_path_tmp = fileDialogPath()
                    if active_map_path_tmp != "":
                        active_map_path = active_map_path_tmp # (2/2)fixed bug for exiting folder window, not sure why tmp is needed
                        # reset state.
                        player_scale, current_frame, current_time_float, paused, player_pos = resetState()
                        # clear old map
                        mapSurface.fill(COLOR_BACKGROUND) 
                        # build new map
                        mapSurface, mapMatrix, tilesize, mapwidth, mapheight = buildMap(active_map_path, mapSurface)
                        # compute/update square meters
                        current_map_sqm = mapSqm(mapMatrix)
                        current_map_exits = mapExits(mapMatrix)

    # render logic
    if active_tab_bools[0]:
        # simulation tab
        if active_map_path == None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            mapSurface.fill(COLOR_BACKGROUND)

            # important blit order
            displaySurface.blit(mapSurface, (0, 55)) # empty here
            displaySurface.blit(rmenuSurface, (909, 45))
        else:
            # right menu (class?)
            rmenuSurface.blit(MENU_RIGHT, (0, 0))
            
            placeText(rmenuSurface, active_map_path, 'Roboto-Regular.ttf', 20, COLOR_BLACK, 12, 5)
            placeText(rmenuSurface, "sqm: " + str(round(current_map_sqm)), 'Roboto-Regular.ttf', 20, COLOR_BLACK, 12, 35)
            placeText(rmenuSurface, "exits: " + str(current_map_exits), 'Roboto-Regular.ttf', 20, COLOR_BLACK, 12, 55)

            # move these to event handler
            rmenuSurface.blit(BUTTON_RUN, (2, 80))
            rmenuSurface.blit(BUTTON_UPLOAD, (28, 640))

            rmenuSurface.blit(TIMER_BACKGROUND, (2, 160))
            
            setClock(rmenuSurface, math.floor(current_time_float))

            # draw players
            playerSurface = drawPlayer(playerSurface, player_pos, tilesize, mapheight, mapwidth, player_scale)
            
            # important blit order
            displaySurface.blit(mapSurface, (0, 55))
            displaySurface.blit(playerSurface, (0, 55))
            displaySurface.blit(rmenuSurface, (909, 45))
    elif active_tab_bools[1]:
        # settings tab
        displaySurface.blit(settingsSurface, (0, 55))
        displaySurface.blit(MENU_FADE, (0, 45))

        # (no) chosen map handler
        if active_map_path == None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            settingsSurface.fill(COLOR_BACKGROUND)
            placeText(settingsSurface, "Choose map first [Settings]", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 200, 300)
        else:
            settingsSurface.fill(COLOR_BACKGROUND)
            placeText(settingsSurface, "Placeholder settingsSurface", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 200, 300)
    elif active_tab_bools[2]:
        # statistics tab
        displaySurface.blit(statisticsSurface, (0, 55))
        displaySurface.blit(MENU_FADE, (0, 45))

        # (no) chosen map handler
        if active_map_path == None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            statisticsSurface.fill(COLOR_BACKGROUND)
            placeText(statisticsSurface, "Choose map first [Stats]", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 200, 300)
        else:
            statisticsSurface.fill(COLOR_BACKGROUND)
            placeText(statisticsSurface, "Placeholder statisticsSurface", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 200, 300)
    else:
        raise NameError('No active tab')
    
    # debugger. remove later, bad fps
    displaySurface.blit(MENU_BACKGROUND, (570, 0)) # bs1
    displaySurface.blit(MENU_FADE, (-120, 45)) # bs^2
    
    placeText(displaySurface, "DEBUGGER", 'Roboto-Regular.ttf', 11, COLOR_BLACK, 570, 0)
    placeText(displaySurface, "+mapwidth: " + str(mapwidth) + "til" + " (" + str(mapwidth*0.5)+ "m)", 'Roboto-Regular.ttf', 11, COLOR_BLACK, 570, 10)
    placeText(displaySurface, "+mapheight: " + str(mapheight) + "til" + " (" + str(mapheight*0.5)+ "m)", 'Roboto-Regular.ttf', 11, COLOR_BLACK, 570, 20)
    placeText(displaySurface, "+tab: " + str(active_tab_bools), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 570, 30)

    placeText(displaySurface, "+p_pos: " + str(player_pos), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 48)
    placeText(displaySurface, "+paused: " + str(paused), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 1)
    placeText(displaySurface, "+elapsed: " + str(counter_seconds), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 14)
    placeText(displaySurface, "+frame_float: " + str(round(current_time_float, 2)), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 26)

    placeText(displaySurface, "+p_scale: " + str(round(player_scale, 2)), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 1)
    placeText(displaySurface, "+populated: " + str(player_pos != []), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 14)
    placeText(displaySurface, "+file: " + str(active_map_path), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 26)

    placeText(displaySurface, "+tilesize: " + str(tilesize), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 910, 1)
    placeText(displaySurface, "+mouse xy: " + str(mouse_x) + "," + str(mouse_y), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 910, 14)
    placeText(displaySurface, "+pipe_in: " + str(pipe_input), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 910, 26)

    # update display if not quitting
    pygame.display.flip() # .update(<surface_args>) instead?

    # fps, remove later
    clock.tick(50)
    pygame.display.set_caption("fps: " + str(clock.get_fps()))
