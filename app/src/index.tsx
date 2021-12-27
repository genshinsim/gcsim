import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import GA4React from "ga-4-react";

const ga4react = new GA4React("G-XXXXXXXXXX");

(async (_) => {
  await ga4react
    .initialize()
    .then((res) => console.log("Analytics Success."))
    .catch((err) => console.log("Analytics Failure."))
    .finally(() => {
      ReactDOM.render(
        <React.StrictMode>
          <App />
        </React.StrictMode>,
        document.getElementById("root")
      );
    });
})();

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
// reportWebVitals();
