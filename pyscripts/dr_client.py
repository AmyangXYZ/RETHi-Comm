import socket
from packet import Packet
import time
import sqlite3
from sqlite3 import Error

IP = "192.168.1.4"
PORT = 10011
DB_FILE = "dr_lite.db"
SQL_CREATE_POWER_FDD = """ CREATE TABLE IF NOT EXISTS P_SOLAR_DUST_HS (
                                timestamp integer PRIMARY KEY,
                                hs_1 double, hs_2 double,
                                hs_3 double, hs_4, double);"""
SQL_INSERT_POWER_FDD = """ INSERT INTO P_SOLAR_DUST_HS(timestamp,hs_1,hs_2,hs_3,hs_4)
                                VALUES(?,?,?,?,?);"""

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.bind((IP, PORT))

start = time.time()
conn = None
try:
    conn = sqlite3.connect(DB_FILE)
    print("sqlite version",sqlite3.version)
    try:
        c = conn.cursor()
        c.execute(SQL_CREATE_POWER_FDD)
    except Error as e:
        print("Error creating table: ", e)
except Error as e:
    print("Error creating connection: ", e)

"""finally:
    if conn:
        conn.close()
"""
cnt = 0
while True:
    
    data, addr = sock.recvfrom(1024)
    cnt+=1
    pkt = Packet()
   
    if len(data) > 7:
        pkt.buf2Pkt(data)
        _src = pkt.header.src
        _dst = pkt.header.dst
        _type = pkt.header.type
        _priority = pkt.header.priority
        _row = pkt.header.row
        _col = pkt.header.col
        _payload = pkt.payload
        ##handler.handle(_src, _payload)
        # timestamp = int(time.time() * 1e6)
        values = (_payload[0], _payload[1], _payload[2], _payload[3], _payload[4])
        try:
            c.execute(SQL_INSERT_POWER_FDD, values)
            conn.commit()
            print("[{}]insert".format(cnt),values,"into P_SOLAR_DUST_HS")
        except Error as e:
            print("Error inserting data: ",e)
  
    
