/** waveguide.js */

const util = require("./util.js");
const upload = require("./upload.js");
const player = require("./player.js");
const stream = require("./stream.js");
const cam = require("./cam.js");
const chat = require("./chat.js");

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

WaveGuide.prototype.ChatWidget = function(elemId, streamID)
{
  var self = this;
  var elem = util.get_id(elemId);
  if(!elem) return null;
  return new chat.LiveChat(elem, streamID);
};

WaveGuide.prototype.Streamer = function(pubkey, experimental)
{
  var self = this;
  if(self._stream) return;
  if(pubkey)
  {
    console.log("streaming "+pubkey);
    if (experimental)
      self._stream = new stream.Streamer(null, pubkey, true);
    else
      self._stream = new stream.Streamer(null, pubkey);
    self._stream.Start();
  }
  else
  {
    console.log("asking for cam");
    cam.getCam(function(err, src) {
      if(err) throw err;
      if(src)
      {
        console.log("Got stream");
        self._stream = new stream.Streamer(src)
        self._stream.Start();
      }
    });
  }
};

module.exports = {
  "WaveGuide": WaveGuide
};
