
#-*- coding: UTF-8 -*- 

import os
import sys
import time
from datetime import date
from datetime import datetime
from BeautifulSoup import BeautifulSoup
    
def getHTML(url):
    import httplib
    import urllib2
    
    try:
        headers={'User-Agent':'Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US; rv:1.9.1.6) Gecko/20091201 Firefox/3.5.6'}
        request = urllib2.Request(url, headers = headers)
        response = urllib2.urlopen(request, timeout = 5)
    except urllib2.URLError, e:
        if hasattr(e, "code"):
            #log.error("Get html failed, server error: %s" % (str(e.code)))
            return None, e.code
        elif hasattr(e, "reason"):
            #log.error("Get html failed, can not reach server: %s" % (str(e.reason)))
            return None, e.reason
    except httplib.BadStatusLine, e:
        #log.error("Get html failed, BadStatusLine")
        return None, "BadStatusLine"
    
    import re
    charset = None
    
    try:
        info = response.info()
        if info.has_key("content-type"):
            pat = r"charset=\s*(\S+)"
            m = re.search(pat, info["content-type"], re.I)
            if m:
                charset = m.group(1)
            else:
                log.warn("html charset not found")
    except:
        log.warn(traceback())
    
    try:
        content = response.read()
    except:
        e = "timeout"
        log.error(traceback())
        return None, e
    
    if charset:
        try:
            content = content.decode(charset)
        except:
            log.error("can not decode html content from charset=%s" % (str(charset)))
    
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
    
def getList():
    url = "http://xueqiu.com/P/ZH000979"
    html_content, err = getHTML(url)
    if not html_content:
        print("get html failed", err)
        return
    list_content = getTargetContent(html_content)
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
            print "\tnew", s
            news = news + s + " "
    
    return news
    
def run():
    last_content = getList()
    
    l = getlist(last_content)
    
    print now(), "start run ..."
    print l
    
    
    while True:
        time.sleep(2)
        
        try:
            content = getList()
            
            news = diff(last_content, content)
            if not news:
                pass
                #print now(), "ok"
            else:
                print now(), "update!"           
                messagebox("Wooow!", "There is an update!\n"+news)
                
            last_content = content
        except Exception as e:
            print "exception:", e

def main():
    run()
    
if __name__ == "__main__":
    main()
