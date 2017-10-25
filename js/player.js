/** player.js */

const WebTorrent = require("webtorrent");
const util = require("./util.js");
const swarm = require("./swarm-viewer.js");

function Player(wg, elem)
{
  this.root = elem;
  this.elem = document.createElement("video");
  this.elem.setAttribute("id", "video-player");
  this.elem.setAttribute("type", "video/mp4");
  this.elem.setAttribute("controls", "true");
  this.client = new WebTorrent();
  this.torrent = null;
  this.err = util.div("player_error", "error");
  this.root.appendChild(this.elem);
  this.root.appendChild(this.err);
  this.swarm = null;
  this._interval = null;
}

Player.prototype.SetInfo = function(url, webseeds, cb)
{
  var self = this;
  
  if(!url.length) {
    /** no torrent available yet */
    self.Error("video not yet ready try again later");
    return;
  }
  self.client.on("error", function(err) {
    self.Error(err);
  });
  self.client.add(url, function(t) {
    
    self.torrent = t;

    self.swarm = new swarm.Viewer(self.root);
    
    for (var idx = 0; idx < webseeds.length; idx ++)
      t.addWebSeed(webseeds[idx]);

    var file = t.files[0];
    file.renderTo(self.elem);

    self._interval = setInterval(function() {
      self.swarm.Update(self.torrent);
    }, 500);
    
    t.on("error", function(err) {
      self.Error(err);
    });

    t.on("done", function() {
      console.log("download is done");
    });
    
    if(cb) cb();
  });
};

Player.prototype.Error = function(msg)
{
  console.log("player error: " + msg);
  var self = this;
  util.set_text(self.err, msg);
};

Player.prototype.Start = function()
{
  var self = this;
  console.log("starting player");
};

Player.prototype.Stop = function()
{
  var self = this;
};


module.exports = {
  "TorrentPlayer" : Player
};
