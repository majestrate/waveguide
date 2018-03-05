/** stream.js */

const Segmenter = require("./segment.js").Segmenter;
const util = require("./util.js");
const settings = require("./settings.js");
const shim = require("./webtorrent-shim.js");

function Streamer(source, key)
{
  this._key = key || null;
  this._source = source;
  this._rewind = 2;
  this._interval = null;
  this._segments = null;
  this._video = null;
  this._logelem = util.get_id("log");
  if(source)
    util.get_id("cam").src = window.URL.createObjectURL(source);
  else
    this._segments = [];
  this._net = new shim.Network();
  this._lastSegmentURL = null;
  this._segmentCounter = 0;
  this._lastPopped = 0;
}

Streamer.prototype.Start = function()
{
  var self = this;
  self._net.Start();
  self._onStarted();
};

Streamer.prototype.log = function(msg)
{
  console.log(msg);
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
  self.log("cleanup storage");
  var exclude = self.SegmentInfohashes()
  if (exclude.length > self._rewind)
  {
    var diff = exclude.length - self._rewind;
    exclude = exclude.slice(exclude.length-self._rewind);
    while(diff > 0)
    {
      var seg = self._segments.shift();
      diff--;
      self.log("discard segment "+seg[1]);
    }
  }
  self._net.Cleanup(exclude);
};

Streamer.prototype.Stop = function()
{
  var self = this;
  if(self._interval) clearInterval(self._interval);
  if(self._segmenter) self._segmenter.Stop();
};

Streamer.prototype._hasSegment = function(idx)
{
  var self = this;
  for(var i = 0; i < self._segments.length; i ++)
  {
    if (self._segments[i][1] == idx) return true;
  }
  return false;
}

Streamer.prototype.SegmentInfohashes = function()
{
  var self = this;
  var list = [];
  for (var i = 0; i < self._segments.length; i ++)
  {
    list.push(self._segments[i][2]);
  }
  return list;
}

Streamer.prototype._queueSegment = function(f, idx, ih)
{
  var self = this;
  if(f)
  {
    if(!self._hasSegment(idx))
    {
      self._segments.push([f, idx, ih]);
      self._segmentCounter ++;
      self._segments = self._segments.sort(function(a, b) {
        return a[1] - b[1];
      });
    }
  }
};

Streamer.prototype._popSegmentBlob = function()
{
  var self = this;
  console.log("segments:" + self._segments);
  var seg = self._segments.shift();
  if(seg && seg[0])
  {
    self.log("pop segment "+seg[1]);
    return URL.createObjectURL(seg[0]);
  }
  else
    return null;
};


Streamer.prototype._nextSegment = function(url)
{
  var self = this;
  if (self._lastSegmentURL != url) {
    self._net.FetchMetadata(url, function(err, metadata) {
      if(err)
      {
        self.log("no stream online: "+err);
        if(!(self._video.src === settings.SegOffline))
        {
          self._video.src = settings.SegOffline;
          self._video.loop = true;
          self._playVideo();
        }
      }
      else
      {
        self.log("segment "+self._segmentCounter);

        var curSeg = self._segmentCounter;
        self._net.AddMetadata(metadata, function(err, blob) {
          if (err) self.log("failed to fetch file: "+err);
          else
          {
            if(self._segmentCounter > curSeg)
            {
              self.log("out of order segment");
            }
            self._queueSegment(blob, curSeg, metadata.infoHash);
            self._lastSegmentURL = url
          }
        });
      }
    });
  }
  if (self._video.src === settings.SegPlaceholder || self._video.src === settings.SegOffline)
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

Streamer.prototype._segmenterCB = function(data, name)
{
  var self = this;
  self._net.SeedData(data, name, function(err, metadata) {
    if(err) self.log(err);
    else
    {
      var ajax = new XMLHttpRequest();
      ajax.onreadystatechange = function() {
        if (ajax.readyState == 4 && ajax.status != 200) {
          self.Stop();
        }
      };
      ajax.open("POST", "/wg-api/v1/authed/stream-update");
      ajax.send(metadata);
      self.Cleanup();
    }
  });
};

Streamer.prototype._getNextSegment = function()
{
  var self = this;
  var ajax = new XMLHttpRequest();
  ajax.open("GET", "/wg-api/v1/stream?u="+self._key);
  ajax.onreadystatechange = function() {
    if (ajax.readyState == 4 && ajax.status == 200) {
      var url = ajax.responseText;
      self._nextSegment(url);
    }
  }
  ajax.send();
};

Streamer.prototype._onStarted = function()
{
  var self = this;
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
    var upload = self._net.UploadRate();
    var download = self._net.DownloadRate();
    var peers = self._net.PeerCount();
    self.BWLabel(upload, download);
    self.PeersLabel(peers);
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
        var oldSrc = self._video.src;
        self._video.src = blob;
        URL.revokeObjectURL(oldSrc);
        self._playVideo();
      }
      else
      {
        self.log("reached the end of segment but no next segment yet");
        self._getNextSegment();
        setTimeout(function() {
          next();
        }, 1000);
      }
    };
    self._video.onended = next;
    self._getNextSegment();
    self._interval = setInterval(function() {
      self._getNextSegment();
    }, settings.RefreshInterval);
  }
};

module.exports = {
  "Streamer": Streamer
};
