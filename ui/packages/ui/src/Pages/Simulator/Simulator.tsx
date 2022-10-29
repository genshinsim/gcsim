import { Callout, Intent } from "@blueprintjs/core";
import { useEffect, useRef, useState } from "react";
import { Viewport, SectionDivider } from "../../Components";
import { ActionList } from "./Components";
import { Team } from "./Team";
import { Trans } from "react-i18next";
import { Toolbox } from "./Toolbox";
import { ActionListTooltip, TeamBuilderTooltip } from "./Tooltips";
import { useAppSelector, RootState, useAppDispatch } from "../../Stores/store";
import { Executor, ExecutorSupplier } from "@gcsim/executors";
import { appActions, defaultStats } from "../../Stores/appSlice";
import { Character } from "@gcsim/types";
import { debounce } from "lodash-es";

export function Simulator({ exec }: { exec: ExecutorSupplier<Executor> }) {
  const dispatch = useAppDispatch();
  const { settings, initCfg } = useAppSelector(
    (state: RootState) => {
      return {
        initCfg: state.app.cfg,
        settings: state.user.settings,
      };
    }
  );

  // use a local cfg and only update redux at a much lower rate (1s after typing stops).
  // Other stuff happens on redux cfg update which lags the editor
  // Note: seems like the editor and/or highlighter we use is just laggy, so it doesn't help much :/
  const [cfg, setCfg] = useState(initCfg);
  const setReduxCfgThrottle = useRef(debounce(
    (cfg) => dispatch(appActions.setCfg({ cfg: cfg, keepTeam: false })), 1000));
  const onChange = (newCfg: string) => {
    setCfg(newCfg);
    setReduxCfgThrottle.current(newCfg);
  };

  // check worker ready state every 250ms so run button becomes available when workers do
  const [isReady, setReady] = useState<boolean | null>(null);
  useEffect(() => {
    const interval = setInterval(() => {
      setReady(exec().ready());
    }, 250);
    return () => clearInterval(interval);
  }, [exec]);

  const [err, setErr] = useState("");

  // will detect changes in the redux config and validate with the executor
  // validated == true means we had a successful validation check run, not that it is valid
  const validated = useConfigValidateListener(exec, cfg, isReady, setErr);

  document.title = "gcsim - simulator";
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

          <ActionList cfg={cfg} onChange={onChange} />

          <div className="sticky bottom-0 bg-bp-bg flex flex-col gap-y-1">
            {err !== "" && cfg !== "" ? (
              <div className="pl-2 pr-2 pt-2 mt-1">
                <Callout intent={Intent.DANGER} title="Error: Config Invalid">
                  <pre className="whitespace-pre-wrap pl-5">{err}</pre>
                </Callout>
              </div>
            ) : null}
            <Toolbox
                exec={exec}
                cfg={cfg}
                isReady={isReady === true}
                isValid={err === "" && validated} />
          </div>
        </div>
      </div>
    </Viewport>
  );
}

export function useConfigValidateListener(
    exec: ExecutorSupplier<Executor>, cfg: string, isReady: boolean | null,
    setErr: (str: string) => void): boolean {
  const dispatch = useAppDispatch();
  const [validated, setValidated] = useState(false);
  const debounced = useRef(debounce((x: () => void) => x(), 200));

  useEffect(() => {
    if (!isReady || cfg === "") {
      return;
    }

    setValidated(false);
    debounced.current(() => {
      exec().validate(cfg).then(
        (res) => {
          console.log("all is good");
          setErr("");
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
            setErr(msg);
          }
          dispatch(appActions.setTeam(team));
          setValidated(true);
        },
        (err) => {
          //set error state
          setErr(err);
          setValidated(false);
        }
      );
    });
  }, [exec, cfg, dispatch, setErr, isReady]);

  return validated;
}