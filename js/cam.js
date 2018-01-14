/** cam.js */


var grabcam = function(cb) {
  navigator.mediaDevices.getUserMedia({video: true, audio:true}).then(function(src) {
    cb(null, src);
  }).catch(function(e) { cb(e, null); });
};


module.exports = {
  "getCam": grabcam
}
