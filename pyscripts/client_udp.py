# Real sub-systems simulating the Simulink
import socket
from packet import Packet
import time
import random

class Socket(object):
    ''' Socket for specific sub-systems
    '''
    def __init__(self, ip, port):
        
        self.ip = ip
        self.port = port
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    def send(self, buf):
        self.sock.sendto(buf, (self.ip, self.port))


class Subsystem(object):
    ''' A sub-system modeled in MCVT
    '''
    def __init__(self, src, dst, type, priority, row, col, length):
        ''' Definitions as per Application Layer Protocol
        '''
        self._src = src
        self._dst = dst
        self._type = type
        self._priority = priority
        self._row = row
        self._col = col
        self._length = length
    
    def to_buffer(self,payload):
        ''' Convert payload with meta-data to Buffer
        '''
        self._payload = payload
        pkt = Packet()
        buf = pkt.pkt2Buf(self._src, self._dst, self._type, \
                self._priority, self._row, self._col, \
                self._length, self._payload)
        return buf



ip = "192.168.1.199"
ip = "localhost"
port_pwr = 10005
port_eclss = 10006
port_str = 10004
port_agnt = 10002

sock_str = Socket(ip, port_str)
sock_pwr = Socket(ip, port_pwr)
sock_eclss = Socket(ip, port_eclss)
sock_agnt = Socket(ip, port_agnt)


EMPTY = -9999
cnt = 1

# payload -> [table_id, testBed_ts, <fdd_info>]
# Table 1 :: FDD_SPG_DUST
spg_fdd = Subsystem(src=5, dst=1, type=3, priority=7, \
                    row=1, col=6, length=6)

# Table 2 :: FDD_ECLSS_DUST
eclss_fdd_dust = Subsystem(src=6, dst=1, type=3, priority=7, \
                            row=1, col=52, length=52)

# Table 3 :: FDD_ECLSS_PAINT
eclss_fdd_paint = Subsystem(src=6, dst=1, type=3, priority=7, \
                            row=1, col=52, length=52)

# Table 4 :: FDD_STR_DMG
str_fdd = Subsystem(src=4, dst=1, type=3, priority=7,\
                    row=1, col=3, length=3)

# Table 5 :: FDD_NPG_DUST
npg_fdd = Subsystem(src=5, dst=1, type=3, priority=7, \
                    row=1, col=3, length=3)

# Table 6 :: STATES_AGENT
agnt = Subsystem(src=2, dst=1, type=3, priority=7, \
                    row=1, col=3, length=3)




sent_freq = 0.001
sleep_freq = 1
while True:
    # if cnt>10:
    #     break

    spg_fdd_payload = [1, cnt] + [0.1]*4
    spg_fdd_buf = spg_fdd.to_buffer(spg_fdd_payload)
    sock_pwr.send(spg_fdd_buf)
    print("[->] Table {} : Count - {} , sent {} bytes".format(\
        1, cnt, len(spg_fdd_buf)))
    time.sleep(sent_freq)
    

    # eclss_fdd_dust_payload = [2, cnt] + [0.1]*50
    # eclss_fdd_dust_buf = eclss_fdd_dust.to_buffer(eclss_fdd_dust_payload)
    # sock_eclss.send(eclss_fdd_dust_buf)
    # print("[->] Table {} : Count - {} , sent {} bytes".format(\
    #     2, cnt, len(eclss_fdd_dust_buf)))
    # time.sleep(sent_freq)

    # eclss_fdd_paint_payload = [3, cnt] + [0.1]*50
    # eclss_fdd_paint_buf = eclss_fdd_paint.to_buffer(eclss_fdd_paint_payload)
    # sock_eclss.send(eclss_fdd_paint_buf)
    # print("[->] Table {} : Count - {} , sent {} bytes".format(\
    #     3, cnt, len(eclss_fdd_paint_buf)))
    # time.sleep(sent_freq)

    # str_fdd_payload = [4, cnt, 0.1]
    # str_fdd_buf = str_fdd.to_buffer(str_fdd_payload)
    # sock_str.send(str_fdd_buf)
    # print("[->] Table {} : Count - {} , sent {} bytes".format(\
    #     4, cnt, len(str_fdd_buf)))
    # time.sleep(sent_freq)

    # npg_fdd_payload = [5, cnt, 0]
    # npg_fdd_buf = npg_fdd.to_buffer(npg_fdd_payload)
    # sock_pwr.send(npg_fdd_buf)
    # print("[->] Table {} : Count - {} , sent {} bytes".format(\
    #     5, cnt, len(npg_fdd_buf)))
    # time.sleep(sent_freq)

    # agnt_payload = [6, cnt, -1.0]
    # agnt_buf = agnt.to_buffer(agnt_payload)
    # sock_agnt.send(agnt_buf)
    # print("[->] Table {} : Count - {} , sent {} bytes".format(\
    #     6, cnt, len(agnt_buf)))
    # time.sleep(sent_freq)


    # bufs = [spg_fdd_buf, eclss_fdd_dust_buf, eclss_fdd_paint_buf, \
    #         str_fdd_buf, npg_fdd_buf, agnt_buf]
    
    # for table_id, buf in enumerate(bufs):
    #     sock.sendto(buf, (IP, PORT))
    #     print("[{}] Table ID: {} sent {} bytes".format(table_id, cnt, len(buf)))
    #     time.sleep(sent_freq)
    
    # time.sleep(sleep_freq)
    time.sleep(100)
    cnt += 1