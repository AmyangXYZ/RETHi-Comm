import socket
from packet import Packet
import time
import random

IP = "127.0.0.1"
PORT = 20002

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.bind((IP, PORT))

cnt = 0
while True:
    cnt += 1
    msg, addr = sock.recvfrom(1024)
    print("[%d] received %d bytes"% (cnt, len(msg)))
    # time.sleep(0.001)
    # time.sleep(2)
    