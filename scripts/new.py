file = 10000
path = "/tmp/watch_test/"
size = 100

msg = ""
for i in range(size):
    msg += "a"

import shutil

for i in range(file):
    shutil.move(path + "{}.log".format(i), path + "{}.log.1".format(i))

for i in range(file):
    with open(path + "{}.log".format(i), "a") as f:
        f.write(msg)
