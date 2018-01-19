/** ipfs-shim.js */

const IPFS = require("ipfs");
const REPO = "waveguide_ipfs_repo";
const toBuffer = require("blob-to-buffer");

const settings = require('./settings.js');

function Shim()
{
  this._segs = {};
  this._node = new IPFS({
    repo: REPO,
    init: true,
    config: {
      Addresses: {
        Swarm: ['/dns4/wrtc-star.discovery.libp2p.io/tcp/443/wss/p2p-webrtc-star']
      }
    },
    EXPERIMENTAL: { 
      pubsub: true,
      sharding: true, 
      dht: true 
    }
  });
  this._started = false;
}

Shim.prototype.Start = function()
{
  var self = this;
  self._node.on("init", function() {
    self._ready();
  });
  self._node.on("start", function() {
    self._started = true;
  });
};

Shim.prototype._ready = function()
{
  var self = this;
  self._node.start();
};

Shim.prototype._do = function(cb)
{
  var self = this;
  if (self._started) cb();
  else setTimeout(function() { cb(); }, 1000);
};

Shim.prototype.Cleanup = function(rewind)
{
 /** TODO */
};

Shim.prototype.UploadRate = function()
{
  return 0;
};

Shim.prototype.DownloadRate = function()
{
  return 0;
};

Shim.prototype.PeerCount= function()
{
  return 0;
};

Shim.prototype.FetchMetadata = function(url, cb)
{
  var ajax = new XMLHttpRequest();
  ajax.onreadystatechange = function() {
    if(ajax.readyState == 4)
    {
      if(ajax.status == 200)
        cb(null, ajax.responseText);
      else
        cb("http "+ajax.status, null);
    }
  };
  ajax.open("GET", url);
  ajax.send();
};

Shim.prototype._put = function(metadata)
{
  var self = this;
  self._segs[metadata] = 1;
};

Shim.prototype._has = function(metadata)
{
  var self = this;
  return self._segs[metadata] == 1;
};


Shim.prototype.AddMetadata = function(metadata, cb)
{
  var self = this;
  if(self._has(metadata)) return;
  self._put(metadata);
  self._do(function() {
    self._node.files.get(metadata, function(err, files) {
      if(err) cb(err, null);
      else
      {
        files.forEach(function(f) {
          cb(null, new Blob( [f.content], {type: settings.SegMime}));
        });
      }
    });
  });
};

Shim.prototype.SeedData = function(d, name, cb)
{
  var self = this;
  self._do(function() {
    var buffer = toBuffer(d, function(err, data) {
      if(err) cb(err, null);
      else self._node.files.add(data, {} , function(err, res) {
        if(err) cb(err, null);
        else cb(null, res[0].hash);
      });
    });
  });
};

module.exports = {
  "Network": Shim
};
