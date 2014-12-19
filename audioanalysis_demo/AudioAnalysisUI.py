#!BPY

# ########################################################################
#
# Audio Analysis - Converts audio power amplitudes into Ipo curves.
# Copyright (C) 2007 Dave Jarvis (http://www.davidjarvis.ca)
#
# This program is free software; you can redistribute it and/or
# modify it under the terms of the GNU General Public License
# as published by the Free Software Foundation; either version 2
# of the License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
# MA  02110-1301, USA.
#
# ########################################################################

# """
# Name: 'Audio Analysis'
# Blender: 243
# Group: 'Misc'
# """

__author__ = 'Dave Jarvis'
__version__ = '1.1.0 2007/04/11'
__url__ = ["Dave's Blender Tools, http://davidjarvis.ca/blender/tools",
   "Support Forum, http://blenderartists.org/forum", "Blender", "elysiun"]
__email__ = ["Dave Jarvis, At no fixed e-mail address", "scripts"]
__bpydoc__ = """\
Analyses an audio file to produce an Ipo curve.
"""

import Blender
from Blender import BGL, Draw, Text

from AudioAnalysis import *
from IpoPlotter import *
from WaveFormPanel import *

X_MARGIN  = 5
X_PADDING = 10
Y_PADDING = 5
Y_MARGIN  = 2

WIDGET_HEIGHT = 25

BUTTON_WIDTH  = 70
BUTTON_HEIGHT = WIDGET_HEIGHT
NUMBER_WIDTH  = 160
NUMBER_HEIGHT = WIDGET_HEIGHT
SLIDER_WIDTH  = 225
SLIDER_HEIGHT = WIDGET_HEIGHT
STRING_WIDTH  = 160
STRING_HEIGHT = WIDGET_HEIGHT
MENU_WIDTH    = 160
MENU_HEIGHT   = WIDGET_HEIGHT

EVENT_BUTTON_PLOT      = 10
EVENT_BUTTON_POWER     = 12
EVENT_BUTTON_FREQUENCY = 14
EVENT_BUTTON_EXIT      = 20
EVENT_BUTTON_OPEN_FILE = 25
EVENT_SLIDER_POWER_MIN = 30
EVENT_SLIDER_POWER_MAX = 40
EVENT_FREQUENCY_MIN    = 50
EVENT_FREQUENCY_MAX    = 60
EVENT_SLIDER_ATTACK    = 70
EVENT_SLIDER_DECAY     = 80
EVENT_SLIDER_OSCILLATE = 90
EVENT_SLIDER_SCALE     = 100
EVENT_SLIDER_INTERVAL  = 110
EVENT_STRING_FILENAME  = 120
EVENT_STRING_IPO_NAME  = 130
EVENT_STRING_CURVE     = 135
EVENT_MENU_IPO_TYPE    = 140

EVENT_TOGGLE_RND_ATTACK    = 150
EVENT_TOGGLE_RND_DECAY     = 160
EVENT_TOGGLE_RND_OSCILLATE = 170

# These lines are used to organise the user interface.
#
lines = []

for i in range( 0, 300, WIDGET_HEIGHT + Y_PADDING ):
  lines.append( i + Y_MARGIN )

powerMin  = Draw.Create( AudioAnalysis.POWER_MIN )
powerMax  = Draw.Create( AudioAnalysis.POWER_MAX )
# freqMin   = Draw.Create( AudioAnalysis.FREQ_MIN )
# freqMax   = Draw.Create( AudioAnalysis.FREQ_MAX )
attack    = Draw.Create( 0.0 )
decay     = Draw.Create( 0.0 )
scale     = Draw.Create( 10.0 )
interval  = Draw.Create( 0.04 )
filename  = Draw.Create( 'audio.wav' )
oscillate = Draw.Create( 0.0 )

# Default to a Lamp's Energ curve
#
ipoType   = Draw.Create( 5 )
curveType = Draw.Create( 'Energ' )
curveName = Draw.Create( 'WavIpo' )

# Unused at the moment; they might not be necessary
#
rndAttack    = Draw.Create( 0 )
rndDecay     = Draw.Create( 0 )
rndOscillate = Draw.Create( 0 )

# Control the extents for the power sliders.
#
powerSliderMin = Draw.Create( AudioAnalysis.POWER_MIN )
powerSliderMax = Draw.Create( AudioAnalysis.POWER_MAX )

# Area that draws the audio wave form.
#
waveFormPanel = WaveFormPanel()

"""
DESCRIPTION
Draws the User Interface for the wave file analysis.
"""
def gui():
  global scale, interval, filename
  global curveName, ipoType, curveType
  global powerMin, powerMax, powerSliderMin, powerSliderMax
  global attack, decay, oscillate, rndAttack, rndDecay, rndOscillate
  global waveFormPanel

  BGL.glClear( BGL.GL_COLOR_BUFFER_BIT )

  # ######################################################################
  #
  # Analyse and Exit buttons
  #
  # ######################################################################
  Draw.PushButton( 'Plot', EVENT_BUTTON_PLOT, X_MARGIN, lines[0], BUTTON_WIDTH, BUTTON_HEIGHT, 'Generate Ipo curve.' )
  Draw.PushButton( 'Power', EVENT_BUTTON_POWER, X_MARGIN + BUTTON_WIDTH + X_PADDING, lines[0], BUTTON_WIDTH, BUTTON_HEIGHT, 'Display power wave form.' )
  # Draw.PushButton( 'Frequency', EVENT_BUTTON_FREQUENCY, X_MARGIN + (BUTTON_WIDTH + X_PADDING) * 2, lines[0], BUTTON_WIDTH, BUTTON_HEIGHT, 'Display frequency levels.' )
  Draw.PushButton( 'Exit', EVENT_BUTTON_EXIT, X_MARGIN + (BUTTON_WIDTH + X_PADDING) * 3, lines[0], BUTTON_WIDTH, BUTTON_HEIGHT, 'Close audio analysis tool.' )

  # ######################################################################
  #
  # Power Minimum and Maximum sliders
  #
  # ######################################################################
  powerMin = Draw.Slider( 'Power Min: ', EVENT_SLIDER_POWER_MIN, X_MARGIN, lines[1], SLIDER_WIDTH, SLIDER_HEIGHT, powerMin.val, powerSliderMin.val, powerSliderMax.val, 1, 'Minimum power amplitude data point.' )
  powerMax = Draw.Slider( 'Power Max: ', EVENT_SLIDER_POWER_MAX, X_MARGIN + SLIDER_WIDTH + X_PADDING, lines[1], SLIDER_WIDTH, SLIDER_HEIGHT, powerMax.val, powerSliderMin.val, powerSliderMax.val, 1, 'Maximum power amplitude data point.' )

  BGL.glRasterPos2i( X_MARGIN, lines[2] )
  Draw.Text( 'Power levels isolate the desired volume range.' )

  # ######################################################################
  #
  # Randomise buttons
  #
  # ######################################################################
  # rndAttack = Draw.Toggle( 'Randomise Attack', EVENT_TOGGLE_RND_ATTACK, X_MARGIN, lines[2], SLIDER_WIDTH, SLIDER_HEIGHT, rndAttack.val, 'Randomise the Attack value.' )
  # rndDecay = Draw.Toggle( 'Randomise Decay', EVENT_TOGGLE_RND_DECAY, X_MARGIN + SLIDER_WIDTH + X_PADDING, lines[2], SLIDER_WIDTH, SLIDER_HEIGHT, rndDecay.val, 'Randomise the Decay value.' )
  # rndOscillate = Draw.Toggle( 'Randomise Oscillate', EVENT_TOGGLE_RND_OSCILLATE, X_MARGIN + SLIDER_WIDTH * 2 + X_PADDING * 2, lines[2], SLIDER_WIDTH, SLIDER_HEIGHT, rndOscillate.val, 'Randomise the Oscillate value.' )

  # ######################################################################
  #
  # Attack and decay sliders
  #
  # ######################################################################
  attack = Draw.Slider( 'Attack:', EVENT_SLIDER_ATTACK, X_MARGIN, lines[3], SLIDER_WIDTH, SLIDER_HEIGHT, attack.val, 0.0, 120.0, 0, 'Time taken to reach power point.' )
  decay = Draw.Slider( 'Decay:', EVENT_SLIDER_DECAY, X_MARGIN + SLIDER_WIDTH + X_PADDING, lines[3], SLIDER_WIDTH, SLIDER_HEIGHT, decay.val, 0.0, 120.0, 0, 'Time taken to reach zero point from peak power.' )
  # oscillate = Draw.Slider( 'Oscillate:', EVENT_SLIDER_OSCILLATE, X_MARGIN + SLIDER_WIDTH * 2 + X_PADDING * 2, lines[3], SLIDER_WIDTH, SLIDER_HEIGHT, oscillate.val, 0.0, 10.0, 0, 'Exponential decay oscillation.' )

  BGL.glRasterPos2i( X_MARGIN, lines[4] )
  Draw.Text( 'Attack and decay control the leading and trailing points on either side of the data point.' )

  # ######################################################################
  #
  # Scale, interval and offset numbers
  #
  # The scale and intervals are sliders, but set to the width of a number
  # field.
  #
  # ######################################################################
  scale = Draw.Slider( 'Scale:', EVENT_SLIDER_SCALE, X_MARGIN, lines[5], NUMBER_WIDTH, NUMBER_HEIGHT, scale.val, 0.0, 1000.0, 0, 'Sets the maximum value for data points on the Ipo curve, as a ratio.' )
  interval = Draw.Slider( 'Interval:', EVENT_SLIDER_INTERVAL, X_MARGIN + NUMBER_WIDTH + X_PADDING, lines[5], NUMBER_WIDTH, NUMBER_HEIGHT, interval.val, 0.0, 120.0, 0, 'Set to the number of seconds between volume events.' )

  # ######################################################################
  #
  # Datablock name, Ipo Type, and Curve Type
  #
  # ######################################################################
  curveName = Draw.String( 'Curve Name:', EVENT_STRING_IPO_NAME, X_MARGIN, lines[6], STRING_WIDTH, STRING_HEIGHT, curveName.val, 255, 'Provide a memorable name for the curve.' )

  menuItems = ''
  count = 0

  # Dynamically create the menu for Ipo types.
  #
  for name in IpoPlotter().getIpoTypeNames():
    menuItems += name + '%x' + str( count ) + '%|'
    count += 1

  ipoType = Draw.Menu( 'Ipo Type %t%|' + menuItems, EVENT_MENU_IPO_TYPE, X_MARGIN + STRING_WIDTH + X_PADDING, lines[6], MENU_WIDTH, MENU_HEIGHT, ipoType.val, 'Name of Ipo to create.' )

  curveType = Draw.String( 'Curve Type:', EVENT_STRING_CURVE, X_MARGIN + STRING_WIDTH + X_PADDING + MENU_WIDTH + X_PADDING, lines[6], STRING_WIDTH, STRING_HEIGHT, curveType.val, 16, 'Type of Ipo curve (e.g., LocX, HorR, Energ).' )

  # ######################################################################
  #
  # Open and Filename
  #
  # ######################################################################
  Draw.PushButton( 'Open', EVENT_BUTTON_OPEN_FILE, X_MARGIN, lines[7], BUTTON_WIDTH, BUTTON_HEIGHT, 'Select a file to load.' )
  filename = Draw.String( 'Filename:', EVENT_STRING_FILENAME, X_MARGIN + BUTTON_WIDTH + X_PADDING, lines[7], STRING_WIDTH * 3 - BUTTON_WIDTH + X_PADDING, STRING_HEIGHT, filename.val, 255, 'File to load (.wav format only).' )

# Default frequency settings at the range of human hearing.
#
#  freqMin = Draw.Slider( 'Freq Min: ', EVENT_SLIDER_FREQ_MIN, X_MARGIN, lines[1], SLIDER_WIDTH, SLIDER_HEIGHT, freqMin.val, AudioAnalysis.FREQUENCY_MIN, AudioAnalysis.FREQUENCY_MAX, 0, 'Minimum frequency data point.' )
#  freqMax = Draw.Slider( 'Freq Max: ', EVENT_SLIDER_FREQ_MAX, X_MARGIN + SLIDER_WIDTH + X_PADDING, lines[1], SLIDER_WIDTH, SLIDER_HEIGHT, freqMax.val, AudioAnalysis.FREQUENCY_MIN, AudioAnalysis.FREQUENCY_MAX, 0, 'Maximum frequency data point.' )

  # Adjust the size of the wave form panel.
  #
  size = Blender.Window.GetAreaSize()
  width = size[0]
  waveFormPanel.setCanvasWidth( width - waveFormPanel.getCanvasX() - 10 )
  waveFormPanel.draw()

"""
DESCRIPTION
Changes Blender's progress bar to indicate percentage complete when analysing
the files.
"""
def updateProgressBar( percent ):
  Blender.Window.DrawProgressBar( percent, 'Analysing: ' + str( int( percent * 100.0 ) ) + '%' )

"""
DESCRIPTION
Plots an Ipo curve corresponding to the audio data.
"""
def plot():
  global powerMin, powerMax, scale, interval, filename
  global curveName, ipoType, curveType
  global attack, decay, oscillate, rndAttack, rndDecay, rndOscillate

  try:
    aa = AudioAnalysis( filename.val )

    #aa.setFreqMin( freqMin.val )
    #aa.setFreqMax( freqMax.val )
    aa.setSamplePowerMin( powerMin.val )
    aa.setSamplePowerMax( powerMax.val )
    aa.setFrameInterval( interval.val )
    aa.analysePower( updateProgressBar )

    Blender.Window.DrawProgressBar( 1, 'Analysis Complete' )

    plotter = IpoPlotter()

    # Options that define the data values.
    #
    plotter.setPlotData( aa.getFramePower() )
    plotter.setDataPointMin( powerMin.val )
    plotter.setDataPointMax( powerMax.val )

    # Options that massage, manipulate and transform the data.
    #
    plotter.setAttack( attack.val )
    plotter.setDecay( decay.val )
    plotter.setScale( scale.val )
    plotter.setOscillateDecay( oscillate.val )
    plotter.setRandomAttack( rndAttack.val )
    plotter.setRandomDecay( rndDecay.val )
    plotter.setRandomOscillate( rndOscillate.val )

    # Options that control the Blender-based Ipo curve.
    #
    plotter.setFramesPerSecond( aa.FRAMES_PER_SECOND )
    plotter.setIpoType( ipoType.val )
    plotter.setCurveType( curveType.val )
    plotter.setName( curveName.val )

    # Plot the data points along an Ipo curve.
    #
    plotter.plot()

  except IOError:
    Draw.PupMenu( 'Error!%t|File not found:' + filename.val )

  except Exception, inst:
    Draw.PupMenu( 'Error!%t|Oops! Check console for error message.' )
    print inst

"""
DESCRIPTION
Capture events from the keyboard.

PARAMETERS
event - The keyboard identifier
value - Unknown
"""
def keyEvent( event, value ):
  if event == Draw.ESCKEY or event == Draw.QKEY:
    if Draw.PupMenu( 'Stop!%t|Exit Audio Analysis? %x1' ) == 1:
      Draw.Exit()

  Blender.Redraw()

"""
DESCRIPTION
Captures non-keyboard events from the user interface.

PARAMETERS
event - The button identifier
"""
def buttonEvent( event ):
  global waveFormPanel, powerMin, powerMax, filename

  if event == EVENT_BUTTON_PLOT:
    plot() 
  elif event == EVENT_BUTTON_OPEN_FILE:
    buttonOpen()
  elif event == EVENT_SLIDER_POWER_MIN:
    sliderPowerMin()
  elif event == EVENT_SLIDER_POWER_MAX:
    sliderPowerMax()
  elif event == EVENT_BUTTON_POWER:
    buttonPower()
  elif event == EVENT_BUTTON_EXIT:
    Draw.Exit()

"""
DESCRIPTION
Called after opening a file (callback function for Window.FileSelector).

PARAMETERS
name - The name of the file that was selected to open.
"""
def fileSelector( name ):
  global filename
  filename.val = name

"""
DESCRIPTION
Called after clicking Open to select a file.
"""
def buttonOpen():
  Blender.Window.FileSelector( fileSelector, 'Open' )

"""
DESCRIPTION
Called after clicking the "Power" button. This ties together the wave form
power analysis with the user interface. It is where the actual analysis
begins.
"""
def buttonPower():
  global waveFormPanel, filename
  global powerMin, powerMax, powerSliderMin, powerSliderMax

  waveFormPanel.clearSamples()

  aa = AudioAnalysis( filename.val )
  aa.doPowerAnalysis( waveFormPanel.appendSample, updateProgressBar )
  waveFormPanel.setPeaks( aa.getPowerMin(), aa.getPowerMax() )
  waveFormPanel.draw()

  powerMin.val = powerSliderMin.val = aa.getPowerMin()
  powerMax.val = powerSliderMax.val = aa.getPowerMax()

  Draw.Redraw( 1 )

"""
DESCRIPTION
Called after changing the minimum power level, using a slider.
"""
def sliderPowerMin():
  global waveFormPanel, powerMin, powerMax

  if powerMin.val >= powerMax.val:
    powerMin.val = powerMax.val - (1 * waveFormPanel.getPixelScale())

  waveFormPanel.setMinThreshold( powerMin.val )

"""
DESCRIPTION
Called after changing the maximum power level, using a slider.
"""
def sliderPowerMax():
  global waveFormPanel, powerMin, powerMax

  if powerMax.val <= powerMin.val:
    powerMax.val = powerMin.val + (1 * waveFormPanel.getPixelScale())

  waveFormPanel.setMaxThreshold( powerMax.val )

# Make Blender aware of this script.
#
Draw.Register( gui, keyEvent, buttonEvent )

