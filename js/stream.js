/** stream.js */

const WebTorrent = require("webtorrent");
const Segmenter = require("./segment.js").Segmenter;
const util = require("./util.js");

function Streamer(source, key)
{
  this._key = key || null;
  this._source = source;
  this._rewind = 5;
  this._interval = null;
  if(source)
    util.get_id("cam").src = window.URL.createObjectURL(source);
}

Streamer.prototype.Start = function()
{
  var self = this;
  self.torrent = new WebTorrent();
  self._onStarted();
};

Streamer.prototype.Cleanup = function()
{
  var self = this;
  console.log("cleanup torrents");
  var torrents = self.torrent.torrents;
  while(torrents.length > self._rewind)
  {
    var ih = torrents[0].infoHash;
    torrents.pop();
    if(ih) {
      self.torrent.remove(ih, function(err) {
        if(err) console.error("failed to remove torrent: "+err);
        else console.log("removed torrent: "+ih);
      });
    }
  }
};

Streamer.prototype.Stop = function()
{
  var self = this;
  if(self._interval) clearInterval(self._interval);
  if(self._segmenter) self._segmenter.Stop();
};

Streamer.prototype._onStarted = function()
{
  var self = this;
  if (self._key)
  {
    self._interval = setInterval(function() {
      var ajax = new XMLHttpRequest();
      ajax.onreadystatechange = function() {
        if (ajax.readyState == 4) {
          if(ajax.status != 200)
          {
            console.log("no such stream");
            self.Stop();
            return;
          }
          var j = JSON.parse(ajax.responseText);
          
          var magnet = j[j.length-1];
          if(self.torrent.get(magnet)) return;
          console.log("getting next magnet: "+magnet);
          var player = util.get_id("player");
          while(player.children.length > 0)
            player.removeChild(player.firstChild);
          self.torrent.add(magnet, function(t) {
            t.files[0].appendTo("#player");
          });
        }
      };
      ajax.open("GET", "/wg-api/v1/stream/"+self._key);
      ajax.send();
      self.Cleanup();
    }, 10000);
  }
  else
  {
    self._segmenter = new Segmenter(self._source);
    self._segmenter.Begin(function(data) {
      self.torrent.seed(data, function(t) {
        console.log("submit magnet: "+ t.magnetURI);
        var ajax = new XMLHttpRequest();
        ajax.onreadystatechange = function() {
          if (ajax.readyState == 4 && ajax.status != 200) {
            
            self.Stop();
          }
        };
        ajax.open("POST", "/wg-api/v1/authed/stream-update");
        ajax.send(t.magnetURI);
        self.Cleanup();
      });
    });
  }
};

module.exports = {
  "Streamer": Streamer
};
