import os
import time


def call_and_measure(line, resolver):
    start = time.time()
    print("dig @" + resolver + " " + line)
    os.system("dig @" + resolver + " " + line)
    end = time.time()
    return end - start


dohResolving = []
normalResolving = []
defaultResolving = []
with open('domains.txt') as fp:
    for line in fp:
        call_and_measure(line, "127.0.0.1")  # Make sure to add to cache
        time1 = call_and_measure(line, "127.0.0.1")
        dohResolving.append(time1)

        call_and_measure(line, "1.1.1.1")  # Make sure to add to cache
        time2 = call_and_measure(line, "1.1.1.1")
        normalResolving.append(time2)

        ip = os.popen(
            "ipconfig getpacket $(networksetup -listallhardwareports "
            "| awk '/Hardware Port: Wi-Fi/{getline; print $2}') "
            "| grep domain_name_server "
            "| sed -n 's/.*{\\([0-9]\\{1,3\\}\\.[0-9]\\{1,3\\}\\.[0-9]\\{1,3\\}\\.[0-9]\\{1,3\\}\\).*/\\1/p'").read()
        call_and_measure(line, ip)  # Make sure to add to cache
        time3 = call_and_measure(line, ip)
        defaultResolving.append(time3)

print(sum(dohResolving) / len(dohResolving))
print(sum(normalResolving) / len(normalResolving))
print(sum(defaultResolving) / len(defaultResolving))
