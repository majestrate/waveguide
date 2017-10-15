/** waveguide.js */

const util = require("./util.js");
const upload = require("./upload.js");
const player = require("./player.js");

function WaveGuide()
{

}

WaveGuide.prototype.UploadWidget = function(id)
{
  var elem = util.get_id(id);
  if(!elem) return null;
  var self = this;
  return new upload.Widget(self, elem);
};


WaveGuide.prototype.VideoPlayer = function(id)
{
  var elem = util.get_id(id);
  if(!elem) return null;
  var self = this;
  return new player.TorrentPlayer(self, elem);
}

module.exports = {
  "WaveGuide": WaveGuide
};
