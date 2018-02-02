/** chat.js */

/** base api endpoint */
const ApiBaseUrl = "http://api.sapphire.moe:7070/";

/** 
    build a chat post element to be inserted 
    TODO: make it pretty
*/
function buildChatPost(chat)
{
  var e = document.createElement("div");
  var user = chat.user;
  var name = "???";
  if(user && user.name)
  {
    name = user.name;
  }
  e.appendChild(document.createTextNode(name+": "+chat.text));
  return e;
}

function Chat(root, key)
{
  this.scrollback = 10;
  this.elem = root;
  this.key = key || 1;
  this.last_id = 0;
  this.update_interval = 500;
  this.token = null;
  this.interval = null;
  this.view = document.createElement("div");
  this.input = document.createElement("input");
  this.sendButton = document.createElement("button");
  this.sendButton.innerText = "send";
  this.elem.appendChild(this.view);
  this.elem.appendChild(this.input);
  this.elem.appendChild(this.sendButton);
}

Chat.prototype.Error = function(msg)
{
  var self = this;
  var e = document.createElement("div");
  e.setAttribute("class", "error");
  e.appendChild(document.createTextNode(msg));
  self.PutLineElem(e);
};

Chat.prototype.PutLineElem = function(e)
{
  var self = this;
  self.view.appendChild(e);
  while(self.view.children.length > self.scrollback)
  {
    self.view.removeChild(self.view.firstChild);
  }
};

Chat.prototype._makeurl = function()
{
  var self = this;
  var url = ApiBaseUrl+"channels/"+self.key+"/messages?";
  if(self.token)
  {
    url += "&access_token="+self.token;
  }
  return url;
}

Chat.prototype.GetLatest = function(cb)
{
  var self = this;
  var ajax = new XMLHttpRequest();
  ajax.onreadystatechange = function()
  {
    if(ajax.readyState == 4)
    {
      if (ajax.status == 200)
      {
        cb(JSON.parse(ajax.responseText));
      }
      else
      {
        self.Error("http "+ajax.status);
      }
    }
  };
  var url = self._makeurl();
  if(self.last_id)
  {
    url += "&last_since="+self.last_id;
  }
  ajax.open("GET", url);
  ajax.send();
};

Chat.prototype.SendChat = function(text)
{
  var self = this;
  if(!self.token) {
    self.Error("not logged in");
    return;
  }
  var ajax = new XMLHttpRequest();
  ajax.onreadystatechange = function()
  {
    if(ajax.readyState == 4)
    {
      if (ajax.status == 200)
      {
        console.log(ajax.responseText);
      }
      else
      {
        self.Error("http "+ajax.status);
      }
    }
  };
  var data = JSON.stringify({text: text});
  var url = self._makeurl();
  ajax.open("POST", url);
  ajax.setRequestHeader("Content-Type", "application/json");
  ajax.send(data);
};

Chat.prototype.SetToken = function(token)
{
  var self = this;
  token = token.trim();
  if(token.length > 0)
  {
    self.token = token;
  }
}

/** start chat updater */
Chat.prototype.Start = function()
{
  var self = this;
  self.sendButton.onclick = function() {
    self.SendChat(self.input.value);
    self.input.value = "";
  };
  if(self.interval) return;
  self.interval = setInterval(function() {
    self.GetLatest(function(j) {
      var data = j.data || [];
      for(var idx = data.length-1; idx > 0; idx--)
      {
        self.PutLineElem(buildChatPost(data[idx]));
      }
      /* update last id */
      if(data.length > 0)
      {
        self.last_id = data[0].id || self.last_id;
      }
    });
  }, self.update_interval);
};


Chat.prototype.Stop = function()
{
  var self = this;
  if(self.interval)
  {
    clearInterval(self.interval);
    self.interval = null;
  }
};

module.exports = {
  "LiveChat": Chat
};
