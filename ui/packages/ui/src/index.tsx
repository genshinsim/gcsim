import {
  Classes,
  Dialog,
  HotkeysProvider,
  Switch as SwitchInput,
} from "@blueprintjs/core";
import { Executor, ExecutorSupplier } from "@gcsim/executors";
import { ReactNode, useEffect, useRef } from "react";
import { Helmet } from "react-helmet";
import { useTranslation } from "react-i18next";
import { Provider } from "react-redux";
import {
  BrowserRouter,
  Redirect,
  Route,
  Switch,
  useHistory,
  useLocation,
} from "react-router-dom";
import "./i18n";
import {
  Dash,
  DBViewer,
  DiscordCallback,
  LocalSample,
  LocalViewer,
  PageUserAccount,
  ShareViewer,
  Simulator,
  UploadSample,
  WebViewer,
} from "./Pages";
import { Footer, Nav } from "./Sectioning";
import { appActions } from "./Stores/appSlice";
import {
  RootState,
  store,
  useAppDispatch,
  useAppSelector,
} from "./Stores/store";

// all the css styling we need (except tailwind)
import "@blueprintjs/core/lib/css/blueprint.css";
import "@blueprintjs/icons/lib/css/blueprint-icons.css";
import "@blueprintjs/popover2/lib/css/blueprint-popover2.css";
import "@blueprintjs/select/lib/css/blueprint-select.css";
import "@gcsim/components/src/index.css";
import "./index.css";

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
  mode: string;
  gitCommit: string;
};

export const UI = (props: UIProps) => {
  return (
    <BrowserRouter>
      <Provider store={store}>
        <HotkeysProvider>
          <Main {...props} />
        </HotkeysProvider>
      </Provider>
    </BrowserRouter>
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

// TODO: Move to its own file?
// TODO: Add tabs for better settings management + extensibility?
const ExecutorSettings = ({ children }: { children: ReactNode }) => {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();
  const { isOpen, sampleOnLoad } = useAppSelector((state: RootState) => {
    return {
      isOpen: state.app.isSettingsOpen,
      sampleOnLoad: state.app.sampleOnLoad,
    };
  });

  return (
    <Dialog
      isOpen={isOpen}
      onClose={() => dispatch(appActions.setSettingsOpen(false))}
      title={t<string>("simple.settings")}
      icon="settings"
      className="!pb-0"
    >
      <div className={Classes.DIALOG_BODY}>
        <>
          {children}
          <SwitchInput
            checked={sampleOnLoad}
            onChange={() => dispatch(appActions.setSampleOnLoad(!sampleOnLoad))}
            className="pt-5"
            labelElement={t<string>("simple.generate_sample")}
          />
        </>
      </div>
    </Dialog>
  );
};

const viewerPaths = ["/web", "/local", "/sh/", "/db/"];

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function movedOffViewer(location: any, prevLocation: any): boolean {
  let prevWasViewer = false;
  let destIsViewer = false;
  for (let i = 0; i < viewerPaths.length; i++) {
    prevWasViewer =
      prevWasViewer || prevLocation.current.pathname.startsWith(viewerPaths[i]);
    destIsViewer = destIsViewer || location.pathname.startsWith(viewerPaths[i]);
  }
  return prevWasViewer && !destIsViewer;
}

const Main = ({ exec, children, gitCommit, mode }: UIProps) => {
  const { t } = useTranslation();
  const content = useRef<HTMLDivElement>(null);
  const location = useLocation();
  const history = useHistory();

  // every time you change location, scroll to top of page. This is necessary since the outer
  // content div will never rerender through the entire lifespan of the app and will always retain
  // its scroll position.
  useEffect(() => {
    const unlisten = history.listen(() => {
      content.current?.scrollTo(0, 0);
    });
    return () => unlisten();
  }, [history]);

  // cancel the run every time we navigate away from the web viewer page
  const prevLocation = useRef(location);
  useEffect(() => {
    if (
      prevLocation.current != location &&
      movedOffViewer(location, prevLocation) &&
      exec().running()
    ) {
      exec().cancel();
    }
    prevLocation.current = location;
  }, [location, exec]);

  return (
    <div className="bp4-dark h-screen flex flex-col">
      <Nav />
      <div
        ref={content}
        className="flex flex-col flex-auto overflow-y-scroll overflow-x-clip"
      >
        <Switch>
          <Route exact path="/">
            <Helmet>
              <title>gcsim - simulation impact</title>
            </Helmet>
            <Dash />
          </Route>

          {/* Simulator */}
          <Route exact path="/simulator">
            <Helmet>
              <title>gcsim - simulator</title>
            </Helmet>
            <Simulator exec={exec} />
          </Route>

          {/* Viewer Routes */}
          <Route path="/web">
            <Helmet>
              <title>gcsim - viewer</title>
            </Helmet>
            <WebViewer exec={exec} gitCommit={gitCommit} mode={mode} />
          </Route>
          <Route path="/local">
            <Helmet>
              <title>gcsim - local viewer</title>
            </Helmet>
            <LocalViewer exec={exec} gitCommit={gitCommit} mode={mode} />
          </Route>
          <Route path="/sh/:id">
            {({ match }) => {
              document.title = "gcsim sh - " + match?.params.id;
              return (
                <ShareViewer
                  exec={exec}
                  id={match?.params.id}
                  gitCommit={gitCommit}
                  mode={mode}
                />
              );
            }}
          </Route>
          <Route path="/db/:id">
            {({ match }) => {
              document.title = "gcsim db - " + match?.params.id;
              return (
                <DBViewer
                  exec={exec}
                  id={match?.params.id}
                  gitCommit={gitCommit}
                  mode={mode}
                />
              );
            }}
          </Route>

          {/* Sample Routes */}
          <Route path="/sample/upload">
            <Helmet>
              <title>gcsim - sample</title>
            </Helmet>
            <UploadSample />
          </Route>
          <Route path="/sample/local">
            <Helmet>
              <title>gcsim - local sample</title>
            </Helmet>
            <LocalSample />
          </Route>

          {/* Redirects */}
          <Route path={["/v3/viewer/share/:id", "/viewer/share/:id", "/s/:id"]}>
            {({ match }) => <Redirect to={"/sh/" + match?.params.id} />}
          </Route>
          <Route path={"/viewer/web"}>
            <Redirect to="/web" />
          </Route>
          <Route path={"/viewer/local"}>
            <Redirect to="/local" />
          </Route>
          <Route path={["/simple", "/advanced", "/viewer"]}>
            <Redirect to="/simulator" />
          </Route>

          {/* DB & Account */}
          <Route path="/db">
            <RedirectDB />
          </Route>
          <Route path="/account">
            <Helmet>
              <title>gcsim - account</title>
            </Helmet>
            <PageUserAccount />
          </Route>
          <Route path="/auth/discord">
            <DiscordCallback />
          </Route>

          {/* Default (404 case) */}
          <Route>
            <Helmet>
              <title>gcsim - simulation impact</title>
            </Helmet>
            <div className="m-2 text-center">
              {t<string>("src.this_page_is")}
            </div>
          </Route>
        </Switch>
        <Footer />
        <ExecutorSettings>{children}</ExecutorSettings>
      </div>
    </div>
  );
};
