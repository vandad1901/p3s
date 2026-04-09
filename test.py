#!/bin/python3

import math
import os
import random
import re
import sys


#
# Complete the 'getOpenApplications' function below.
#
# The function is expected to return a STRING_ARRAY.
# The function accepts STRING_ARRAY commands as parameter.
#


def getOpenApplications(commands):
    n = len(commands)
    firstOpenAfterClear = 0
    for i in range(n):
        if commands[n - i - 1] == "clear":
            firstOpenAfterClear = n - i
            break

    openPrograms = []
    currentIndex = 0
    for i in range(firstOpenAfterClear, n):
        verb, argument = commands[i].split(" ")
        if verb == "open":
            if len(openPrograms) <= currentIndex:
                openPrograms.append(argument)
            else:
                openPrograms[currentIndex] = argument

            currentIndex += 1
        elif verb == "close":
            currentIndex = min(currentIndex - int(argument), 0)

    return openPrograms[:currentIndex]


if __name__ == "__main__":
    fptr = open(os.environ["OUTPUT_PATH"], "w")

    commands_count = int(input().strip())

    commands = []

    for _ in range(commands_count):
        commands_item = input()
        commands.append(commands_item)

    result = getOpenApplications(commands)

    fptr.write("\n".join(result))
    fptr.write("\n")

    fptr.close()
