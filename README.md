## How to run

  make ping

This will run the augmented tool and append the ttest output in the ttest_output.txt file.
If you want to clear the ttest_output.txt file, run the following command
  make clean
And then run the make ping command again.
## Setup
+---------------------+              +------------------------+                +---------+
|         LAN         |              | Raspberry Pi Router    |                |  Client |
|---------------------|              |------------------------|                |---------|
|                     |----(eth0)----|                        |----(br-lan)----|         |
|                     |              |------------------------|                |         |
|                     |              |        (br-lan)        |                |         |
|                     |              +------------------------+                |         |
|                     |                                                        |         |
+---------------------+                                                        +---------+


## How to change the upload rate 
  tc qdisc add dev eth0 root netem rate 30mbit

## How to change the download rate
To shape the download traffic, we add a virtual interface (ifb0) to capture and redirect incoming traffic from the physical network interface . Effectively, we capture incoming traffic on the virtual interface, shape it according to requirement and redirect it into the physical interface.

  modprobe ifb numifbs=1
  ip link add name ifb0 type ifb
  ip link set dev ifb0 up
  tc qdisc add dev eth0 handle ffff: ingress
  tc filter add dev eth0 parent ffff: protocol ip u32 match u32 0 0 action mirred egress redirect dev ifb0
  tc qdisc replace dev ifb0 root handle 1: htb default 1
  tc class add dev ifb0 parent 1: classid 1:1 htb rate 100mbit


+---------------------+                    +------------------------+                +---------+
|         LAN         |                    | Raspberry Pi Router    |                |  Client |
|---------------------|                    |------------------------|                |---------|
|                     |---(ifb0)-(eth0)----|                        |----(br-lan)----|         |
|                     |                    |------------------------|                |         |
|                     |                    |        (br-lan)        |                |         |
|                     |                    +------------------------+                |         |
|                     |                                                              |         |
+---------------------+                                                              +---------+


## References
  https://dl.acm.org/doi/pdf/10.1145/3230543.3230549
  https://www.mdpi.com/1424-8220/23/2/923
  https://www.researchgate.net/publication/354194462_Measuring_and_Localising_Congestion_in_Mobile_Broadband_Networks

