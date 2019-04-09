import os, time

def call_and_measure(line, resolver):
    start = time.time()
    os.system("dig @" + resolver + " " + line)
    end = time.time()
    return end - start


dohResolving = []
normalResolving = []
with open('domains.txt') as fp:
    for line in fp:
        time1 = call_and_measure(line, "127.0.0.1")
        dohResolving.append(time1)
        time2 = call_and_measure(line, "1.1.1.1")
        normalResolving.append(time2)

print(sum(dohResolving) / len(dohResolving))
print(sum(normalResolving) / len(normalResolving))
