/* upload.js */

const util = require("./util.js");

function UploadWidget(parent, elem)
{
  this.waveguide = parent;
  this.elem = elem;
  this.submit = null;
  this.file = null;
  this.form = null;
}

UploadWidget.prototype.Setup = function()
{
  var self = this;
  self.form = util.form();
  
  var button_id = "upload_button";
  self.submit = util.submit(button_id, "upload-button");
  self.submit.innerHTML = "upload";
  self.form.appendChild(self.submit);
  
  var file_id = "upload_file";
  self.file = util.file(file_id, "upload-file");
  self.form.appendChild(self.file);

  self.elem.appendChild(self.form);
};

module.exports = {
  "Widget": UploadWidget
}
