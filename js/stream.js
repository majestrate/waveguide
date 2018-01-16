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
    msg += "\n";
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
    if(err)
    {
      self.log("no stream online");
      if(!(self._video.src === settings.SegOffline))
      {
        self._video.src = settings.SegOffline;
        self._playVideo();
      }
    }
    else if (!self.torrent.get(tfile.infoHash))
    {
      self.log("add torrent "+tfile.infoHash);
      self.torrent.add(parse_torrent.toTorrentFile(tfile), function(t) {
        t.files[0].getBlobURL(function(err, url) {
          if(err) self.log(err);
          else self._queueSegment(url);
        });
      });
    }
  });
  if (self._video.src === settings.SegPlaceholder)
  {
    var blob = self._popSegmentBlob();
    if(blob)
    {
      self.log("got segment: "+blob);
      self._video.loop = false;
      self._video.src = blob;
      self._playVideo();
    }
    else
      self.log("no segment yet");
  }
  self.Cleanup();
};



Streamer.prototype._playVideo = function()
{
  var self = this;
  if (self._video)
  {
    if (util.isAndroid())
      return;
    self._video.play();
  }
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
    e.innerHTML = "Active Peers: "+peers;
};

Streamer.prototype._segmenterCB = function(torrent, data)
{
  var self = this;
  torrent.seed(data, function(t) {
    self.log("submit torrent: "+ t.infoHash);
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

  if(!WebTorrent.WEBRTC_SUPPORT)
  {
    self.log("no webrtc support this stream will not work D:");
    return;
  }
  if (util.isAndroid())
  {
    var playbutton = util.button();
    playbutton.value = "play video";
    playbutton.onclick = function(ev) {
      self._video.play();
      self.log("play pressed");
    };
    document.body.appendChild(playbutton);
  }
  setInterval(function() {
    self.BWLabel(self.torrent.uploadSpeed, self.torrent.downloadSpeed);
    var numPeers = 0;
    for(var idx = 0 ; idx < self.torrent.torrents.length; idx ++)
    {
      var peers = self.torrent.torrents[idx].numPeers || 0;
      if(numPeers < peers) numPeers = peers;
    }
    self.PeersLabel(numPeers);
  }, 1000);
  if (self._key)
  {    
    self._video = util.get_id("player");
    self._video.src = settings.SegPlaceholder;
    self._video.loop = true;
    self._playVideo();
    var next = function() {
      var blob = self._popSegmentBlob();
      if(blob)
      {
        self.log("popped next segment");
        self._video.loop = false;
        self._video.src = blob;
        self._playVideo();
        self._video.onended = next;
      }
      else
      {
        self.log("reached the end of segment but no next segment yet");
        setTimeout(function() {
          next();
        }, 1000);
      }
    };
    self._video.onended = next;
    var url = "https://"+location.host+"/wg-api/v1/stream/"+self._key;
    self._interval = setInterval(function() {
      self._nextSegment(url);
    }, 2500);
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
