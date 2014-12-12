
import sys
import os.path
import string

def extract(srt_content):
    srt_content = srt_content.replace("\r\n", "\n")
    srt_lines = srt_content.split("\n")
    
    reversal_lines = []
    txt_lines = []
    i = 0
    length = len(srt_lines)
    expect = "line_num"
    while i < length:
        line = srt_lines[i]
        line.strip()
        if expect == "line_num":
            if line.isdigit():
                expect = "time_range"
            line += "\n"
            reversal_lines.append(line)
            txt_lines = []
            
        elif expect == "time_range":
            # must be skip
            expect = "text"
            
            line += "\n"
            reversal_lines.append(line)
            txt_lines = []
            
        elif expect == "text":
            if not line or line.isspace():
                expect = "line_num"
                txt_lines.reverse()
                reversal_lines += txt_lines
                
            if line.isdigit():
                expect = "time_range"
                
                txt_lines.reverse()
                reversal_lines += txt_lines
                
                line = "\n" + line + "\n"
                reversal_lines.append(line)
                txt_lines = []
            else:
                # text
                line += "\n"
                txt_lines.append(line)
        i += 1
    print len(reversal_lines)
    txt_content = "".join(reversal_lines)
    return txt_content

def main():
    file_path = "1 - 5 - Song Form (1042).srt"
    
    if not os.path.exists(file_path):
        print "File not exists!", file_path
        return False
        
    f = open(file_path, "r")
    if not f:
        print "Open file failed!", file_path
        return False
    
    srt_content = f.read()
    f.close()
    
    txt_content = extract(srt_content)
    
    out_path = file_path + ".txt"
    f = open(out_path, "w")
    f.write(txt_content)
    f.close()
    
    return True


if __name__ == "__main__":
    main()
