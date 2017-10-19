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
#field_cmd="cmd"
field_key="key"
field_value="value"

cmd_type_get="get"
cmd_type_set="set"

svr_ip="127.0.0.1"
svr_port=9090

CMD_TYPE_GET=0
CMD_TYPE_GET_ACK=1
CMD_TYPE_SET=2
CMD_TYPE_SET_ACK=3

def get_msg():
    msgdict={}
    #infodict={}
    #infodict[field_key]="testkey"
    #msgdict["cmd"]=cmd_type_get
    #msgdict["info"]=infodict
    msgdict[field_key]="testkey1"
    msg=json.dumps(msgdict)
    length=len(msg)+8
    str = struct.pack('!II', length,CMD_TYPE_GET)
    str=str+msg
    #print len(msg),len(str)
    return str
def produce_msg_get(key):
    msgdict={}
    msgdict[field_key]=key
    msg=json.dumps(msgdict)
    length=len(msg)+8
    str = struct.pack('!II', length,CMD_TYPE_GET)
    str=str+msg
    return str

def produce_msg_set(key,value):
    msgdict={}
    msgdict[field_key]=key
    msgdict[field_value]=value
    msg=json.dumps(msgdict)
    length=len(msg)+8
    str = struct.pack('!II', length,CMD_TYPE_SET)
    str=str+msg
    return str

def parse_msg(buffer):
    headbytes=buffer[0:8]
    length, cmdtype= struct.unpack('!II',headbytes)
    info=buffer[8:]
    return length, cmdtype, info

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
    #time.sleep(1)
    #sock.close()
    msgrecv=sock.recv(10000)
    print "msgrecv: ",len(msgrecv)
    length,cmdtype, info=parse_msg(msgrecv)
    print length, cmdtype, info
    sock.close()
    return

def client_short_set_get():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((svr_ip, svr_port))
    msg=produce_msg_set("testkey1","testvalue1")
    sock.send(msg)
    #time.sleep(1)
    #sock.close()
    msgrecv=sock.recv(10000)
    length,cmdtype, info=parse_msg(msgrecv)
    print "set finish ",length, cmdtype, info
    #return
    msg=produce_msg_get("testkey1")
    sock.send(msg)
    msgrecv=sock.recv(10000)
    length,cmdtype, info=parse_msg(msgrecv)
    print "get finish ",length, cmdtype, info
    sock.close()
    return


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
    #multi_client()
    #client_req_reply()
    #client_short()
    client_short_set_get()

































