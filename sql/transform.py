#!/usr/bin/python

import sys, re

for line in sys.stdin:
    print(re.sub("\$(.*)\$", "static.whooplist.com/assets/\\1", line.strip()))
