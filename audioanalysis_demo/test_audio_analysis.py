

import sys, wave
import AudioAnalysis

FILE_NAME = "snippet.wav"

def testWavWrite():
    try:
        f = wave.open(FILE_NAME, "rb")
    except Exception, e:
        print e
        print "File type is not wav!"
        return
    c = wave.open("cnv_" + FILE_NAME, "wb")
        
    print f.getnchannels()
    print f.getsampwidth()
    print f.getframerate()
    print f.getnframes()
    #print f.getparams()
    
    total = f.getnframes()
    read_count = total / 2
    
    c.setnchannels(f.getnchannels())
    c.setsampwidth(f.getsampwidth())
    c.setframerate(f.getframerate())
    c.setnframes(read_count)
    c.setcomptype(f.getcomptype(), f.getcompname())
    
    frames = f.readframes(read_count)
    
    print len(frames)
    print "bytes per frame: ", len(frames) / read_count
    
    #for b in frames:
    #    i = int(b.encode("hex"), 16)
    #    print b.encode("hex")
        #print '#' * (i / 10)
    
    c.writeframes(frames)
    
    print "----------"
    
    f.close()
    c.close()

def process(p):
    print p

def testAudioAnalysis():
    
    a = AudioAnalysis.AudioAnalysis(FILE_NAME)
    
    print a.getFilename()
    print a.getFileType()
    
    a.setFrameInterval(0.01)
    print a.analysePower(process)
    
    print a.getPowerMin(), "\tgetPowerMin"
    print a.getPowerMax(), "\tgetPowerMax"
    print a.getSamplePowerMin(), "\tgetSamplePowerMin"
    print a.getSamplePowerMax(), "\tgetSamplePowerMax"
    print a.getFrameRate(), "\tgetFrameRate"
    print a.getSampleWidth(), "\tgetSampleWidth"
    print a.getDuration(), "\tgetDuration"
    print a.getFrameInterval(), "\tgetFrameInterval"
    print a.getSamples(), "\tgetSamples"
    
    powers = a.getFramePower()
    for p in powers:
        print "%04lf" % p[0], "%-6d" % p[1] ,'#' * (p[1] / 100)

def main():
    
    f = open(FILE_NAME, "rb")
    if not f:
        print "Open file failed!"
        return
    try:
        w = wave.open(f)
    except Exception, e:
        print e
        print "File type is not wav!"
        return
        
    print "get channels\t", w.getnchannels()  # channels, single or double
    print "frame rate\t", w.getframerate()  # rate, frames per sec
    print "samp width\t", w.getsampwidth()  # maybe:  channels * width = bytes per frame
    print "get n frames\t", w.getnframes()  # total frames 
    print "comp type\t", w.getcomptype()   # compress
    print "params\t", w.getparams()
    
    total = w.getnframes()
    read_count = 100
    
    frames = w.readframes(read_count)
    
    print "len(frames)\t", len(frames)
    print "bytes per frame\t", len(frames) / read_count
    
    #for b in frames:
        #i = int(b.encode("hex"), 16)
        #print b.encode("hex")
        #print '#' * (i / 10)
    
    print "----------"
    
    w.close()
    f.close()
    




if __name__ == "__main__":
    main()
    #testAudioAnalysis()
    #testWavWrite()
    