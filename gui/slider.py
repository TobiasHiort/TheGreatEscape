import pygame
from pygame.locals import *
import sys
import os
from pygame import gfxdraw # use later, AA

SCREEN_WIDTH = 1024
SCREEN_HEIGHT = 768

WHITE = (255, 255, 255)
#RED   = (255,   0,   0)
RED = (255, 0, 0)
GREY  = (45,   45,  45)

PADDING = 7

CIRCLE_RADIUS = 8

pygame.init()
screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))

#circle = pygame.draw.circle(screen, GREY, (round(SCREEN_WIDTH/2) + CIRCLE_RADIUS, 0), CIRCLE_RADIUS)
circle = pygame.rect.Rect(100-CIRCLE_RADIUS, 0, 16, 16)

rectangle_dragging = False
clock = pygame.time.Clock()
running = True
while running:
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            running = False
        elif event.type == KEYDOWN:
            if event.key == K_r:
                python = sys.executable
                os.execl(python, python, * sys.argv)
        elif event.type == MOUSEBUTTONDOWN:
            if event.button == 1:
                if circle.collidepoint(event.pos):
                    rectangle_dragging = True
                    mouse_x, _ = event.pos
                    offset_x = circle.x - mouse_x

        elif event.type == MOUSEBUTTONUP:
            if event.button == 1:
                rectangle_dragging = False

        elif event.type == MOUSEMOTION:
            if rectangle_dragging:
                mouse_x, _ = event.pos
                circle.x = mouse_x + offset_x
                print(circle)
                print(round(circle[0]/(SCREEN_WIDTH - PADDING * 2), 2))

    screen.fill(WHITE)

####

    pygame.draw.line(screen, RED, (PADDING, PADDING), (SCREEN_WIDTH-PADDING, PADDING), 4)
    pygame.draw.line(screen, GREY, (circle.center[0]-1, circle.center[1]-1), (SCREEN_WIDTH-PADDING, PADDING), 4)

    # circle.center[0], circle.center[1] = (x, y)
    #pygame.draw.line(screen, GREY, (PADDING, PADDING), (SCREEN_WIDTH-PADDING, PADDING), 4)

    #pygame.draw.line(screen, RED, (PADDING, PADDING+5), (SCREEN_WIDTH-PADDING, PADDING+5), 4)
    #pygame.draw.circle(screen, GREY, (circle.center), 8)

    pygame.gfxdraw.aacircle(screen, circle.center[0], circle.center[1], 8, RED)
    pygame.gfxdraw.filled_circle(screen, circle.center[0], circle.center[1], 8, RED)


####

    pygame.draw.line(screen, RED, (PADDING, PADDING), (SCREEN_WIDTH-PADDING, PADDING), 4)
    pygame.draw.line(screen, GREY, (circle.center[0]-1, circle.center[1]-1), (SCREEN_WIDTH-PADDING, PADDING), 4)

    # circle.center[0], circle.center[1] = (x, y)
    #pygame.draw.line(screen, GREY, (PADDING, PADDING), (SCREEN_WIDTH-PADDING, PADDING), 4)

    #pygame.draw.line(screen, RED, (PADDING, PADDING+5), (SCREEN_WIDTH-PADDING, PADDING+5), 4)
    #pygame.draw.circle(screen, GREY, (circle.center), 8)

    pygame.gfxdraw.aacircle(screen, circle.center[0], circle.center[1], 8, RED)
    pygame.gfxdraw.filled_circle(screen, circle.center[0], circle.center[1], 8, RED)



    pygame.display.flip()
    clock.tick(30)
pygame.quit()
