import os
import subprocess
import re
import socket


hops = 30
def getip(ht):
    try:
        ip = socket.gethostbyname(ht)
        return ip
    except socket.gaierror:
        return None

destination = input("Enter hostname: ")
if(destination[0]=='w'):
        ipo=getip(destination)
else:
        ipo=destination

timetolive= 1
while timetolive != hops + 1:
    ttl = timetolive
    nping_command = f"nping -c 1 --udp -p 33434 --ttl {ttl} {destination}"
    try:
        nping_output = subprocess.check_output(nping_command, shell=True, text=True)
        _ip = re.findall(r'\d+\.\d+\.\d+\.\d+', nping_output)
        _rtt = re.findall(r'Avg rtt: [-+]?\d*\.\d+', nping_output)
        if _ip and _rtt:
            hop_ip = _ip[2]
            rtt = float(_rtt[0][9:])
            print(f"Hop {ttl}: ip = {hop_ip}, rtt = {rtt:.2f} ms")
            if hop_ip == ipo:
                break
        else:
            print(f"Hop {ttl}: *")
    except subprocess.CalledProcessError:
        print(f"Hop {ttl}: Error executing nping")
        break
    timetolive=timetolive+1