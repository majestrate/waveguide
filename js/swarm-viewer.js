/** swarm-viewer.js */

const util = require("./util.js");

function Viewer(elem)
{
  var e;
  this.root = util.div("swarm");
  this.rxLabel = util.span(null, "rx");
  this.txLabel = util.span(null, "tx");
  this.peersLabel = util.span(null, "peers");
  this.ratioLabel = util.span(null, "ratio");

  e = util.div();
  e.appendChild(this.ratioLabel);
  e.appendChild(this.peersLabel);
  this.root.appendChild(e);

  e = util.div();
  e.appendChild(this.rxLabel);
  e.appendChild(this.txLabel);
  this.root.appendChild(e);
  
  elem.appendChild(this.root);
}

Viewer.prototype.Update = function(torrent)
{
  var self = this;
  util.set_text(self.rxLabel, "RX: " + util.format_rate(torrent.downloadSpeed));
  util.set_text(self.txLabel, "TX: " + util.format_rate(torrent.uploadSpeed));
  util.set_text(self.peersLabel, "Peers: " + torrent.numPeers);
  util.set_text(self.ratioLabel, "Ratio: " + util.format_float(torrent.ratio, 100));
};


module.exports = {
  "Viewer": Viewer
};
