import os
import sys
import re

def sort_human(l):
  convert = lambda text: float(text) if text.isdigit() else text
  alphanum_key = lambda key: [ convert(c) for c in re.split('([0-9]+)', key) ] 
  l.sort( key=alphanum_key )
  return l

def make(val):
  ret = `val`
  if val < 10000:
    ret = "0" + ret
  if val < 1000:
    ret = "0" + ret
  if val < 100:
    ret = "0" + ret
  if val < 10:
    ret = "0" + ret
  return ret

files = []
for f in os.listdir(sys.argv[1]):
    files.append(f)

files = sort_human(files)
print files
count = int(sys.argv[2])
for f in files:
  if f.endswith('.png'):
    out = make(count)
    os.popen("cp %s%s finished/%s.png" % (sys.argv[1], f,out)).readlines()
    count += 1
