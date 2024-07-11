import random

N = 10
m = 5
for _ in range(N):

    val = ""
    for _ in range(m):
        val += str(hex(random.randint(0, 9999999999999999999999)))
    print(val)

