
import urllib2

'wget --user-agent="Mozilla/5.0" --post-file=hello.flac --header="Content-Type: audio/x-flac; rate=16000" "http://www.google.com/speech-api/v1/recognize?xjerr=1&client=chromium&lang=zh-CN&maxresults=1"'

API_URL = "https://www.google.com/speech-api/v1/recognize?xjerr=1&client=chromium&lang=en-US&results=5"

FILE = "hello.flac"

def test():
    f = open(FILE, 'rb').read()
    if not f:
        print "Open file error!"
    headers = {'Content-Type':'audio/x-flac; rate=44100'}
    req = urllib2.Request(API_URL, f, headers)
    response = urllib2.urlopen(req)
    print "[", response.read().decode('UTF-8'), "]"

def main():
    test()
    
if __name__ == "__main__":
    main()
