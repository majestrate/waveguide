/** chat.js */


function Chat(key)
{
  this.key = key;
}

Chat.prototype.Error = function(msg)
{
};

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
  ajax.open("GET", "/wg-api/v1/comments?id="+self.key);
  ajax.send();
};

Chat.prototype.SendChat = function(text)
{
  var self = this;
  var ajax = new XMLHttpRequest();
  ajax.onreadystatechange = function()
  {
    if(ajax.readyState == 4)
    {
      if (ajax.status == 200)
      {
        
      }
      else
      {
        self.Error("http "+ajax.status);
      }
    }
  };
  var data = new FormData();
  data.set("comment", text);
  ajax.open("POST", "/wg-api/v1/authed/comment");
  ajax.send(data);
};

modules.exports = {
  "LiveChat": Chat
};
