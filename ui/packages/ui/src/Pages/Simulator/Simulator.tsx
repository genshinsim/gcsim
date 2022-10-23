import { Callout } from "@blueprintjs/core";
import React, { useEffect, useState } from "react";
import { Viewport, SectionDivider } from "../../Components";
import { ActionList } from "./Components";
import { Team } from "./Team";
import { Trans } from "react-i18next";
import { Toolbox } from "./Toolbox";
import { ActionListTooltip, TeamBuilderTooltip } from "./Tooltips";
import { useAppSelector, RootState, useAppDispatch } from "../../Stores/store";
import { ExecutorSupplier } from "@gcsim/executors";
import { appActions, defaultStats } from "../../Stores/appSlice";
import { Character } from "@gcsim/types";

export function Simulator({ exec }: { exec: ExecutorSupplier }) {
  const dispatch = useAppDispatch();

  const { settings, cfg, cfgErr } = useAppSelector(
    (state: RootState) => {
      return {
        cfg: state.app.cfg,
        cfgErr: state.app.cfg_err,
        settings: state.user.settings ?? {
          showTips: false,
          showBuilder: false,
        },
      };
    }
  );

  // check worker ready state every 250ms so run button becomes available when workers do
  const [isReady, setReady] = React.useState<boolean | null>(null);
  useEffect(() => {
    const interval = setInterval(() => {
      setReady(exec().ready());
    }, 250);
    return () => clearInterval(interval);
  }, [exec]);

  // will detect changes in the redux config and validate with the executor
  // validated == true means we had a successful validation check run, not that it is valid
  const validated = useConfigValidateListener(exec, cfg, isReady);

  return (
    <Viewport className="flex flex-col gap-2">
      <div className="flex flex-col gap-2">
        <div className="flex flex-col">
          {settings.showBuilder ? (
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

          <ActionList
            cfg={cfg}
            onChange={(newCfg) => dispatch(appActions.setCfg({ cfg: newCfg, keepTeam: false }))}
          />

          <div className="sticky bottom-0 bg-bp-bg flex flex-col gap-y-1">
            {cfgErr !== "" ? (
              <div className="pl-2 pr-2 pt-2 mt-1">
                <Callout intent="warning" title="Error parsing config">
                  <pre className=" whitespace-pre-wrap">{cfgErr}</pre>
                </Callout>
              </div>
            ) : null}
            <Toolbox
                exec={exec}
                cfg={cfg}
                canRun={cfgErr === "" && isReady === true && validated} />
          </div>
        </div>
      </div>
    </Viewport>
  );
}

export function useConfigValidateListener(
    exec: ExecutorSupplier, cfg: string, isReady: boolean | null): boolean {
  const dispatch = useAppDispatch();
  const [validated, setValidated] = useState(false);

  useEffect(() => {
    if (!isReady) {
      return;
    }

    exec().validate(cfg).then(
      (res) => {
        console.log("all is good");
        dispatch(appActions.setCfgErr(""));
        //if successful then we're going to update the team based on the parsed results
        let team: Character[] = [];
        if (res.characters) {
          team = res.characters.map((c) => {
            return {
              name: c.base.key,
              level: c.base.level,
              element: c.base.element,
              max_level: c.base.max_level,
              cons: c.base.cons,
              weapon: c.weapon,
              talents: c.talents,
              stats: c.stats,
              snapshot: defaultStats,
              sets: c.sets,
            };
          });
        }
        //check if there are any warning msgs
        if (res.errors) {
          let msg = "";
          res.errors.forEach((err) => {
            msg += err + "\n";
          });
          dispatch(appActions.setCfgErr(msg));
        }
        dispatch(appActions.setTeam(team));
        setValidated(true);
      },
      (err) => {
        //set error state
        dispatch(appActions.setCfgErr(err));
        setValidated(false);
      }
    );
  }, [exec, cfg, dispatch, isReady]);

  return validated;
}