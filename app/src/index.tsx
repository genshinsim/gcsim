import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import { store } from "~src/store";
import { Provider } from "react-redux";
import "@blueprintjs/core/lib/css/blueprint.css";
import "@blueprintjs/popover2/lib/css/blueprint-popover2.css";
import "@blueprintjs/select/lib/css/blueprint-select.css";
import { HotkeysProvider } from "@blueprintjs/core";

console.log(store);

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <HotkeysProvider>
        <App />
      </HotkeysProvider>
    </Provider>
  </React.StrictMode>,
  document.getElementById("app")
);
