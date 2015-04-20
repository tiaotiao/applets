
"""
1. Download and install py2exe from here. 
http://sourceforge.net/projects/py2exe/files/py2exe/0.6.9/

2. Build an exe file with the cmd in the current path:
>> python setup.py py2exe

3. A new folder named 'dist' will be generated. Feel free to pack and distribute it.
"""

from distutils.core import setup 
import py2exe 

setup(console=["monitor.py"])
