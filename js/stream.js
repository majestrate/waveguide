/** stream.js */

const WebTorrent = require("webtorrent");
const Segmenter = require("./segment.js").Segmenter;
const util = require("./util.js");
const parse_torrent = require('parse-torrent');
const settings = require("./settings.js");

function Streamer(source, key)
{
  this._key = key || null;
  this._source = source;
  this._rewind = 3;
  this._interval = null;
  this._segments = null;
  this._video = null;
  this._logelem = util.get_id("log");
  this.torrent = new WebTorrent();
  if(source)
    util.get_id("cam").src = window.URL.createObjectURL(source);
  else
    this._segments = [];
}

Streamer.prototype.Start = function()
{
  var self = this;
  self._onStarted();
};

Streamer.prototype.log = function(msg)
{
  var self = this;
  if(self._logelem)
  {
    var e = util.span();
    e.appendChild(document.createTextNode(msg));
    self._logelem.appendChild(e);
    while(self._logelem.children.length > 10)
    {
      self._logelem.removeChild(self._logelem.firstChild);
    }
  }
};

Streamer.prototype.Cleanup = function()
{
  var self = this;
  self.log("cleanup torrents");
  var torrents = self.torrent.torrents;
  while(torrents.length > self._rewind)
  {
    var ih = torrents[0].infoHash;
    torrents.pop();
    if(ih) {
      self.torrent.remove(ih, function(err) {
        if(err) self.log("failed to remove torrent: "+err);
        else self.log("removed torrent: "+ih);
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

Streamer.prototype._popSegmentBlob = function()
{
  var self = this;
  var seg = self._segments.pop(0);
  if(seg)
  {
    return seg;
  }
  else
    return null;
};

Streamer.prototype._nextSegment = function(url)
{
  var self = this;
  parse_torrent.remote(url, function(err, tfile) {
    if(err) self.log(err);
    else if (!self.torrent.get(tfile.infoHash))
    {
      self.log("add torrent "+tfile.infoHash);
      self.torrent.add(parse_torrent.toTorrentFile(tfile), function(t) {
        t.files[0].getBlobURL(function(err, url) {
          if(err) self.log(err);
          else self._queueSegment(url);
        });
        if (self._video.src === settings.SegPlaceholder)
        {
          self.log("pop segment");
          var blob = self._popSegmentBlob();
          self.log("got segment: "+blob);
          if(blob)
          {
            self._video.loop = false;
            self._video.src = blob;
            self._video.play();
          }
        }   
      });
    }
  });
  self.Cleanup();
};

Streamer.prototype.BWLabel = function(upload, download)
{
  var e = util.get_id("tx");
  if(e)
    e.innerHTML = "Upload: " + util.format_rate(upload);
  e = util.get_id("rx");
  if(e)
    e.innerHTML = "Download: " + util.format_rate(download);
};

Streamer.prototype.PeersLabel = function(peers)
{
  var e = util.get_id("peers");
  if(e)
    e.innerHTML = "Viewers: "+peers;
};

Streamer.prototype._segmenterCB = function(torrent, data)
{
  var self = this;
  torrent.seed(data, function(t) {
    self.log("submit magnet: "+ t.magnetURI);
    var ajax = new XMLHttpRequest();
    ajax.onreadystatechange = function() {
      if (ajax.readyState == 4 && ajax.status != 200) {
        self.Stop();
      }
    };
    ajax.open("POST", "/wg-api/v1/authed/stream-update");
    ajax.send(t.torrentFile);
    self.Cleanup();
    self._segmenter.Begin(function(data) {
      self._segmenterCB(torrent, data);
    });
  });
};

Streamer.prototype._onStarted = function()
{
  var self = this;
  
  setInterval(function() {
    self.BWLabel(self.torrent.uploadSpeed, self.torrent.downloadSpeed);
    self.PeersLabel(self.torrent.numPeers || 0);
  }, 1000);
  if (self._key)
  {    
    self._video = util.get_id("player");
    self._video.src = settings.SegPlaceholder;
    self._video.loop = true;
    self._video.play();
    var next = function() {
      var blob = self._popSegmentBlob();
      if(blob)
      {
        self.log("popped next segment");
        self._video.loop = false;
        self._video.src = blob;
        self._video.play();
        self._video.onended = next;
      }
      else
        setTimeout(function() {
          next();
        }, 1000);
    };
    self._video.onended = next;
    var url = "https://"+location.host+"/wg-api/v1/stream/"+self._key;
    self._interval = setInterval(function() {
      self._nextSegment(url);
    }, 10000);
    self._nextSegment(url);
  }
  else
  {
    self._segmenter = new Segmenter(self._source);
    self._segmenter.Begin(function(data) {
      self._segmenterCB(self.torrent, data);
    });
  }
};

module.exports = {
  "Streamer": Streamer
};
