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
  ev.data.name = "segment" + settings.SegExt;
  cb(ev.data);
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
      self.MakeData(ev, cb);
      self._collector.stop();
    }
  };
  self._collector.start(settings.SegLen);
};

module.exports = {
  Segmenter: Segmenter
};
