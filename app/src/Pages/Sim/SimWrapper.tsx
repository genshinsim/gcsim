import { Callout, Spinner } from "@blueprintjs/core";
import React from "react";
import { useAppSelector, RootState, useAppDispatch } from "~src/store";
import { loadWorkers } from ".";
import { Main } from "./Components";

export function SimWrapper({ children }: { children: React.ReactNode }) {
  const { ready } = useAppSelector((state: RootState) => {
    return {
      ready: state.sim.ready,
    };
  });
  const dispatch = useAppDispatch();

  React.useEffect(() => {
    dispatch(loadWorkers());
  }, []);

  if (ready === 0) {
    return (
      <Main>
        <Callout intent="primary" title="Loading simulator. Please wait">
          <Spinner />
        </Callout>
      </Main>
    );
  }
  return <div>{children}</div>;
}
