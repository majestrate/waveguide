/** segment.js */

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
  try {
    self._collector = new MediaRecorder(self._source, {mimeType: "video/webm"});
  } catch ( ex ) {
    console.error("failed to begin capture: " + ex);
  }
  if(self._collector)
  {
    console.log("starting...");
    self._collector.ondataavailable = function(ev) {
      console.log("got chunk of size "+ev.data.size);
      ev.data.name = "segment.webm";
      cb(ev.data);
    };
    self._collector.start(1000 * 30);
  }
};

module.exports = {
  Segmenter: Segmenter
};
