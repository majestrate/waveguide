/** waveguide.js */

const util = require("./util.js");
const upload = require("./upload.js");
const player = require("./player.js");
const stream = require("./stream.js");

function WaveGuide()
{
  this._stream = null;
}

WaveGuide.prototype.UploadWidget = function(id)
{
  var elem = util.get_id(id);
  if(!elem) return null;
  var self = this;
  return new upload.Widget(self, elem);
};


WaveGuide.prototype.VideoPlayer = function(id)
{
  var elem = util.get_id(id);
  if(!elem) return null;
  var self = this;
  return new player.TorrentPlayer(self, elem);
};

WaveGuide.prototype.Streamer = function(pubkey)
{
  var self = this;
  if(self._stream) return;
  if(pubkey)
  {
    console.log("streaming "+pubkey);
    self._stream = new stream.Streamer(null, pubkey);
    self._stream.Start();
  }
  else
  {
    console.log("asking for cam");
    navigator.mediaDevices.getUserMedia({video: true, audio:true}).then(function(st) {
      console.log("Got stream");
      self._stream = new stream.Streamer(st)
      self._stream.Start();
    }).catch(function(e) { console.log("failed to grab cam: "+e);} );
  }
};

module.exports = {
  "WaveGuide": WaveGuide
};
