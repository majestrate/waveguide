/** player.js */

const WebTorrent = require("webtorrent");
const util = require("./util.js");

function Player(wg, elem)
{
  this.root = elem;
  this.elem = document.createElement("video");
  this.client = new WebTorrent();
  this.torrent = null;
  this.err = util.div("player_error", "error");
  this.root.appendChild(this.elem);
  this.root.appendChild(this.err);
}

Player.prototype.SetInfo = function(url, webseeds, cb)
{
  var self = this;
  
  if(!url.length) {
    /** no torrent available yet */
    self.Error("video not yet ready try again later");
    return;
  }
  self.elem.setAttribute("controls", "controls");
  self.client.on("error", function(err) {
    self.Error(err);
  });
  self.client.add(url, function(t) {
    
    self.torrent = t;
    
    for (var idx = 0; idx < webseeds.length; idx ++)
      t.addWebSeed(webseeds[idx]);
    
    t.on("ready", function() {
      console.log("torrent ready");
      t.files[0].renderTo(self.elem);
    });
    t.on("error", function(err) {
      self.Error(err);
    });

    t.on("done", function() {
      console.log("download is done");
      t.files[0].renderTo(self.elem);
    });
    
    if(cb) cb();
  });
};

Player.prototype.Error = function(msg)
{
  console.log("player error: " + msg);
  var self = this;
  var e = util.div();
  e.appendChild(document.createTextNode(msg));
  while(self.err.children.length > 0)
    self.err.firstChild.remove();
  self.err.appendChild(e);
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
