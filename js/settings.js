/** settings.js */


const SegmentLength = 30 * 1000;
const RefreshInterval = SegmentLength / 8;
//const SegMime = 'video/mp4; codecs="avc1.42E01E, mp4a.40.2"';
const SegMime = 'video/webm';
const SegExt = ".webm";
const SegPlaceholder = "https://"+location.host+"/static/loading.webm";
const SegOffline = "https://"+location.host+"/static/offline.webm";
const SegBitrate = 100000;

module.exports = {
  SegLen: SegmentLength,
  RefreshInterval: RefreshInterval,
  SegMime: SegMime,
  SegExt: SegExt,
  SegPlaceholder: SegPlaceholder,
  SegOffline: SegOffline,
  SegBitrate: SegBitrate
};
