#-*-coding:UTF-8-*-
import sys
import json
import struct
import socket
import time
import gevent
#from gevent.coros import BoundedSemaphore
from gevent import monkey
gevent.monkey.patch_all()
field_cmd="cmd"
field_key="key"
field_value="value"

cmd_type_get="get"
cmd_type_set="set"

svr_ip="127.0.0.1"
svr_port=9090

def get_msg():
    msgdict={}
    infodict={}
    infodict[field_key]="testkey"
    msgdict["cmd"]=cmd_type_get
    msgdict["info"]=infodict
    msg=json.dumps(msgdict)
    length=len(msg)+4
    str = struct.pack('!I', length)
    str=str+msg
    #print len(msg),len(str)
    return str

def client():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((svr_ip, svr_port))
    count=0
    while True:
        msg=get_msg()
        sock.send(msg)
        count=count+1
        if count%1000==1:
            time.sleep(1000)
            print 'send data'

def client_short():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((svr_ip, svr_port))
    msg=get_msg()
    sock.send(msg)
    time.sleep(1)
    sock.close()

def client_req_reply():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((svr_ip, svr_port))
    msg=get_msg()
    msg=msg+msg+msg+msg+msg+msg+msg+msg+msg+msg+msg+msg+msg
    sock.send(msg)
    data=sock.recv(1000)
    print data
    sock.close()

def multi_client():

    count=1
    while True:
        eventlist=[]
        for item in range(0,100,1):
            #if count%100==0:
            #    time.sleep(6)
            #event = gevent.spawn(client_short)
            event = gevent.spawn(client_req_reply)
            eventlist.append(event)
        #time.sleep(3)
        #break
        gevent.joinall(eventlist)
        count=count+1
        print 'send data ',count


if __name__=="__main__":
    print "hello"
    #client()
    multi_client()
    #client_req_reply()


































