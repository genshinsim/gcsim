import { Executor } from "@gcsim/executors";
import { useEffect, useRef } from "react";
import { Redirect, Route, Switch, useLocation } from "wouter";
import { HotkeysProvider } from "@blueprintjs/core";
import { Provider } from "react-redux";
import { useTranslation } from "react-i18next";
import { store } from "./Stores/store";
import { Footer, Nav } from "./Components";
import {
  Dash,
  Simulator,
  ViewerLoader,
  ViewTypes,
  PageUserAccount,
  DiscordCallback,
} from "./Pages";
import "./Translation/i18n";

// all the css styling we need (except tailwind)
import "@blueprintjs/core/lib/css/blueprint.css";
import "@blueprintjs/icons/lib/css/blueprint-icons.css";
import "./index.css";

function RedirectDB() {
  window.location.replace("https://db.gcsim.app");
  return (
    <div>
      Please visit the new db at{" "}
      <a href="https://db.gcsim.app">https://db.gcsim.app</a>
    </div>
  );
}

let pool: Executor;

export function SetExecutor(executor: Executor) {
  pool = executor;
}

function UI() {
  const { t } = useTranslation();
  const content = useRef<HTMLDivElement>(null);
  const [location] = useLocation();

  useEffect(() => {
    let loc = window.location.href;
    if (loc.includes("www.gcsim.app")) {
      loc = loc.replace("www.gcsim.app", "gcsim.app");
      window.location.href = loc;
    }
  }, []);

  // every time you change location, scroll to top of page. This is necessary since the outer
  // content div will never rerender through the entire lifespan of the app and will always retain
  // its scroll position.
  useEffect(() => {
    content.current?.scrollTo(0, 0);
  }, [location]);

  return (
    <div className="bp4-dark h-screen flex flex-col">
      <Nav />
      <div
        ref={content}
        className="flex flex-col flex-auto overflow-y-scroll overflow-x-clip"
      >
        <Switch>
          <Route path="/" component={Dash} />
          <Route path="/simple">
            <Redirect to="/simulator" />
          </Route>
          <Route path="/advanced">
            <Redirect to="/simulator" />
          </Route>
          <Route path="/simulator">
            <Simulator pool={pool} />
          </Route>

          {/* Viewer Routes */}
          <Route path="/viewer">
            <ViewerLoader pool={pool} type={ViewTypes.Landing} />
          </Route>
          <Route path="/viewer/upload">
            <ViewerLoader pool={pool} type={ViewTypes.Upload} />
          </Route>
          <Route path="/viewer/web">
            <ViewerLoader pool={pool} type={ViewTypes.Web} />
          </Route>
          <Route path="/viewer/local">
            <ViewerLoader pool={pool} type={ViewTypes.Local} />
          </Route>
          <Route path="/viewer/share/:id">
            {(params) => (
              <ViewerLoader pool={pool} type={ViewTypes.Share} id={params.id} />
            )}
          </Route>

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
              {t<string>("src.this_page_is")}
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

export function App() {
  return (
    <Provider store={store}>
      <HotkeysProvider>
        <UI />
      </HotkeysProvider>
    </Provider>
  );
}
