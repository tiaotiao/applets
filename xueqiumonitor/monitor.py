#-*- coding: UTF-8 -*- 

import os
import sys
import time
import logging
from datetime import date
from datetime import datetime
from BeautifulSoup import BeautifulSoup
import ConfigParser

global cfg
global ids
global url_prefix
global log_level

def getHTML(url):
    import httplib
    import urllib2
    
    try:
        headers={'User-Agent':'Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US; rv:1.9.1.4) Gecko/20091201 Firefox/3.5.6'}
        request = urllib2.Request(url, headers = headers)
        response = urllib2.urlopen(request, timeout = 5)
    except urllib2.URLError, e:
        if hasattr(e, "code"):
            return None, e.code
        elif hasattr(e, "reason"):
            return None, e.reason
    except httplib.BadStatusLine, e:
        return None, "BadStatusLine"
    
    import re
    charset = "utf-8"

    try:
        info = response.info()
        if info.has_key("content-type"):
            pat = r"charset=\s*(\S+)"
            m = re.search(pat, info["content-type"], re.I)
            if m:
                charset = m.group(1)
            else:
                print "ERROR: html charset not found"
    except:
        print "ERROR:", traceback()

    
    try:
        content = response.read()
    except:
        e = "timeout"
        print "ERROR:", traceback()
        return None, e
    
    if charset:
        try:
            content = content.decode(charset)
        except:
            print "ERROR: can not decode html content from charset=", str(charset)
    
    return content, None
    

def getTargetContent(content):
    soup = BeautifulSoup(content)   # TODO

    cubeWeight = soup.body.findAll("div", {"id":"cube-weight"})
    if len(cubeWeight) != 1:
        return None
    cubeWeight = cubeWeight[0]
    
    weightlist = cubeWeight.findAll("div", {"class":"weight-list"})
    if len(weightlist) != 1:
        return None
    
    return weightlist[0]
    
def messagebox(title, msg):
    import ctypes
    MessageBox = ctypes.windll.user32.MessageBoxA
    MessageBox(None, msg, title, 0)

def getList(pid):
    url = url_prefix + pid
    html_content, err = getHTML(url)
    if not html_content:
        print "ERROR: get html failed", err
        return
    list_content = getTargetContent(html_content)
    if not list_content:
        print "ERROR: get target content failed", html_content.encode("gbk")
        return
    return list_content
    
def now():
    return time.strftime("[%H:%M:%S]", time.localtime())
    
def getlist(content):
    list = content.findAll("span", {"class": "stock-name"})
    l = ""
    for v in list:
        s = v.contents[0].encode("gbk")
        l = l + s + " "
    return l
        
def diff(last, content):    
    listNew = content.findAll("span", {"class": "stock-name"})
    listOld = last.findAll("span", {"class": "stock-name"})
   
    list = ""
    news = ""
    
    #if listOld == listNew:
    #    return None
    
    for v in listNew:
        s = v.contents[0].encode("gbk")
        list = list + s + " "
        
        find = False
        for u in listOld:
            if u == v:
                find = True
                break
        if not find:
            news = news + s + " "
    
    return news

class Monitor(object):
    
    def __init__(self, pid):
        super(Monitor, self).__init__()
        self.pid = pid
        self.last_content = getList(pid)
        l = getlist(self.last_content)
        print "Start:", pid, ":", l

    def check(self):
        content = getList(self.pid)

        news = diff(self.last_content, content)
        if news:
            print "Update:", pid, ">>>>", news
            messagebox("Monitor", "There is an update:\n"+news)

        l = getlist(content)

        self.last_content = content

        return news
        

def loadConfig(filepath):
    global cfg
    global url_prefix
    global ids
    global log_level

    config = ConfigParser.ConfigParser()
    config.read(filepath)

    cfg = {}

    url_prefix = config.get("monitor", "url_prefix")
    if not url_prefix:
        print "ERROR: config [monitor:prefix] not found"
        return None
    cfg["prefix"] = url_prefix

    pid = config.get("monitor", "id")
    if not pid:
        print "ERROR: config [monitor:id] not found"
        return None
    ids = pid.split(",")
    cfg["ids"] = ids

    interval = config.getint("monitor", "interval")
    if not interval:
        interval = 10    
    cfg["interval"] = 10

    return cfg

    
def run():

    loadConfig("./config.ini")
    
    # init
    monitors = []
    for pid in ids:
        m = Monitor(pid)
        monitors.append(m)

    print "Start run ..."

    while True:
        time.sleep(10)

        print "Check."

        for m in monitors:
            try:
                m.check()
            except Exception as e:
                print "Exception:", e

def main():
    run()
    
if __name__ == "__main__":
    main()
