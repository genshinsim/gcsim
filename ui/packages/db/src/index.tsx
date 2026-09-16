import { initI18n } from "@gcsim/localization";
import App from "./App";

// all the css styling we need (except tailwind)
import "@blueprintjs/core/lib/css/blueprint.css";
import "@blueprintjs/icons/lib/css/blueprint-icons.css";
import "@blueprintjs/popover2/lib/css/blueprint-popover2.css";
import "@blueprintjs/select/lib/css/blueprint-select.css";
import "@gcsim/components/src/index.css";
import "./index.css";

import React from "react";
import ReactDOM from "react-dom/client";

initI18n();

const root = ReactDOM.createRoot(
	document.getElementById("root") as HTMLElement,
);
root.render(
	<React.StrictMode>
		<App />
	</React.StrictMode>,
);
