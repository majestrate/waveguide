/** player.js */

const WebTorrent = require("webtorrent");
const util = require("./util.js");

function Player(wg, elem)
{
  this.root = elem;
  this.elem = document.createElement("video");
  this.elem.setAttribute("id", "video-player");
  this.elem.setAttribute("type", "video/mp4");
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

    t.files[0].getBlobURL(function(err, url) {
      if(err) self.Error(err);
      else
      {
        self.elem.src = url;
      }
    });
    
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
