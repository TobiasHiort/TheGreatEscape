from __future__ import print_function
import pygame
import os
import subprocess
import sys
from subprocess import Popen, PIPE
import json
import pickle

import time

from pygame.locals import *
#from signal import signal, SIGPIPE, SIG_DFL
#signal(SIGPIPE, SIG_DFL)

# initialize game engine
#pygame.init()
# set screen width/height and caption
#size = [640, 480]
#screen = pygame.display.set_mode(size)
#pygame.display.set_caption('My Game')
# initialize clock. used later in the loop.
#clock = pygame.time.Clock()

# Loop until the user clicks close button
#done = False

# Create child with a pipe to STDIN and STDOUT
child = Popen('./gotest', stdin=subprocess.PIPE, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True)

# Serialize and send JSON as a single line (the default is no indentation).
list1 = [1, 2]
#print(type(list1))

#list1_bytes = json.dumps(list1).encode('utf-8')
#print(type(list1_bytes))

#child.stdin.write(b'hej')

#child.stdin.write("hej\n")

#print(type(json.dumps([1, 2]).encode('utf-8')))
#bytes1 = json.dumps([1, 2]).encode('utf-8')
child.stdin.flush()
print(json.dumps([1, 2]), file=child.stdin)
#child.stdin.write(b'\n')

#json.dump([1, 2], child.stdin)
#child.stdin.write('hej')

#print(list1_bytes, file=child.stdin)


#ttest = json.dump(list1, child.stdin)
#print(type(ttest))
#child.stdin.write('\n')

#child.stdin.flush()
#child.stdin.write("confirm_go_start")
#child.communicate("confirm_go_start".encode("utf-8"))
#child.communicate(str.encode("confirm_go_start"))
#print('From go:', child.stdout.readline().rstrip('\n'))

child.stdout.flush()
fromgo_json = child.stdout.readline().rstrip('\n')
print(fromgo_json)
print(type(fromgo_json))

data1 = json.loads(fromgo_json)

print(data1[0])
#print(type(fromgo_json))

#child.stdin.write("confirm_go_start")

#print("1", file=child.stdin)
#print('From go2:', child.stdout.readline().rstrip('\n'))

#for command in commandlist:
#    print('From PIPE: Q:', child.stdout.readline().rstrip('\n'))
#    print(command, file=child.stdin)
#    #### child.stdin.flush()
#    if command != 'Exit':
#        print('From PIPE: A:', child.stdout.readline().rstrip('\n'))


child.stdin.close()
child.stdout.close()
