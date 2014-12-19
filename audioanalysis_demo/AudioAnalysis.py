import sys, wave, audioop

"""
DESCRIPTION
This class provides functionality to read a .WAV file and filter the
data samples within it.
"""
class AudioAnalysis:
  BITS_PER_BYTE = 8
  FRAME_BUFFER_SIZE = 16384
  FRAMES_PER_SECOND = 25

  POWER_MAX = 32768
  POWER_MIN = 0

  __filename = ''
  __powerMin = POWER_MAX
  __powerMax = POWER_MIN

  # When adding samples into __framePower, these variables determine whether
  # or not the value will be added. This occurs during analysis.
  #
  __samplePowerMin = POWER_MIN
  __samplePowerMax = POWER_MAX

  __frameInterval = 1.0 / FRAMES_PER_SECOND

  __frameRate   = 0
  __samples     = 0
  __sampleWidth = 0
  __channels    = 0

  __framePower = []

  def __init__( self, filename ):
    self.setFilename( filename )

  """
  DESCRIPTION
  Sets the name of the file to analyse.
  """
  def setFilename( self, filename ):
    self.__filename = filename

  """
  DESCRIPTION
  Returns the name of the file to analyse.
  """
  def getFilename( self ):
    return self.__filename

  """
  DESCRIPTION
  Returns the type of file, using its filename extension.
  """
  def getFileType( self ):
    return str.split( self.getFilename(), '.' )[ -1 ].lower()

  """
  DESCRIPTION
  Sets the minimum power limit determined by the signal.
  """
  def setPowerMin( self, power ):
    self.__powerMin = power

  """
  DESCRIPTION
  Returns the minimum power limit determined by the signal.
  """
  def getPowerMin( self ):
    return self.__powerMin

  """
  DESCRIPTION
  Sets the maximum power limit determined by the signal.
  """
  def setPowerMax( self, power ):
    self.__powerMax = power

  """
  DESCRIPTION
  Returns the maximum power limit determined by the signal.
  """
  def getPowerMax( self ):
    return self.__powerMax

  """
  DESCRIPTION
  Sets the maximum level for power values added to the samples.
  """
  def setSamplePowerMax( self, power ):
    self.__samplePowerMax = power

  """
  DESCRIPTION
  Returns the maximum level for power values added to the samples.
  """
  def getSamplePowerMax( self ):
    return self.__samplePowerMax

  """
  DESCRIPTION
  Sets the minimum level for power values added to the samples.
  """
  def setSamplePowerMin( self, power ):
    self.__samplePowerMin = power

  """
  DESCRIPTION
  Returns the minimum level for power values added to the samples.
  """
  def getSamplePowerMin( self ):
    return self.__samplePowerMin

  def setFrameRate( self, rate ):
    self.__frameRate = rate

  def getFrameRate( self ):
    return self.__frameRate

  def setSamples( self, samples ):
    self.__samples = samples

  def getSamples( self ):
    return self.__samples

  def setSampleWidth( self, sampleWidth ):
    self.__sampleWidth = sampleWidth

  def getSampleWidth( self ):
    return self.__sampleWidth

  """
  DESCRIPTION
  Controls the rate at which data samples are added into the list of
  frame-power tuples. Although sound files can run at 44100 Hz, video
  files in Blender run at 25 frames per second. Making this value smaller
  than 1/25 (0.04) is likely not necessary.
  
  If you know roughly how often the sound occurs within the given power
  limits (e.g., every three seconds), then you can change this value
  accordingly (e.g., 3).
  """
  def setFrameInterval( self, sfr ):
    self.__frameInterval = sfr

  def getFrameInterval( self ):
    return self.__frameInterval

  def setChannels( self, channels ):
    self.__channels = channels

  def getChannels( self ):
    return self.__channels

  def getFramePower( self ):
    return self.__framePower

  def __resetFramePower( self ):
    self.__framePower = []

  """
  DESCRIPTION
  Returns the length of the media, in seconds.
  """
  def getDuration( self ):
    return float( self.getSamples() ) / float( self.getFrameRate() )

  """
  DESCRIPTION
  Returns true if the given amount of power equals or exceeds the minimum
  power limit. If the minimum power limit has not been set, this
  will always return true.
  """
  def aboveSampleMin( self, power ):
    min = self.getSamplePowerMin()

    if min == -1:
      power = min

    return power >= min

  """
  DESCRIPTION
  Returns true if the given amount of power is less than or equal to the
  maximum power limit. If the maximum power limit has not been set, this
  will always return true.
  """
  def belowSampleMax( self, power ):
    max = self.getSamplePowerMax()

    if max == -1:
      max = power

    return power <= max

  """
  DESCRIPTION
  Returns true if the given amount of power is between the minimum and
  maximum power limits.
  """
  def withinSamplePowerLimit( self, power ):
    return self.aboveSampleMin( power ) and self.belowSampleMax( power )

  """
  DESCRIPTION
  Changes the minimum or maximum power limit if the given amount of power
  exceeds either limit.
  """
  def __setPowerLimits( self, power ):
    if power < self.getPowerMin(): self.setPowerMin( power )
    if power > self.getPowerMax(): self.setPowerMax( power )

  """
  DESCRIPTION
  Returns the time, in seconds, that a specific frame occurred.
  """
  def __calculateTime( self, frame ):
    return float( frame ) / float( self.getSamples() ) * self.getDuration()

  """
  DESCRIPTION
  Store the values that fit the limits.
  """
  def __appendFramePower( self, frame, power ):
    if self.withinSamplePowerLimit( power ):
      t = self.__calculateTime( frame )
      framePower = self.getFramePower()

      if framePower:
        tuple = framePower[-1]

        # Only append the tuple if the difference in time between this power
        # sample and the previous frame exceeds the sample frame rate.
        #
        if t - tuple[0] >= self.getFrameInterval():
          framePower.append( (t, power) )
      else:
        framePower.append( (t, power) )

  """
  DESCRIPTION
  Creates a tally of tuples: time (in seconds) and its associated power level.
  """
  def analysePower( self, progress ):
    self.__resetFramePower()
    self.doPowerAnalysis( self.__appendFramePower, progress )

  """
  DESCRIPTION
  Used to perform various operations on the power values within the media file.
  The method 'doWork' is called every frame with both the frame number and
  its corresponding amount of power.
  """
  def doPowerAnalysis( self, doWork, progress ):
    filename = self.getFilename()

    try:
      signal = wave.open( filename, 'rb' )

      # Record the frame data so that power values can be associated with
      # an interval of seconds into the audio stream.
      #
      self.setFrameRate( signal.getframerate() )
      self.setSamples( signal.getnframes() )
      self.setSampleWidth( signal.getsampwidth() )
      self.setChannels( signal.getnchannels() )

      width = self.getSampleWidth()

      # There must be a way to do this more efficiently ...?
      #
      frame = signal.readframes( 1 )
      frameNumber = 0

      """
      Somewhere around this loop is where you would want to use a windowing
      technique and Fast Fourier Transform to ascertain the frequencies that
      are running through the data set.

      The windowing should probably take place over groups of 1764 frames.
      The result of the FFT would be a list of frequencies for the data
      set.
      
      1764 = 44100 Hz / 25 FPS
      """
      while len( frame ):

        # 16383 is arbitrary; but it is one less than a power of two, which
        # allows this condition to be true ... sometimes. In other words, this
        # controls how often the progress bar is updated.
        #
        if (frameNumber & 16383) == 0:
          progress( (float( frameNumber ) / self.getSamples()) )

        # To keep this method generic (i.e., loop through a file and
        # determine the power value for a specific frame), the following line
        # calls a function with a frame number and a power level. See
        # analysePower for example usage.
        #
        power = audioop.rms( frame, width )
        doWork( frameNumber, power )

        # Track the minimum and maximum values found in the file.
        #
        self.__setPowerLimits( power )

        frame = signal.readframes( 1 )
        frameNumber += 1

      signal.close()

    except IOError:
      print 'File not found:', filename

    except Exception, inst:
      print inst

