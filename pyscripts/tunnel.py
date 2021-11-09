import socket
import json
# redirect output to tcp tunnel (a tcp client connection)


class Tunnel:
    def __init__(self, ip, port):
        self.ip = ip
        self.port = port
        self.conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.conn.connect((self.ip, self.port))
        print("established tunnel")

    # msg is a json/dict
    # {
    #   type: , -1-heartbeat, 0-normal log, 1-statistics
    #   msg: , string
    #   stats_dr: int array, statistics for dr (entries counter of each table, ordered as table id)
    #   stats_c2: statistics for c2, format not determined yet
    # }
    def send(self, msg):
        self.conn.send(bytes(json.dumps(msg), encoding="utf-8"))

# tunnel = Tunnel("localhost", 60001)
# tunnel.send(msg)
