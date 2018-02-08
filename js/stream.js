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
  self._net.Cleanup(self._rewind);
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
  if(f)
  {
    self._segments.push(f);
    self._segmentCounter ++;
  }
};

Streamer.prototype._popSegmentBlob = function()
{
  var self = this;
  var seg = self._segments.pop(0);
  if(seg)
  {
    return URL.createObjectURL(seg);
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
          self._playVideo();
        }
      }
      else
      {
        self.log("segment "+self._segmentCounter);
        if(self._segmentCounter > 0)
        {
          self._net.AddMetadata(metadata, function(err, blob) {
            if (err) self.log("failed to fetch file: "+err);
            else
            {
              self._queueSegment(blob);
              self._lastSegmentURL = url
            }
          });
        }
        else
        {
          self._net.Stream(metadata, self._video, true);
          self._video.loop = false;
          self._segmentCounter +=  1;
          self._lastSegmentURL = url;
        }
      }
    });
    
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
  ajax.open("GET" , "https://"+location.host+"/wg-api/v1/stream?u="+self._key);
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
        self._video.src = blob;
        self._playVideo();
        /* self._video.onended = next; */
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
  else
  {
    self._segmenter = new Segmenter(self._source);
    try {
      self._segmenter.Begin(function(data, name) {
        self._segmenterCB(data, name);
      });
    } catch(ex) {
      self.log("error starting video recorder: "+ex);
      self._segmenter.Stop();
    }
  }
};

module.exports = {
  "Streamer": Streamer
};
