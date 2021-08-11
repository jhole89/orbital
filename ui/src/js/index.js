"use strict";

require("/src/js/graph.js");
var Elm = require("/src/elm/Main.elm").Elm;

var app = Elm.Main.init({
  node: document.getElementById("orbital-ui"),
});
