# Purpose: To construct payload 
# Date Created: 1 March 2021
# Date Last Modified: 1 March 2021
# Modeler Name: Jiachen Wang (UConn)
# Funding Acknowledgement: Funded by the NASA RETHi Project

from ctypes import *

class Header(Structure):
    _fields_ = [("src", c_uint8),
                ("dst", c_uint8)]

class Packet:
    def __init__(self):
        pass
    
    # payload is a double list
    def pkt2Buf(self, _src, _dst):
        header_buf = Header(_src, _dst, )
        buf = bytes(header_buf)
        return buf
        
    def buf2Pkt(self, buffer):
        self.header = Header.from_buffer_copy(buffer[:8])
        double_arr = c_double * self.header.length
        self.payload = double_arr.from_buffer_copy(buffer[8:8+8*self.header.length])

# usage:
# _src = 1
# _dst = 2
# _type = 3
# _priority = 4
# _row = 1
# _col = 4
# _length = 5
# _payload = [1, 2.2, -3, 400, -800]

# pkt = Packet()
# buf = pkt.pkt2Buf(_src, _dst, _type, _priority, _row, _col, _length, _payload)
# print(buf)

# pkt2 = Packet()
# pkt2.buf2Pkt(buf)
# print(pkt2.header.src, pkt2.payload[2])

