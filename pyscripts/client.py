import socket
from packet import Packet
import time
import random

IP = "127.0.0.1"
PORT = 10001

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

cnt = 0
while True:
    # if cnt>10:
    #     break

    _src = 1
    _dst = 2

    cnt += 1
    pkt = Packet()
    buf = pkt.pkt2Buf(_src, _dst)
    sock.sendto(buf, (IP, PORT))
    print("[{}] sent {} bytes".format(cnt, len(buf)))
    time.sleep(2)
    # time.sleep(2)
    