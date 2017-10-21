/** utils.js */


var get_id = function(id) {
  return document.getElementById(id) || null;
};

var for_each_class = function(cl, f) {
  var elems = document.getElementsByClassName(cl);
  for (var idx = 0; idx < elems.length; idx ++)
    f(elems[idx]);
};

var new_elem = function(tag, id, cl) {
  var e = document.createElement(tag);
  if(cl)
    e.setAttribute("class", cl);
  if(id)
    e.setAttribute("id", id);
  return e;
};

var new_div = function(id, cl) {
  return new_elem("div", id, cl);
};

var new_button = function(id, cl) {
  var btn = new_elem("input", id, cl);
  btn.setAttribute("type", "button");
  return btn;
};

var new_submit = function(id, cl) {
  var submit = new_elem("input", id, cl);
  submit.setAttribute("type", "submit");
  return submit;
};

var new_form = function(url, method) {
  var b = new_elem("form");
  b.action = url || window.location.href;
  b.method = method || "POST";
  b.enctype = "multipart/form-data";
  return b;
};

var new_input = function(id, cl) {
  return new_elem("input", id, cl);
};

var new_form_file = function(id, cl) {
  var f = new_elem("input", id, cl);
  f.setAttribute("type", "file");
  return f;
};

module.exports = {
  "get_id": get_id,
  "for_each_class": for_each_class,
  "div": new_div,
  "button": new_button,
  "form": new_form,
  "file": new_form_file,
  "submit": new_submit,
  "input": new_input
};
