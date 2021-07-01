import React from "react";
import { useDispatch, useSelector } from "react-redux";
import { Route, Switch, useLocation } from "wouter";
import { openSocket } from "app/appSlice";
import { RootState } from "app/store";
import { Button, Callout, Classes, Dialog, Spinner } from "@blueprintjs/core";
import Sim from "features/sim/Sim";
import Nav from "features/nav/Nav";
import Debug from "features/debug/Debug";
import Results from "features/results/Results";

function App() {
  const { socketOpen, loading, msg, hasErr, haveResults, haveDebug } =
    useSelector((state: RootState) => {
      return {
        socketOpen: state.app.isOpen,
        loading: state.sim.isLoading,
        msg: state.sim.msg,
        hasErr: state.sim.hasErr,
        haveResults: state.results.haveResult,
        haveDebug: state.debug.haveDebug,
      };
    });
  const dispatch = useDispatch();
  const [open, setIsOpen] = React.useState<boolean>(false);
  const [, setLocation] = useLocation();

  const navigate = (url: string) => {
    return () => {
      setIsOpen(false);
      setLocation(url);
    };
  };

  React.useEffect(() => {
    dispatch(openSocket());
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  React.useEffect(() => {
    if (loading === true) {
      setIsOpen(true);
    }
  }, [loading]);

  if (!socketOpen) {
    return (
      <div className="App">
        <div style={{ marginLeft: "100px", marginTop: "50px" }}>
          <div className="row center-xs">
            <Callout intent="warning">
              Connecting to server... Please wait.
            </Callout>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="App">
      <Dialog
        isOpen={open}
        canEscapeKeyClose={false}
        canOutsideClickClose={false}
      >
        {loading ? (
          <div className={Classes.DIALOG_BODY}>
            Running sim. Please wait...
            <Spinner></Spinner>
          </div>
        ) : (
          <div>
            <div className={Classes.DIALOG_BODY}>
              <Callout intent={hasErr ? "danger" : "success"}>{msg}</Callout>
            </div>
            <div className={Classes.DIALOG_FOOTER}>
              <div className={Classes.DIALOG_FOOTER_ACTIONS}>
                <Button onClick={() => setIsOpen(false)}>Close</Button>
                <Button
                  onClick={navigate("/results")}
                  intent="success"
                  disabled={!haveResults}
                >
                  View Results
                </Button>
                <Button
                  onClick={navigate("/debug")}
                  intent="primary"
                  disabled={!haveDebug}
                >
                  Debug
                </Button>
              </div>
            </div>
          </div>
        )}
      </Dialog>
      <Nav />
      <Switch>
        <Route path="/" component={Sim} />
        <Route path="/results" component={Results} />
        <Route path="/debug" component={Debug} />
      </Switch>
    </div>
  );
}

export default App;
