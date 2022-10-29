import { Executor, ExecutorSupplier } from "@gcsim/executors";
import { ReactNode, useEffect, useRef } from "react";
import { Redirect, Route, Switch, useLocation } from "wouter";
import { Classes, Dialog, HotkeysProvider } from "@blueprintjs/core";
import { Provider } from "react-redux";
import { useTranslation } from "react-i18next";
import { RootState, store, useAppDispatch, useAppSelector } from "./Stores/store";
import { appActions } from "./Stores/appSlice";
import { Footer, Nav } from "./Sectioning";
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
import "@blueprintjs/select/lib/css/blueprint-select.css";
import "@blueprintjs/popover2/lib/css/blueprint-popover2.css";
import "./index.css";

// helper functions
export { useLocalStorage } from "./Util";

// Two things must always be supplied to the UI for it to work
//  1. ExecutorSupplier (Should rarely/never change. Must be react safe)
//  2. ExecutorSettings (passed as children, options and executor state management owned by app)
//
// We use a supplier to give the owning app the opporunity to add their own logic every time we
// try to access the executor. This means they can defer creation to only when it is used
// (improving performance and UX), auto refresh/recreate it if it detects a stale state, or change
// what pool instance is in use
//
// ExecutorSettings is a dialog which will be added to the DOM as part of the main content
// so that it is always loaded in. This is how the app can supply state and decide how it wants to
// construct and configure the executors (and which executors are available to use).
type UIProps = {
  exec: ExecutorSupplier<Executor>;
  children: ReactNode;
};

export const UI = (props: UIProps) => {
  return (
    <Provider store={store}>
      <HotkeysProvider>
        <Main {...props} />
      </HotkeysProvider>
    </Provider>
  );
};

function RedirectDB() {
  window.location.replace("https://db.gcsim.app");
  return (
    <div>
      Please visit the new db at{" "}
      <a href="https://db.gcsim.app">https://db.gcsim.app</a>
    </div>
  );
}

const ExecutorSettings = ({ children }: { children: ReactNode }) => {
  const dispatch = useAppDispatch();
  const { isOpen } = useAppSelector((state: RootState) => {
    return {
      isOpen: state.app.isSettingsOpen,
    };
  });

  const close = () => dispatch(appActions.setSettingsOpen(false));

  return (
    <Dialog
        isOpen={isOpen}
        onClose={close}
        title="Executor Settings"
        icon="settings">
      <div className={Classes.DIALOG_BODY}>{children}</div>
    </Dialog>
  );
};

const Main = ({ exec, children }: UIProps) => {
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

  // cancel the run every time we navigate away from the viewer page
  const prevLocation = useRef(location);
  useEffect(() => {
    if (prevLocation.current != location && prevLocation.current.startsWith("/viewer/")
        && exec().running()) {
      exec().cancel();
    }
    prevLocation.current = location;
  }, [location, exec]);

  return (
    <div className="bp4-dark h-screen flex flex-col">
      <Nav />
      <div ref={content} className="flex flex-col flex-auto overflow-y-scroll overflow-x-clip">
        <Switch>
          <Route path="/" component={Dash} />
          <Route path="/simple">
            <Redirect to="/simulator" />
          </Route>
          <Route path="/advanced">
            <Redirect to="/simulator" />
          </Route>
          <Route path="/simulator">
            <Simulator exec={exec} />
          </Route>

          {/* Viewer Routes */}
          <Route path="/viewer">
            <ViewerLoader exec={exec} type={ViewTypes.Landing} />
          </Route>
          <Route path="/viewer/upload">
            <ViewerLoader exec={exec} type={ViewTypes.Upload} />
          </Route>
          <Route path="/viewer/web">
            <ViewerLoader exec={exec} type={ViewTypes.Web} />
          </Route>
          <Route path="/viewer/local">
            <ViewerLoader exec={exec} type={ViewTypes.Local} />
          </Route>
          <Route path="/viewer/share/:id">
            {(params) => (
              <ViewerLoader exec={exec} type={ViewTypes.Share} id={params.id} />
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
        <Footer />
        <ExecutorSettings>
          {children}
        </ExecutorSettings>
      </div>
    </div>
  );
};