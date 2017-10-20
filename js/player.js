/** player.js */

const WebTorrent = require("webtorrent");


function Player(wg, elem)
{
  this.elem = elem;
  this.client = new WebTorrent();
  this.torrent = null;
}

Player.prototype.SetInfo = function(url, webseeds, cb)
{
  var self = this;
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
