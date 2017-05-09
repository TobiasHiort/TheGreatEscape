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


from sys import getsizeof
from subprocess import Popen, PIPE
from utils import *
from pygame.locals import *
from tkinter import filedialog # remove?
from PIL import Image
#from pygame import gfxdraw # use later, AA

# init game
pygame.init()

#print("splitPipeData: " + str(splitPipeData("abcdefg12345678")))

# set window icon and program name
icon = pygame.image.load(os.path.join('gui', 'window_icon.png'))
pygame.display.set_icon(icon)
pygame.display.set_caption(GAME_NAME)

# game time
TIMER1000 = USEREVENT + 1
TIMER100 = USEREVENT + 2
pygame.time.set_timer(TIMER1000, 1000)
pygame.time.set_timer(TIMER100, 100)
target_fps = 60
prev_time = time.time() # for fps

# variables
counter_seconds = 0 # counter for TIMER1000
current_frame = 0 # which time frame for movement, int
current_time_float = 0.0 # float time for accurate time frame measurement, right now 0.1s per time frame.
paused = True
player_scale = 1.0
player_count = 0
pop_percent = 0.1 # init as this later?

player_pos = [] # might use this as indicator to not populate instead of players_movement?
players_movement = []


#How much data is sent in each pipe
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

# create the display surface, the overall main screen size that will be rendered
displaySurface = pygame.display.set_mode((GAME_RES)) # ,DOUBLEBUF?
displaySurface.fill(COLOR_BACKGROUND) # no color leaks in the beginning
displaySurface.set_colorkey(COLOR_KEY, pygame.RLEACCEL) # RLEACCEL unstable?

# create surfaces
mapSurface = createSurface(907, 713-PADDING_MAP)
playerSurface = createSurface(907, 713-PADDING_MAP)
fireSurface = createSurface(907, 713-PADDING_MAP)
rmenuSurface = createSurface(115, 723)
statisticsSurface = createSurface(1024, 713)
settingsSurface = createSurface(1024, 713)

# load and blit menu components to main surface (possibly remove blit, blit then in the game loop?)
MENU_FADE = pygame.image.load(os.path.join('gui', 'menu_fade.png')).convert()
displaySurface.blit(MENU_FADE, (0, 45)) # remove? blit game_loop

MENU_BACKGROUND = pygame.image.load(os.path.join('gui', 'menu_background.png')).convert()
displaySurface.blit(MENU_BACKGROUND, (0, 0)) # remove? blit game_loop

MENU_RIGHT = pygame.image.load(os.path.join('gui', 'menu_right.png')).convert()
rmenuSurface.blit(MENU_RIGHT, (0, 0)) # remove? blit game_loop

# load buttons in init state
BUTTON_SIMULATION_ACTIVE, BUTTON_SIMULATION_BLANK, BUTTON_SIMULATION_HOVER = buildButton("simulation", True)
BUTTON_SETTINGS_ACTIVE, BUTTON_SETTINGS_BLANK, BUTTON_SETTINGS_HOVER = buildButton("settings", True)
BUTTON_STATISTICS_ACTIVE, BUTTON_STATISTICS_BLANK, BUTTON_STATISTICS_HOVER = buildButton("statistics", True)

## warning for blits, expand later with hover etc., do buildButton but for rmenu. Sort this out later.
BUTTON_RUN_BLANK, BUTTON_RUN_HOVER = buildButton("run", False)

BUTTON_UPLOAD_SMALL = pygame.image.load(os.path.join('gui', 'upload_small.png')).convert_alpha()
#rmenuSurface.blit(BUTTON_UPLOAD_SMALL, (28, 640)) # remove  later, blit in gameloop.
BUTTON_UPLOAD_SMALL0 = pygame.image.load(os.path.join('gui', 'upload_small0.png')).convert_alpha()

BUTTON_UPLOAD_LARGE = pygame.image.load(os.path.join('gui', 'upload_large.png')).convert_alpha()
#rmenuSurface.blit(BUTTON_UPLOAD_LARGE, (58, 640)) # remove  later, blit in gameloop.
BUTTON_UPLOAD_LARGE0 = pygame.image.load(os.path.join('gui', 'upload_large0.png')).convert_alpha()

TIMER_BACKGROUND = pygame.image.load(os.path.join('gui', 'timer.png')).convert()
#rmenuSurface.blit(TIMER_BACKGROUND, (2, 276)) # remove  later, blit in gameloop.

DIVIDER_LONG = pygame.image.load(os.path.join('gui', 'divider_long.png')).convert()
DIVIDER_SHORT = pygame.image.load(os.path.join('gui', 'divider_short.png')).convert()

BUTTON_SCALE = pygame.image.load(os.path.join('gui', 'scale.png')).convert_alpha()
BUTTON_SCALE_PLUS = pygame.image.load(os.path.join('gui', 'scale_plus.png')).convert()
BUTTON_SCALE_MINUS = pygame.image.load(os.path.join('gui', 'scale_minus.png')).convert()

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

            # just demonstrate time movement, but read from pipe changes above. Sort out the if clauses
            #if False: # deactivate
            if players_movement != []: # warning, handling unpopulated map
                if active_map_path is not None or active_map_path != "": # "" = cancel on choosing map
                    if current_frame < len(players_movement[0]) - 1 and not paused: # no more movement coords and not paused
                        current_frame += 1
                        current_time_float += 0.1

                        for player in range(len(player_pos)):
                            player_pos[player] = players_movement[player][current_frame]
                        #for player in range(len(player_pos)):
                        #    player_pos[player] = players_movement[player][current_frame]   # change this to pipe var later.
                                                                                            # handle empty or let Go fill it
                                                                                            # with dummy (same pos) data?
                                                                                            # like movement[player][dir at frame] == "up", go up
                        # pause if last frame, problem with go's pipe?
                        if current_frame == len(players_movement[0]) - 1:
                            paused = True

        # keyboard events, later move to to mouse click event
        elif event.type == KEYDOWN:
            if active_tab_bools[0] and active_map_path is not None: # do not add time/pos if no map
                # these two need to read from _saved_ pipe movement, cant go back otherwise. and only possible when paused
                # add 'not' for not populated, time runs anyhow for these
                if event.key == K_g and paused and player_pos != []: # forwards player movement from players_movement, move later to timed game event
                        if current_frame < len(players_movement[0])-1: # no (more) movement tuples
                            current_frame += 1
                            current_time_float += 0.1
                            for player in range(len(player_pos)):
                                player_pos[player] = players_movement[player][current_frame]

                elif event.key == K_f and paused and player_pos != []: # backwards player movement from players_movement, move later to timed game event
                        if current_frame > 0: # no (more) movement tuples
                            current_frame -= 1
                            current_time_float -= 0.1
                            for player in range(len(player_pos)):
                                player_pos[player] = players_movement[player][current_frame]

                elif event.key == K_m and paused:
                    # read stdout through pipe TEST
                    #popen = subprocess.call('./hello') # just a call
                    child = Popen('./gotest', stdin=subprocess.PIPE, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True)
                    child.stdout.flush()
                    child.stdin.flush()
                    print(getsizeof(json.dumps(mapMatrix.tolist())))


                    #converts stuff to into int
                    map_matrixInt = copy.deepcopy(mapMatrix).astype(int)

                    #pipes length of pipe data
                    print(makeItr(byte_limit, json.dumps(mapMatrix.tolist())), file=child.stdin)
                    print("itrprint: ")
                    print(makeItr(byte_limit, json.dumps(mapMatrix.tolist())))

                    map_list = splitPipeData(byte_limit, json.dumps(mapMatrix.tolist()))
                    #wait for inpput from go
                    #
                    #the real deal pipe mother
                    for map_part in range(map_list)
                        #get input from go

                        #send stuff to go
                        print
                    #print(splitPipeData(byte_limit, json.dumps(mapMatrix.tolist())), file=child.stdin)
                    fromgo_json = child.stdout.readline().rstrip('\n')

                    print(getsizeof(fromgo_json))

                    data1 = json.loads(fromgo_json)
                    print(type(data1))

                    #child.stdin.close()
                    #child.stdout.close()


                elif event.key == K_s and paused and player_pos != []:
                    #print(len(players_movement[0][0]))
                    #print(current_frame)
                    if players_movement != [] and current_frame < len(players_movement[0]) - 1:  # do not start time frame clock if not pupulated.
                                                                                                 # problems if we have no people?
                                                                                                 # shaky logic with current frame, can otherwise
                                                                                                 # run/unpause at last frame
                        paused = False

                elif event.key == K_p: # for use with cursorHitBox
                    paused = True # if paused == True -> False?

                elif event.key == K_a: # populate, warning. use after randomizing init pos
                    paused = True
                    player_scale = 1.0
                    # function of this? maybe scrap for direction movement instead
                    player_count = len(player_pos)
                    if current_frame == 0:
                    #if players_movement != [] and current_frame == 0: # warning, cannot run sim without people due to this.
                                                                       # shitty handling for no respawn (current_frame)?,
                                                                       # if respawn is needed, remove current_frame
                        player_pos, player_count = populateMap(mapMatrix, pop_percent)

                        # remove, for testing. creates a 1 frame movement (players_movement from player_pos).
                        # MUST BE DONE BEFORE TIMER100 EVENT/K_s, not current players_movement otherwise
                        #print(player_pos)
                        player_pos_test1 = copy.deepcopy(player_pos)
                        player_pos_test2 = copy.deepcopy(player_pos)
                        player_pos_test3 = [["foo" for i in range(1)] for j in range(player_count)]
                        for x in range(player_count):
                            #player_pos_test3 = [[],[]]
                            player_pos_test3[x] = [player_pos_test1[x], player_pos_test2[x]]
                        #print(player_pos_test3[0][0][0])
                        for player in range( player_count ):
                            for frame in range(1):
                                 player_pos_test3[player][1][1] += 1
                                 #player_pos_test3[player][frame][0] += 1
                        #print(player_pos_test3)
                        #print(player_pos_test2)
                        #print(player_pos_test2[0])
                        players_movement = copy.deepcopy(player_pos_test3)

                    else:
                        print('Depop first')
                elif event.key == K_z: # depopulate
                    _, current_frame, current_time_float, paused, player_pos, player_count = resetState()
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
                if cursorBoxHit(mouse_x, mouse_y, 450, 574, 335, 459, active_tab_bools[0]) and active_map_path is None:
                    #active_map_path_tmp = fileDialogPath()
                    active_map_path_tmp = "map2.png"
                    if active_map_path_tmp != "": #and active_map_path != "/":
                        active_map_path = active_map_path_tmp # (2/2)fixed bug for exiting folder window, not sure why tmp is needed
                        # reset state.
                        player_scale, current_frame, current_time_float, paused, player_pos, player_count = resetState()
                        # clear old map
                        mapSurface.fill(COLOR_BACKGROUND)
                        # build new map
                        mapSurface, mapMatrix, tilesize, mapwidth, mapheight = buildMap(active_map_path, mapSurface)
                        # compute sqm/exits
                        current_map_sqm = mapSqm(mapMatrix)
                        current_map_exits = mapExits(mapMatrix)

                        player_pos, player_count = populateMap(mapMatrix, pop_percent)
                        players_movement = []
                # upload button routine rmenu
                if cursorBoxHit(mouse_x, mouse_y, 937, 999, 685, 747, active_tab_bools[0]) and active_map_path is not None:
                    active_map_path_tmp = fileDialogPath()
                    if active_map_path_tmp != "": #and active_map_path != "/":
                        active_map_path = active_map_path_tmp # (2/2)fixed bug for exiting folder window, not sure why tmp is needed
                        # reset state.
                        player_scale, current_frame, current_time_float, paused, player_pos, player_count = resetState()
                        # clear old map
                        mapSurface.fill(COLOR_BACKGROUND)
                        # build new map
                        mapSurface, mapMatrix, tilesize, mapwidth, mapheight = buildMap(active_map_path, mapSurface)
                        # compute sqm/exits
                        current_map_sqm = mapSqm(mapMatrix)
                        current_map_exits = mapExits(mapMatrix)

                        player_pos, player_count = populateMap(mapMatrix, pop_percent)
                        players_movement = []
                # scale plus/minus
                if cursorBoxHit(mouse_x, mouse_y, 932, 946, 364, 378, active_tab_bools[0]) and active_map_path is not None:
                    if player_scale > 0.5: # crashes if negative radius, keep it > zero
                        player_scale *= 0.8
                if cursorBoxHit(mouse_x, mouse_y, 993, 1008, 364, 378, active_tab_bools[0]) and active_map_path is not None:
                    if player_scale < 5: # not to big?
                        player_scale *= 1.25


    # render logic
    if active_tab_bools[0]: # simulation tab
        # no chosen map
        if active_map_path is None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            mapSurface.fill(COLOR_BACKGROUND)

            # large upload button
            if cursorBoxHit(mouse_x, mouse_y, 450, 574, 335, 459, active_tab_bools[0]):
                mapSurface.blit(BUTTON_UPLOAD_LARGE, (450, 280))
            else:
                mapSurface.blit(BUTTON_UPLOAD_LARGE0, (450, 280))

            # important blit order
            displaySurface.blit(mapSurface, (0, 55)) # empty here
        # chosen map
        else:
            if current_frame == 0:
                if counter_seconds % 2 == 0: # even
                    placeText(rmenuSurface, '--', 'digital-7-mono.ttf', 45, COLOR_YELLOW, 71, 249)
                    placeText(rmenuSurface, '--', 'digital-7-mono.ttf', 45, COLOR_YELLOW, 8, 249)
                else:
                    rmenuSurface.blit(TIMER_BACKGROUND, (2, 245))

            # right menu
            displaySurface.blit(rmenuSurface, (909, 45))
            rmenuSurface.blit(MENU_RIGHT, (0, 0))

            # dividers
            #rmenuSurface.blit(DIVIDER_SHORT, (23, y))
            #rmenuSurface.blit(DIVIDER_LONG, (5, y))
            rmenuSurface.blit(DIVIDER_LONG, (5, 33))
            rmenuSurface.blit(DIVIDER_SHORT, (23, 78))
            rmenuSurface.blit(DIVIDER_LONG, (5, 350))
            rmenuSurface.blit(DIVIDER_SHORT, (23, 478))
            rmenuSurface.blit(DIVIDER_LONG, (5, 620))

            placeCenterText(rmenuSurface, active_map_path[:-4], 'Roboto-Regular.ttf', 20, COLOR_BLACK, 116, 19)

            placeText(rmenuSurface, str(round(current_map_sqm)), 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 35)
            placeText(rmenuSurface, str(current_map_exits), 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 55)

            placeText(rmenuSurface, "str1", 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 85)
            placeText(rmenuSurface, "str2", 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 105)
            placeText(rmenuSurface, "str3", 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 125)
            placeText(rmenuSurface, "str4", 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 145)

            # run button hover/blank
            if cursorBoxHit(mouse_x, mouse_y, 900, 1024, 236, 270, active_tab_bools[0]):
                rmenuSurface.blit(BUTTON_RUN_HOVER, (2, 191))
            else:
                rmenuSurface.blit(BUTTON_RUN_BLANK, (2, 191))
            # upload button hover/blank
            if cursorBoxHit(mouse_x, mouse_y, 937, 999, 685, 747, active_tab_bools[0]):
                rmenuSurface.blit(BUTTON_UPLOAD_SMALL, (28, 640))
            else:
                rmenuSurface.blit(BUTTON_UPLOAD_SMALL0, (28, 640))

            rmenuSurface.blit(BUTTON_SCALE, (49, 313))
            rmenuSurface.blit(BUTTON_SCALE_MINUS, (34-10, 319))
            rmenuSurface.blit(BUTTON_SCALE_PLUS, (75+10, 319))

            # rmenu statistics
            placeCenterText(rmenuSurface, "Total", 'Roboto-Regular.ttf', 18, COLOR_GREY2, 116, 375)
            placeCenterText(rmenuSurface, str(player_count), 'Roboto-Regular.ttf', 28, COLOR_BLACK, 116, 400)
            placeCenterText(rmenuSurface, "Left", 'Roboto-Regular.ttf', 18, COLOR_GREY2, 116, 425)
            placeCenterText(rmenuSurface, "int01", 'Roboto-Regular.ttf', 28, COLOR_BLACK, 116, 450)
            placeCenterText(rmenuSurface, "Survivors", 'Roboto-Regular.ttf', 18, COLOR_GREY2, 116, 500)
            placeCenterText(rmenuSurface, "int02", 'Roboto-Regular.ttf', 28, COLOR_BLACK, 116, 525)
            placeCenterText(rmenuSurface, "Dead", 'Roboto-Regular.ttf', 18, COLOR_GREY2, 116, 550)
            placeCenterText(rmenuSurface, "int03", 'Roboto-Regular.ttf', 28, COLOR_BLACK, 116, 575)

            rmenuSurface.blit(TIMER_BACKGROUND, (2, 245))
            if current_frame > 0:
                setClock(rmenuSurface, math.floor(current_time_float))

            # draw players
            playerSurface = drawPlayer(playerSurface, player_pos, tilesize, mapheight, mapwidth, player_scale) # add health here? from player_pos

            # draw fire
            fire_pos = [[3,3,1],[3,4,2],[3,5,3]]
            #fireSurface = drawFire(fireSurface, fire_pos, tilesize, mapheight, mapwidth)

            # important blit order
            displaySurface.blit(mapSurface, (0, 55))
            displaySurface.blit(playerSurface, (0, 55))
            #displaySurface.blit(fireSurface, (0, 55))

    elif active_tab_bools[1]: # settings tab
        # no chosen map
        if active_map_path == None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            settingsSurface.fill(COLOR_BACKGROUND)
            placeText(settingsSurface, "Choose map first [Settings], id01", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 200, 300)
        # map chosen
        else:
            settingsSurface.fill(COLOR_BACKGROUND)
            if player_pos != []:
                placeText(settingsSurface, "Populated sim, but paused, id02", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 200, 400)
            paused = True
            placeText(settingsSurface, "Placeholder settingsSurface, id03", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 200, 300)

        displaySurface.blit(settingsSurface, (0, 55))
        displaySurface.blit(MENU_FADE, (0, 45))

    elif active_tab_bools[2]: # statistics tab
        # no chosen map
        if active_map_path == None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            statisticsSurface.fill(COLOR_BACKGROUND)
            placeText(statisticsSurface, "Choose map first [Stats], id04", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 200, 300)
        # map chosen
        else:
            statisticsSurface.fill(COLOR_BACKGROUND)
            if player_pos != []:
                placeText(statisticsSurface, "Populated sim, but paused, id05", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 200, 400)
            paused = True
            placeText(statisticsSurface, "Placeholder statisticsSurface, id06", 'Roboto-Regular.ttf', 24, COLOR_BLACK, 200, 300)

        displaySurface.blit(statisticsSurface, (0, 55))
        displaySurface.blit(MENU_FADE, (0, 45))
    else:
        raise NameError('No active tab')

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
    p_oob = None
    p_oob_id = []
    if player_pos != []:
        for player in range(len(players_movement)):
            if mapMatrix[player_pos[player][1]][player_pos[player][0]] == 1 or mapMatrix[player_pos[player][1]][player_pos[player][0]] == 3:
                p_oob_id.append(player)
    if p_oob_id == []:
        p_oob = False
    else:
        p_oob = True

    #placeText(displaySurface, "+p_pos: " + str(player_pos), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 31)
    placeText(displaySurface, "+p_oob: " + str(p_oob), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 31)
    placeText(displaySurface, "+oob_id: " + str(p_oob_id), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 43)
    placeText(displaySurface, "+paused: " + str(paused), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 0)
    placeText(displaySurface, "+elapsed: " + str(counter_seconds), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 11)
    placeText(displaySurface, "+frame_float: " + str(round(current_time_float, 2)), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 21)

    placeText(displaySurface, "+p_scale: " + str(round(player_scale, 2)), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 0)
    placeText(displaySurface, "+populated: " + str(player_pos != []), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 11)
    placeText(displaySurface, "+file: " + str(active_map_path), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 810, 22)

    placeText(displaySurface, "+tilesize: " + str(tilesize), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 910, 0)
    placeText(displaySurface, "+mouse xy: " + str(mouse_x) + "," + str(mouse_y), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 910, 11)
    placeText(displaySurface, "+pipe_in: " + str(pipe_input), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 910, 21)
    # /debugger

    # update displaySurface
    pygame.display.flip() # .update(<surface_args>) instead?

    # fps calc, remove later
    curr_time = time.time() # so now we have time after processing
    diff = curr_time - prev_time # frame took this much time to process and render
    delay = max(1.0/target_fps - diff, 0) # if we finished early, wait the remaining time to desired fps, else wait 0 ms
    time.sleep(delay)
    fps = 1.0/(delay + diff) # fps is based on total time ("processing" diff time + "wasted" delay time)
    prev_time = curr_time
    pygame.display.set_caption("{0}: {1:.2f}".format(GAME_NAME, fps))
