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
import wx

from sys import getsizeof
from subprocess import Popen, PIPE
from utils import *
from pygame.locals import *
from PIL import Image
#from pygame import gfxdraw # use later, AA

# file dialog init
app = wx.App()
frame_wx = wx.Frame(None, -1, 'win.py')

# REMOVE
print("splitPipeData: " + str(splitPipeData(5, "abcdefg12345678")))

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
target_fps = 60
prev_time = time.time() # for fps

# variables
counter_seconds = 0 # counter for TIMER1000
counter_10ms = 0
current_frame = 0 # which time frame for movement, int
current_time_float = 0.0 # float time for accurate time frame measurement, right now 0.1s per time frame.
paused = True
player_scale = 1.0
player_count = 0
pop_percent = 0.1 # init as this later?

player_pos = [] # might use this as indicator to not populate instead of players_movement?
players_movement = []

opacity = 0
opacity2 = 0
opacity3 = 0

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

# create the display surface, the overall main screen size that will be rendered
displaySurface = pygame.display.set_mode((GAME_RES)) # FULLSCREEN, DOUBLEBUF?
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
MENU_FADE = loadImage('gui', 'menu_fade.png')
displaySurface.blit(MENU_FADE, (0, 45)) # blit in game_loop?

MENU_BACKGROUND = loadImage('gui', 'menu_background.png')

MENU_RIGHT = loadImage('gui', 'menu_right.png')

# load buttons in init state
BUTTON_SIMULATION_ACTIVE = loadImage('gui', 'simulation_active.png')
BUTTON_SIMULATION_BLANK = loadImage('gui', 'simulation_blank.png')
BUTTON_SIMULATION_HOVER = loadImage('gui', 'simulation_hover.png')

BUTTON_SETTINGS_ACTIVE = loadImage('gui', 'settings_active.png')
BUTTON_SETTINGS_BLANK = loadImage('gui', 'settings_blank.png')
BUTTON_SETTINGS_HOVER = loadImage('gui', 'settings_hover.png')

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
        elif event.type == TIMER10:
            counter_10ms += 1
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

                        playerSurface = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale)
        # keyboard events, later move to to mouse click event
        elif event.type == KEYDOWN:
            if active_tab_bools[0] and active_map_path is not None: # do not add time/pos if no map
                # these two need to read from _saved_ pipe movement, cant go back otherwise. and only possible when paused
                # add 'not' for not populated, time runs anyhow for these
                if event.key == K_g and paused and players_movement != []: # forwards player movement from players_movement, move later to timed game event
                        if current_frame < len(players_movement[0])-1: # no (more) movement tuples
                            current_frame += 1
                            current_time_float += 0.1
                            for player in range(len(player_pos)):
                                player_pos[player] = players_movement[player][current_frame]
                            playerSurface = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale)

                elif event.key == K_f and paused and player_pos != []: # backwards player movement from players_movement, move later to timed game event
                        if current_frame > 0: # no (more) movement tuples
                            current_frame -= 1
                            current_time_float -= 0.1
                            for player in range(len(player_pos)):
                                player_pos[player] = players_movement[player][current_frame]
                            playerSurface = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale)

                elif event.key == K_m and paused:
                    # read stdout through pipe TEST
                    #popen = subprocess.call('./hello') # just a call
                    child = Popen('./gotest', stdin=subprocess.PIPE, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True)
                    child.stdout.flush()
                    child.stdin.flush()
                    #print(getsizeof(json.dumps(mapMatrix.tolist())))

                    map_jsons = json.dumps(mapMatrix.tolist())
                    test54 = splitPipeData(byte_limit, map_jsons)
                    print(test54[0])

                    #print(, file=child.stdin)

                    #fromgo_json = child.stdout.readline().rstrip('\n')

                    #print(getsizeof(fromgo_json))

                    #data1 = json.loads(fromgo_json)
                    #print(getsizeof(data1))

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
                        playerSurface = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale)
                    else:
                        print('Depop first')
                elif event.key == K_z: # depopulate
                    _, current_frame, current_time_float, paused, player_pos, player_count = resetState()
                    playerSurface = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale)
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
                    openFileDialog = wx.FileDialog(frame_wx, "Open", "", "", "PNG Maps (*.png)|*.png", wx.FD_OPEN | wx.FD_FILE_MUST_EXIST)
                    openFileDialog.ShowModal()
                    active_map_path_tmp = openFileDialog.GetPath()
                    openFileDialog.Destroy()
                    if active_map_path_tmp != "": #and active_map_path != "/":
                        active_map_path = active_map_path_tmp # (2/2)fixed bug for exiting folder window, not sure why tmp is needed
                        # reset state.
                        player_scale, current_frame, current_time_float, paused, player_pos, player_count = resetState()
                        # clear old map and players
                        mapSurface.fill(COLOR_BACKGROUND)
                        playerSurface.fill(COLOR_KEY)
                        # build new map
                        mapSurface, mapMatrix, tilesize, mapwidth, mapheight = buildMap(active_map_path, mapSurface)
                        mapSurface.set_alpha(0)
                        opacity3 = 0

                        # precalc (better performance) for scaling formula
                        coord_x, coord_y, radius_scale = calcScaling(PADDING_MAP, tilesize, mapheight, mapwidth)

                        # compute sqm/exits
                        current_map_sqm = mapSqm(mapMatrix)
                        current_map_exits = mapExits(mapMatrix)

                        player_pos, player_count = populateMap(mapMatrix, pop_percent)
                        players_movement = []
                        
                        playerSurface = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale)
                # upload button routine rmenu
                if cursorBoxHit(mouse_x, mouse_y, 937, 999, 685, 747, active_tab_bools[0]) and active_map_path is not None:
                    openFileDialog = wx.FileDialog(frame_wx, "Open", "", "", "PNG Maps (*.png)|*.png", wx.FD_OPEN | wx.FD_FILE_MUST_EXIST)
                    openFileDialog.ShowModal()
                    active_map_path_tmp = openFileDialog.GetPath()
                    openFileDialog.Destroy()
                    if active_map_path_tmp != "": #and active_map_path != "/":
                        active_map_path = active_map_path_tmp # (2/2)fixed bug for exiting folder window, not sure why tmp is needed
                        # reset state.
                        player_scale, current_frame, current_time_float, paused, player_pos, player_count = resetState()
                        # clear old map and players
                        mapSurface.fill(COLOR_BACKGROUND)
                        playerSurface.fill(COLOR_KEY)
                        # build new map
                        mapSurface, mapMatrix, tilesize, mapwidth, mapheight = buildMap(active_map_path, mapSurface)
                        mapSurface.set_alpha(0)
                        opacity3 = 0

                        # precalc (better performance) for scaling formula
                        coord_x, coord_y, radius_scale = calcScaling(PADDING_MAP, tilesize, mapheight, mapwidth)

                        # compute sqm/exits
                        current_map_sqm = mapSqm(mapMatrix)
                        current_map_exits = mapExits(mapMatrix)

                        player_pos, player_count = populateMap(mapMatrix, pop_percent)
                        players_movement = []

                        playerSurface = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale)
                # scale plus/minus
                if cursorBoxHit(mouse_x, mouse_y, 932, 946, 364, 378, active_tab_bools[0]) and active_map_path is not None:
                    if player_scale > 0.5: # crashes if negative radius, keep it > zero
                        player_scale *= 0.8
                        playerSurface = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale)
                if cursorBoxHit(mouse_x, mouse_y, 993, 1008, 364, 378, active_tab_bools[0]) and active_map_path is not None:
                    if player_scale < 5: # not to big?
                        player_scale *= 1.25
                        playerSurface = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale)


    # render logic
    if active_tab_bools[0]: # simulation tab
        # no chosen map
        if active_map_path is None or active_map_path == "": # if no active map (init), "" = cancel on choosing map
            mapSurface.fill(COLOR_BACKGROUND)

            # large upload button
            if cursorBoxHit(mouse_x, mouse_y, 450, 574, 335, 459, active_tab_bools[0]):
                mapSurface.blit(BUTTON_UPLOAD_LARGE, (450, 280))
            else:
                if counter_10ms % 2 == 0 and opacity < 255:
                    opacity += 5
                    BUTTON_UPLOAD_LARGE0.set_alpha(opacity)
                mapSurface.blit(BUTTON_UPLOAD_LARGE0, (450, 280))

            displaySurface.blit(mapSurface, (0, 55)) # empty here
        # chosen map
        else:



            if current_frame == 0:
                if counter_seconds % 2 == 0: # even
                    rmenuSurface.blit(TIMER_BACKGROUND, (2, 245))
                    placeText(rmenuSurface, '--', 'digital-7-mono.ttf', 45, COLOR_YELLOW, 71, 249)
                    placeText(rmenuSurface, '--', 'digital-7-mono.ttf', 45, COLOR_YELLOW, 8, 249)
                else:
                    rmenuSurface.blit(TIMER_BACKGROUND, (2, 245))

            # dividers
            #rmenuSurface.blit(DIVIDER_SHORT, (23, y))
            #rmenuSurface.blit(DIVIDER_LONG, (5, y))
            rmenuSurface.blit(DIVIDER_LONG, (5, 33))
            rmenuSurface.blit(DIVIDER_SHORT, (23, 78))
            rmenuSurface.blit(DIVIDER_LONG, (5, 350))
            rmenuSurface.blit(DIVIDER_SHORT, (23, 483))
            rmenuSurface.blit(DIVIDER_LONG, (5, 620))

            placeCenterText(rmenuSurface, active_map_path[-9:-4], 'Roboto-Regular.ttf', 20, COLOR_BLACK, 116, 19)

            placeText(rmenuSurface, str(round(current_map_sqm)), 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 35)
            placeText(rmenuSurface, str(current_map_exits), 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 55)

            placeText(rmenuSurface, "str1", 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 85)
            placeText(rmenuSurface, "str2", 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 105)
            placeText(rmenuSurface, "str3", 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 125)
            placeText(rmenuSurface, "str4", 'Roboto-Regular.ttf', 18, COLOR_BLACK, 29, 145)

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
            
            rmenuSurface.blit(BUTTON_SCALE, (49, 313))
            rmenuSurface.blit(BUTTON_SCALE_MINUS, (34-10, 319))
            rmenuSurface.blit(BUTTON_SCALE_PLUS, (75+10, 319))

            # rmenu statistics
            placeCenterText(rmenuSurface, "Total", 'Roboto-Regular.ttf', 18, COLOR_GREY2, 116, 375)
            placeCenterText(rmenuSurface, str(player_count), 'Roboto-Regular.ttf', 28, COLOR_BLACK, 116, 400)
            placeCenterText(rmenuSurface, "Left", 'Roboto-Regular.ttf', 18, COLOR_GREY2, 116, 425)
            placeCenterText(rmenuSurface, "int01", 'Roboto-Regular.ttf', 28, COLOR_BLACK, 116, 450)
            placeCenterText(rmenuSurface, "Survivors", 'Roboto-Regular.ttf', 18, COLOR_GREY2, 116, 500+10)
            placeCenterText(rmenuSurface, "int02", 'Roboto-Regular.ttf', 28, COLOR_BLACK, 116, 525+10)
            placeCenterText(rmenuSurface, "Dead", 'Roboto-Regular.ttf', 18, COLOR_GREY2, 116, 550+10)
            placeCenterText(rmenuSurface, "int03", 'Roboto-Regular.ttf', 28, COLOR_BLACK, 116, 575+10)

            if current_frame > 0:
                rmenuSurface.blit(TIMER_BACKGROUND, (2, 245))
                setClock(rmenuSurface, math.floor(current_time_float))

            # draw players. Removed because it's not necessary to drawplayers each frame! Same for fire and other things.
            #playerSurface = drawPlayer(playerSurface, player_pos, tilesize, player_scale, coord_x, coord_y, radius_scale) # add health here? from player_pos

            # draw fire
            fire_pos = [[3,3,1],[3,4,2],[3,5,3]]
            #fireSurface = drawFire(fireSurface, fire_pos, tilesize, mapheight, mapwidth)

            # important blit order
            # all right menu above. warning, move most of this out of the render logic
            if counter_10ms % 2 == 0 and opacity2 < 255:
                opacity2 += 5
                rmenuSurface.set_alpha(opacity2)
            displaySurface.blit(rmenuSurface, (909, 45))
            rmenuSurface.blit(MENU_RIGHT, (0, 0))

            if counter_10ms % 2 == 0 and opacity3 < 255:
                opacity3 += 5
                mapSurface.set_alpha(opacity3)
            displaySurface.blit(mapSurface, (0, 55))

            #displaySurface.blit(mapSurface, (0, 55))
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
    #placeText(displaySurface, "+oob_id: " + str(p_oob_id), 'Roboto-Regular.ttf', 11, COLOR_BLACK, 710, 43)
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


