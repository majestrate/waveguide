/** segment.js */

var Buffer = require("buffer").Buffer;

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

Segmenter.prototype.Begin = function(cb)
{
  var self = this;
  self._collector = new MediaRecorder(self._source, {mimeType: "video/webm"});
  console.log("starting...");
  var c = function(ev) {
    console.log("got chunk of size "+ev.data.size);
    ev.data.name = "segment.webm";
    self._collector.stop();
    cb(ev.data);
  };
  self._collector.ondataavailable = c;
  self._collector.start(1000 * 30);
};

module.exports = {
  Segmenter: Segmenter
};
