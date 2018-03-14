/** swarm-viewer.js */

const util = require("./util.js");

function Viewer(elem)
{
  this.root = document.getElementById("swarm");
  this.rxLabel = document.getElementById("rx");
  this.txLabel = document.getElementById("tx");
  this.peersLabel = document.getElementById("peers");
  this.ratioLabel = document.getElementById("ratio");
  /*
  e = util.div();
  e.appendChild(this.ratioLabel);
  e.appendChild(this.peersLabel);
  this.root.appendChild(e);

  e = util.div();
  e.appendChild(this.rxLabel);
  e.appendChild(this.txLabel);
  this.root.appendChild(e);
  
  elem.appendChild(this.root);
  */
}

Viewer.prototype.Update = function(torrent)
{
  var self = this;
  util.set_text(self.rxLabel, " " + util.format_rate(torrent.downloadSpeed));
  util.set_text(self.txLabel, " " + util.format_rate(torrent.uploadSpeed));
  util.set_text(self.peersLabel, " " + torrent.numPeers);
  util.set_text(self.ratioLabel, " " + util.format_float(torrent.ratio, 100));
};


module.exports = {
  "Viewer": Viewer
};
