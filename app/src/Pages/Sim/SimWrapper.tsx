import { Callout, Spinner } from "@blueprintjs/core";
import React from "react";
import { Viewport } from "~src/Components/Viewport";
import { useAppSelector, RootState, useAppDispatch } from "~src/store";
import { setTotalWorkers } from ".";
import { useTranslation } from "react-i18next";

export function SimWrapper({ children }: { children: React.ReactNode }) {
  let { t } = useTranslation()

  const { ready, workers } = useAppSelector((state: RootState) => {
    return {
      ready: state.sim.ready,
      workers: state.sim.workers,
    };
  });
  const dispatch = useAppDispatch();

  React.useEffect(() => {
    dispatch(setTotalWorkers(workers));
  }, []);

  if (ready === 0) {
    return (
      <Viewport>
        <Callout intent="primary" title={t("sim.loading_simulator_please")}>
          <Spinner />
        </Callout>
      </Viewport>
    );
  }
  return <div>{children}</div>;
}
