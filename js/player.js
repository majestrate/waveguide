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
      t.files[0].renderTo(self.elem, {}, function(err, elem) {
        if(err) throw err;
      });
    });
    if(cb) cb();
  });
}

Player.prototype.Start = function()
{
  var self = this;
  if(self.torrent)
  {
    self.torrent.start();
    self.torrent.play();
  }
}

Player.prototype.Stop = function()
{
  var self = this;
  if(self.torrent)
  {
    self.elem.pause();
    self.torrent.stop();
  }
}


module.exports = {
  "TorrentPlayer" : Player
};
