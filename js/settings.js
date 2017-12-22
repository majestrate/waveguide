/** settings.js */


const SegmentLength = 30;
const RefreshInterval = SegmentLength / 4;
const SegMime = "video/x-matroska;codecs=avc1";
const SegExt = ".m4v";
const SegPlaceholder = "/static/waiting" + SegExt;

module.exports = {
  SegLen: SegmentLength,
  RefreshInterval: RefreshInterval,
  SegMime: SegMime,
  SegExt: SegExt,
  SegPlaceholder: SegPlaceholder
};
