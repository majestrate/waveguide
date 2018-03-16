/** videostream.js */

const settings = require("./settings.js");
const toArrayBuffer = require('to-arraybuffer');
//const mp4box = require("mp4box");


function Player(videoElem)
{
  this._video = videoElem;
  if(!window.MediaSource)
  {
    throw new Error("MediaSource not supported");
  }
  this._mediasource = new MediaSource();
  this._sourceBuffer = null;
  this._buffers = [];
  this._stalled = false;
  this._offline = true;
  this._mp4 = new mp4box();
};

Player.prototype._openHandler = function()
{
  var self = this;
  console.log("open buffer");
  self._sourceBuffer = self._mediasource.addSourceBuffer(settings.SegMime);
  self._sourceBuffer.addEventListener('updateend', self._flowHandler);
  self._flowHandler();
}

Player.prototype._flowHandler = function()
{
  console.log('flow');
  var self = this;
  if(self._sourceBuffer)
  {
    var chunk = self._buffers.shift();
    if (chunk)
    {
      console.log("append buffer "+chunk.length);
      self._sourceBuffer.appendBuffer(chunk);
    }
    else
      console.log("no segment");
  }
  else
    console.log("no source yet");
}

/** start processing of video fragments */
Player.prototype.Start = function()
{
  var self = this;
  self._offline = false;
  self._mp4.onReady = function(info)
  {
    self._mp4.onSegment = function(id, user, buffer)
    {
    };
  };
  if(MediaSource.isTypeSupported(settings.SegMime))
  {
    if(!self._sourceBuffer)
    {
      console.log("media supported, opening");
      self._mediasource.addEventListener('sourceopen', function() {
        console.log("opened source");
        self._openHandler();
      });
      self._mediasource.addEventListener('timeupdate', function() {
        console.log("time update");
        self._flowHandler();
      });
      self._video.src = window.URL.createObjectURL(self._mediasource);
      self._flowHandler();
    }
  }
  else
    throw new Error("browser does not support required video format");
};

/** play underlying video element */
Player.prototype.Play = function()
{
  var self = this;
  self._video.play();
};

/** set player as stalled */
Player.prototype.Stalled = function()
{
  var self = this;
  return self._stalled;
};

Player.prototype.Started = function()
{
  var self = this;
  return self._sourceBuffer != null;
}

/** add a segment to the video segment queue */
Player.prototype.AddSegment = function(buff)
{
  var self = this;
  console.log("new segment");
  self._buffers.push(buff);
};

/** return true if we think the stream is offline */
Player.prototype.IsOffline = function()
{
  var self = this;
  return self._offline;
};

/** set the stream as explicitly offline */
Player.prototype.SetOffline = function()
{
  var self = this;
  self._offline = true;
};

Player.prototype.SetOnline = function()
{
  var self = this;
  self._offline = false;
}

module.exports = {
  "Player": Player
};
