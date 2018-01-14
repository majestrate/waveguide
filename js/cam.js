/** cam.js */


var grabcam = function(cb) {
  var vidopts = {
    manditory: {
      chromeMediaSource: 'screen'
    }
  };
  navigator.mediaDevices.getUserMedia({video: vidopts, audio:true}).then(function(src) {
    cb(null, src);
  }).catch(function(e) { cb(e, null); });
};


module.exports = {
  "getCam": grabcam
}
