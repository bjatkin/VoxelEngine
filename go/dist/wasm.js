'use strict';

const WASM_URL = 'main.wasm';
var wasm;

// Prevent the context menue on the canvas
document.getElementById("gocanvas").oncontextmenu = (e) => {
  e.preventDefault();
}

// Render one frame of the animation
function renderFrame(evt) {
  wasm.exports.renderFrame(evt);
}

// Load and run the wasm
function init() {
  const go = new Go();
  if ('instantiateStreaming' in WebAssembly) {
    WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then(function (obj) {
      wasm = obj.instance;
      go.run(wasm); // Initial setup
    })
  } else {
    fetch(WASM_URL).then(resp =>
      resp.arrayBuffer()
    ).then(bytes =>
      WebAssembly.instantiate(bytes, go.importObject).then(function (obj) {
        wasm = obj.instance;
        go.run(wasm);
      })
    )
  }
}

init();