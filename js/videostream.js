/** videostream.js */

const Readable = require('readable-stream').Readable;
const util = require('util');

function Stream() {
  Readable.call(this);
  this.buffers = [];
  this.name = "stream.mp4";
};

util.inherits(Stream, Readable);

Stream.prototype._read = function(size)
{
  
};

Stream.prototype.createReadStream = function(opts)
{
  var self = this;
  console.log(opts);
  return self;
};

Stream.prototype.Put = function(blob, idx)
{
  var self = this;
  self.buffers.push([blob, idx]);
  self.buffers.sort(function(a, b) {
    return a[1] - b[1];
  });
};

module.exports = {
  "Stream": Stream
};
