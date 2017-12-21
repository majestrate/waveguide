/** stream.js */

const WebTorrent = require("webtorrent");
const Segmenter = require("./segment.js").Segmenter;
const util = require("./util.js");
const parse_torrent = require('parse-torrent');

function Streamer(source, key)
{
  this._key = key || null;
  this._source = source;
  this._rewind = 5;
  this._interval = null;
  this._segments = null;
  this._video = null;
  if(source)
    util.get_id("cam").src = window.URL.createObjectURL(source);
  else
    this._segments = [];
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

Streamer.prototype._queueSegment = function(f)
{
  var self = this;
  self._segments.push(f);
};

Streamer.prototype._popSegmentBlob = function(cb)
{
  var self = this;
  var seg = self._segments.pop();
  if(seg)
  {
    seg.getBlobURL(function(err, url) {
      cb(err, url);
    });
  }
  else
    cb(null, null);
};

Streamer.prototype._nextSegment = function(url)
{
  var self = this;
  parse_torrent.remote(url, function(err, tfile) {
    if(err) console.log(err);
    else if (!self.torrent.get(tfile.infoHash))
    {
      console.log("add torrent "+tfile.infoHash);
      self.torrent.add(parse_torrent.toTorrentFile(tfile), function(t) {
        self._queueSegment(t.files[0]);
        console.log(self._video.src);
        if (!self._video.src)
        {
          console.log("pop segment");
          self._popSegmentBlob(function(err, blob) {
            if(err) console.error(err);
            else if(blob)
            {
              self._video.src = blob;
              self._video.play();
            }
            else
            {
              console.log("no blob");
            }
          });
        }   
      });
    }
  });
  self.Cleanup();
};

Streamer.prototype._onStarted = function()
{
  var self = this;
  
  setInterval(function() {
    console.log(self.torrent.uploadSpeed, self.torrent.downloadSpeed);
  }, 1000);
  if (self._key)
  {    
    self._video = util.get_id("player");
    self._video.onended = function() {
      console.log("pop next segment");
      self._popSegmentBlob(function(err, url) {
        if(err) console.error(err);
        else
        {
          self._video.src = url;
          self._video.play();
        }
      });
    };
    var url = location.protocol+"//"+location.host+"/wg-api/v1/stream/"+self._key;
    self._interval = setInterval(function() {
      self._nextSegment(url);
    }, 10000);
    self._nextSegment(url);
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
        ajax.send(t.torrentFile);
        self.Cleanup();
      });
    });
  }
};

module.exports = {
  "Streamer": Streamer
};
