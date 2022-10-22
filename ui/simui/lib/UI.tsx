import React, { useRef } from "react";
import { Redirect, Route, Switch, useLocation } from "wouter";
import Footer from "./Components/Footer/Footer";
import Nav from "./Components/Nav/Nav";
import { Dash } from "./Pages/Dash";
import { Simulator } from "./Pages/Simulator";
import "./Translation/i18n";
import { Trans, useTranslation } from "react-i18next";
import { DiscordCallback } from "./Pages/User/DiscordCallback";
import { PageUserAccount } from "./Pages/User/PageUserAccount";
import "./index.css";
import { store } from "~/Stores/store";
import { Provider } from "react-redux";
import "@blueprintjs/core/lib/css/blueprint.css";
import "@blueprintjs/popover2/lib/css/blueprint-popover2.css";
import "@blueprintjs/select/lib/css/blueprint-select.css";
import { HotkeysProvider } from "@blueprintjs/core";

function RedirectDB() {
  window.location.replace("https://db.gcsim.app");
  return (
    <div>
      Please visit the new db at <a href="https://db.gcsim.app">https://db.gcsim.app</a>
    </div>
  );
}

function UI() {
  useTranslation();
  const content = useRef<HTMLDivElement>(null);
  const [location] = useLocation();

  React.useEffect(() => {
    let loc = window.location.href;

    if (loc.includes("www.gcsim.app")) {
      loc = loc.replace("www.gcsim.app", "gcsim.app");
      window.location.href = loc;
    }
  }, []);

  // every time you change location, scroll to top of page. This is necessary since the outer
  // content div will never rerender through the entire lifespan of the app and will always retain
  // its scroll position.
  React.useEffect(() => {
    content.current?.scrollTo(0, 0);
  }, [location]);

  return (
    <div className="bp4-dark h-screen flex flex-col">
      <Nav />
      <div ref={content} className="flex flex-col flex-auto overflow-y-scroll overflow-x-clip">
        <Switch>
          <Route path="/" component={Dash} />
          <Route path="/simulator">
            <Simulator />
          </Route>

          {/* Viewer Routes */}
          {/* <Route path="/viewer">
            <ViewerLoader type={ViewTypes.Landing} />
          </Route>
          <Route path="/viewer/upload">
            <ViewerLoader type={ViewTypes.Upload} />
          </Route>
          <Route path="/viewer/web">
            <ViewerLoader type={ViewTypes.Web} />
          </Route>
          <Route path="/viewer/local">
            <ViewerLoader type={ViewTypes.Local} />
          </Route>
          <Route path="/viewer/share/:id">
            {(params) => <ViewerLoader type={ViewTypes.Share} id={params.id} />}
          </Route> */}

          {/* reroute v3 -> new viewer */}
          <Route path="/v3/viewer/share/:id">
            {(params) => <Redirect to={"/viewer/share/" + params.id} />}
          </Route>

          <Route path="/db">
            <RedirectDB />
          </Route>
          <Route path="/account">
            <PageUserAccount />
          </Route>
          <Route path="/auth/discord">
            <DiscordCallback />
          </Route>
          <Route>
            <div className="m-2 text-center">
              <Trans>src.this_page_is</Trans>
            </div>
          </Route>
        </Switch>
        <div className="w-full pt-4 pb-4 md:pl-4">
          <Footer />
        </div>
      </div>
    </div>
  );
}

export default function App() {
  return (
    <Provider store={store}>
      <HotkeysProvider>
        <UI />
      </HotkeysProvider>
    </Provider>
  );
}
