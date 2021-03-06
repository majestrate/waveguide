/* upload.js */

const util = require("./util.js");

const state_Submit = "submit";
const state_Error = "error";
const state_Ready = "ready";

function UploadWidget(parent, elem)
{
  this.waveguide = parent;
  this.elem = elem;
  this.submit = null;
  this.file = null;
  this.form = null;
  this.message = null;
  this.webseed_url = null;
  this._state = state_Ready;
}

UploadWidget.prototype.IsSumbitting = function()
{
  return this._state === state_Submit;
};

UploadWidget.prototype._EnterState = function(state)
{
  this._state = state;
};

UploadWidget.prototype.Setup = function()
{
  var self = this;
  self.form = util.form();

  var button_id = "upload_button";
  self.submit = util.button(button_id, "upload-button");
  self.submit.value = "upload";
  self.submit.onclick = function() { self.Submit(); };

  var e = util.div(null, "upload");
  e.appendChild(self.submit);
  self.form.appendChild(e);
  
  var file_id = "upload_file";
  self.file = util.file(file_id, "upload-file");
  var e = util.div(null, "upload");
  e.appendChild(self.file);
  self.form.appendChild(e);

  var message_id = "upload_message";
  self.message = util.div(message_id, "upload-label");
  self.form.appendChild(self.message);

  var webseed_id = "upload_webseed";
  self.webseed_url = util.input(webseed_id, "upload-label");

  e = util.div(null, "upload");
  e.appendChild(self.webseed_url);
  e.appendChild(document.createTextNode("webseed url (optional)"));
  self.form.appendChild(e);
  
  self.elem.appendChild(self.form);
  
};

UploadWidget.prototype.Error = function(msg)
{
  var self = this;
  self._EnterState(state_Error);
  self.SetMessage(msg, "error");
};

UploadWidget.prototype.Clear = function()
{
  var self = this;
  self.file.value = "";
  self.SetMessage("", "");
};

UploadWidget.prototype.Success = function(url)
{
  var self = this;
  self._EnterState(state_Ready);
  self.SetMessage("upload okay", "success");
  setTimeout(function() {
    window.location = url;
  }, 1000);
};

UploadWidget.prototype.SetMessage = function(msg, cls)
{
  var self = this;
  while(self.message.children.length > 0)
  {
    self.message.firstChild.remove();
  }
  var e = util.div();
  e.appendChild(document.createTextNode(msg));
  self.message.appendChild(e);
  self.message.setAttribute("class", cls);
};

UploadWidget.prototype.Progress = function(progress)
{
  var self = this;
  self.SetMessage("uploading... "+(parseInt(progress*1000)/10)+"%", "upload");
};

UploadWidget.prototype.Submit = function()
{
  var self = this;
  if (self.IsSumbitting()) return;

  var ajax = new XMLHttpRequest();

  self._EnterState(state_Submit);

  var formdata = new FormData();
  if(self.file.files.length == 1)
  {
    formdata.append("video", self.file.files[0]);
  }
  else if(self.webseed_url.value.length > 0)
  {
    var parts = self.webseed_url.value.split("/");
    console.log(parts);
    formdata.append("webseed", self.webseed_url.value);
    formdata.append("title", parts[parts.length-1]);
  }
  else 
  {
    self.Error("no video provided");
    return;
  }
  
  ajax.onreadystatechange = function()
  {
    if(ajax.readyState == 4)
    {
      var err = null;
      var j = null;
      if(ajax.status == 201 || ajax.status == 200)
      {
        try
        {
          j = JSON.parse(ajax.responseText);
        }
        catch ( ex )
        {
          err = ex;
        }
        err = err || j.error;
        if(err) self.Error(err);
        else if(ajax.status == 201) self.Success(j.url);
        else self.Error("failed to create video");
      }
      else
      {
        self.Error("http "+ajax.status);
      }
    }
  };
  ajax.upload.addEventListener("progress", function(ev) {
    if(ev.lengthComputable)
    {
      self.Progress(ev.loaded / ev.total);
    }
  });
  ajax.open("POST", "/wg-api/v1/authed/upload");
  ajax.send(formdata);
};

module.exports = {
  "Widget": UploadWidget
}
