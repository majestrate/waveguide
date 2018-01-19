/** segment.js */

var Buffer = require("buffer").Buffer;
var settings = require("./settings.js");

function Segmenter(source)
{
  this._source = source;
  this._collector = null;
  this.fps = 30;
}

Segmenter.prototype.Stop = function()
{
  var self = this;
  if(self._collector)
  {
    self._collector.stop();
    self._collector = null;
  }
}

Segmenter.prototype.MakeData = function(ev, cb)
{
  var self = this;
  //console.log("got chunk of size "+ev.data.size);
  cb(ev.data, "segment" + settings.SegExt);
};

Segmenter.prototype.Begin = function(cb)
{
  var self = this;
  self._collector = new MediaRecorder(self._source, {mimeType: settings.SegMime, bitsPerSecond: settings.SegBitrate});
  //console.log("starting...");
  self.cb = cb;
  self._collector.ondataavailable = function(ev) {
    if (self._collector.state === 'recording')
    {
      console.log('make data '+ev);
      self.MakeData(ev, cb);
      self._collector.stop();
    }
    else if(self._collector.state === 'inactive')
    {
      self._collector.start(settings.SegLen);
    }
    console.log(self._collector.state);
  };
  self._collector.start(settings.SegLen);
};

module.exports = {
  Segmenter: Segmenter
};
