/** webtorrent-shim.js */

const parse_torrent = require("parse-torrent");
const WebTorrent = require("webtorrent");

function Shim()
{
  this.torrent = new WebTorrent();
  this._lastInfohash = null;
}

Shim.prototype.Start = function()
{
};

Shim.prototype.DownloadRate = function()
{
  var self = this;
  return self.torrent.downloadSpeed;
}

Shim.prototype.UploadRate = function()
{
  var self = this;
  return self.torrent.uploadSpeed;
}

Shim.prototype.PeerCount = function()
{
  var self = this;
  var numPeers = 0;
  for(var idx = 0; idx < self.torrent.torrents.length; idx++)
  {
    var peers = self.torrent.torrents[idx].numPeers;
    if(numPeers < peers) numPeers = peers;
  }
  return numPeers;
}

Shim.prototype.Cleanup = function(rewind)
{
  var self = this;
  var torrents = self.torrent.torrents;
  while(torrents.length > rewind)
  {
    var ih = torrents[0].infoHash;
    if(ih == self._lastInfohash) continue;
    torrents.pop();
    if(ih)
      self.torrent.remove(ih, function(err) { });
  }
};


Shim.prototype.FetchMetadata = function(url, cb)
{
  var self = this;
  parse_torrent.remote(url, function(err, tfile) {
    if(err)
      cb(err, null);
    else
      cb(null, tfile);
  });
};

Shim.prototype.Stream = function(metadata, elem, first)
{
  var self = this;
  if(!first)
  {
    if(self.torrent.get(metadata.infoHash)) return;
    if(self._lastInfohash == metadata.infoHash) return;
  }
  self._lastInfohash = metadata.infoHash;
  self.torrent.add(metadata, function(t) {
    t.files[0].renderTo(elem);
  });
}

Shim.prototype.AddMetadata = function(metadata, cb)
{
  var self = this;
  if (self.torrent.get(metadata.infoHash)) return;
  if(self._lastInfohash == metadata.infoHash) return;
  self._lastInfohash = metadata.infoHash;
  self.torrent.add(metadata, function(t) {
    t.files[0].getBlob(function(err, blob) {
      if(err) cb(err, null);
      else cb(null, blob);
    });
  });
};

Shim.prototype.SeedData = function(data, name, cb)
{
  var self = this;
  data.name = name;
  self.torrent.seed(data, function(t) {
    cb(null, t.torrentFile);
  });
};

module.exports = {
  "Network": Shim
};
