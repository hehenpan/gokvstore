#-*-coding:UTF-8-*-
import sys
import json
import struct
import socket
import time
import gevent

import gokvstorecli

#from gevent.coros import BoundedSemaphore
from gevent import monkey


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


def test():
    cli = gokvstorecli.gokvstoreclient('127.0.0.1',9090)
    cli.connect()
    print cli.set('testkey','testvalue1')
    print cli.get('testkey')
    cli.close()


if __name__=="__main__":
    print "hello"
    test()

    #client()
    #multi_client()
    #client_req_reply()
    #client_short()
   # client_short_set_get()

































