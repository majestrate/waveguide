/** settings.js */


const SegmentLength = 30 * 1000;
const RefreshInterval = SegmentLength / 4;
const SegMime = "video/webm;codecs=vp8";
const SegExt = ".webm";
const SegPlaceholder = "https://"+location.host+"/static/loading" + SegExt;
const SegBitrate = 20 * 1000;

module.exports = {
  SegLen: SegmentLength,
  RefreshInterval: RefreshInterval,
  SegMime: SegMime,
  SegExt: SegExt,
  SegPlaceholder: SegPlaceholder,
  SegBitrate: SegBitrate
};
