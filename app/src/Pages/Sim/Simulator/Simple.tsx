import { Callout, useHotkeys } from "@blueprintjs/core";
import React from "react";

import { SectionDivider } from "~src/Components/SectionDivider";
import { Viewport } from "~src/Components/Viewport";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import { simActions } from "..";
import { ActionList, SimProgress } from "../Components";
import { Team } from "./Team";
import { Trans, useTranslation } from "react-i18next";
import { Toolbox } from "./Toolbox";
import { ActionListTooltip, TeamBuilderTooltip } from "./Tooltips";
import { updateCfg } from "../simSlice";

export function Simple() {
  let { t } = useTranslation();

  const { cfg, cfg_err, showBuilder } = useAppSelector((state: RootState) => {
    return {
      cfg: state.sim.cfg,
      cfg_err: state.sim.cfg_err,
      showBuilder: state.sim.showBuilder,
    };
  });
  const dispatch = useAppDispatch();

  const [open, setOpen] = React.useState<boolean>(false);

  const hotkeys = React.useMemo(
    () => [
      {
        combo: "Esc",
        global: true,
        label: t("simple.exit_edit"),
        onKeyDown: () => {
          dispatch(simActions.editCharacter({ index: -1 }));
        },
      },
    ],
    []
  );
  useHotkeys(hotkeys);

  return (
    <Viewport className="flex flex-col gap-2">
      <div className="flex flex-col gap-2">
        <div className="flex flex-col">
          {showBuilder ? (
            <>
              <SectionDivider>
                <Trans>simple.team</Trans>
              </SectionDivider>
              <TeamBuilderTooltip />
              <Team />
            </>
          ) : null}

          <SectionDivider>
            <Trans>simple.action_list</Trans>
          </SectionDivider>

          <ActionListTooltip />

          <ActionList cfg={cfg} onChange={(v) => dispatch(updateCfg(v))} />

          <div className="sticky bottom-0 bg-bp-bg flex flex-col gap-y-1">
            {cfg_err !== "" ? (
              <div className="basis-full p-1">
                <Callout intent="warning" title="Error parsing config">
                  <pre className=" whitespace-pre-wrap">{cfg_err}</pre>
                </Callout>
              </div>
            ) : null}

            <Toolbox canRun={cfg_err === ""} />
          </div>
        </div>
      </div>
      <SimProgress isOpen={open} onClose={() => setOpen(false)} />
    </Viewport>
  );
}
