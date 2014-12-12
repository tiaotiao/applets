# coding:utf-8
import sys
import os.path
import string
import re
reload(sys)
sys.setdefaultencoding('utf-8')

def is_zh(str):
    len1 = len(str)
    len2 = len(unicode(str, 'utf-8'))
    #print len1 != len2, str
    return len1 != len2

def remove_font_label(line):
    p = "<.*?>"
    return re.sub(p, '', line)
    
def remove_en_lines(txt_lines):
    ret_lines = []
    status = "en"
    for line in txt_lines[::-1]:
    #for line in txt_lines:
        if status == "en":
            if is_zh(line):
                ret_lines.insert(0, line)
                #ret_lines.append(line)
                status = "zh"
        else:
            ret_lines.insert(0, line)
            #ret_lines.append(line)
    if not ret_lines:
        return txt_lines
    return ret_lines
    
def convert(srt_content):
    srt_content += "\n"
    srt_content = srt_content.replace("\r\n", "\n")
    srt_lines = srt_content.split("\n")
    
    cnv_lines = []
    txt_lines = []
    i = 0
    length = len(srt_lines)
    expect = "line_num"
    
    while i < length:
        line = srt_lines[i]
        #print ">>> ", expect, "\t", line
        if expect == "line_num":
            if line.isdigit():
                cnv_lines.append("\n")
                expect = "time_range"
            cnv_lines.append(line)
        elif expect == "time_range":
            # must be skip
            cnv_lines.append(line)
            expect = "text"
        elif expect == "text":
            if not line or line.isspace() or i == length - 1:
                expect = "line_num"
                # append txt_lines
                #print "-------------"
                #print txt_lines
                zh_lines = remove_en_lines(txt_lines)
                cnv_lines += zh_lines
                txt_lines = []
            else:
                # text
                line = remove_font_label(line)
                txt_lines.append(line)
        i += 1
        
    i = 0
    for line in cnv_lines:
        cnv_lines[i] = line.strip() + "\n"
        i += 1
    
    txt_content = "".join(cnv_lines)
    return txt_content

def check_file_encode(content):
    if content.startswith("\xEF\xBB\xBF"):
        content = content.lstrip("\xEF\xBB\xBF")
    try:
        content.decode('utf-8')
    except Exception, e:
        u = content.decode('gbk')
        content = u.encode('utf-8')
    return content
    
    
def convert_file(file_path):
    #file_path = check_file_encode(file_path)
    file_path = file_path.encode('gbk')
    if not os.path.exists(file_path):
        print "File not exists!", file_path
        return False
        
    f = open(file_path, "rb")
    if not f:
        print "Open file failed!", file_path
        return False
    
    srt_content = f.read()
    f.close()
    
    srt_content = check_file_encode(srt_content)
    
    
    txt_content = convert(srt_content)
    
    out_path = file_path + "_rm_en.srt"
    f = open(out_path, "w")
    f.write(txt_content)
    f.close()
    
    return True

def convert_dir(dir_path):
    files = os.listdir(dir_path)
    for file in files:
        _, ext = os.path.splitext(file)
        if ext.lower() == ".srt":
            path = os.path.join(dir_path, file)
            #print path
            print path
            #convert_file(path)
            
        
def main():
    #convert_dir(ur"E:\Klive\GMC\sw3\prefix")
    convert_file(ur"E:\Desktop\for_use_coursera-guitar-003_srt-45_en.srt")

if __name__ == "__main__":
    main()
    #print remove_font_label("<asd>abc<sdfs>")
    
