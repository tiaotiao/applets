
import sys, os
import wave, audioop
import AudioAnalysis


# callback: onFrame(frameNumber, power, totalFrames, rate, params)
def analysisWavePower(file_path, onFrame):
    
    if not os.path.exists(file_path):
        print "File not exists!", file_path
        return
    f = open(file_path, "rb")
    if not f:
        print "Open file failed!"
        return
    try:
        wave_file = wave.open(f)
    except Exception, e:
        print e
        print "File type is not wav!"
        return
    
    frame_rate = wave_file.getframerate()
    frame_count = wave_file.getnframes()
    frame_width = wave_file.getsampwidth()
    channels = wave_file.getnchannels()
    params = wave_file.getparams()

    # start
    frame = wave_file.readframes( 1 )
    frameNumber = 0
    
    while len( frame ):
        #progress = float( frameNumber ) / self.getSamples()
        power = audioop.rms( frame, frame_width )
        frameNumber, power  # TODO
        
        # callback
        onFrame(frame, power, frameNumber, frame_count, frame_rate, params)

        frame = wave_file.readframes( 1 )
        frameNumber += 1
        
    wave_file.close()

def frameToSecond(frame_count, rate):
    return frame_count / (1.0 * rate)

def createSplitWaveFunc(out_dir, power_threshold = 800, silence_time = 0.2, min_split_time = 0.1):
    if not os.path.isdir(out_dir):
        os.makedirs(out_dir)
    if not os.path.isdir(out_dir):
        print "Create split wave dir failed. %s" % out_dir
        return None
    
    # init
    value = {}
    value["status"] = "slience"  # or "recording"
    value["buffer"] = []
    value["split_files"] = []
    value["record_start_time"] = 0
    value["last_sound_time"] = 0
    value["rate"] = 0
    value["params"] = None
    value["max_power"] = 0
    
    def flushBufferToFile():
        num = len(value["split_files"])
        start_time = value["record_start_time"]
        length = frameToSecond(len(value["buffer"]), value["rate"])
        if length < min_split_time:
            print "flush: %.4lf len=%.4lf skiped" % (start_time, length)
            return
        
        out_filename = "%03d_%.4lf.wav" % (num + 1, start_time)
        out_path = os.path.join(out_dir, out_filename)
        value["split_files"].append(out_path)
        
        print "flush: %.4lf len=%.4lf" % (start_time, length)
        
        out_file = wave.open(out_path, "wb")
        out_file.setparams(value["params"])
        
        for frame in value["buffer"]:
            out_file.writeframes(frame)
        out_file.close()
        
        # clear
        value["buffer"] = []
    
    # process function
    def splitWave_onFrameFunc(frame, power, frame_num, total_frames, rate, params):
        
        value["rate"] = rate
        value["params"] = params
        current_time = frameToSecond(frame_num, rate)
        hasSound = power > power_threshold
        
        if value["max_power"] < power:
            value["max_power"] = power
        
        if frame_num % 100 == 0:
            pass
            #print "%d %.5lf: \t[%d] %s max=%d" % (frame_num, frameToSecond(frame_num, rate), power, hasSound, value["max_power"])
        
        if hasSound or value["status"] == "recording":
            value["buffer"].append(frame)
        
        if hasSound:
            if value["status"] == "slience":
                value["record_start_time"] = current_time
            value["status"] = "recording"
            value["last_sound_time"] = current_time
        else:
            if value["status"] == "recording":
                if current_time - value["last_sound_time"] >= silence_time:
                    value["status"] = "slience"
                    flushBufferToFile()
    
        if frame_num == total_frames - 1:   # last frame
            if len(value["buffer"]):
                flushBufferToFile()
    
    return splitWave_onFrameFunc
    
def main():
    file_path = "demo.wav"
    ss = file_path.split(".")
    onFrame = createSplitWaveFunc(ss[0]+"_split")
    
    analysisWavePower(file_path, onFrame)

if __name__ == "__main__":
    main()
    #test()
    