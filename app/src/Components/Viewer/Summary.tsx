import Graphs from "./Graphs/Graphs";
import { SimResults } from "./DataType";
import TeamView from "./Team/TeamView";
import DPSOverTime from "./Graphs/DPSOverTime";
import { Trans, useTranslation } from "react-i18next";
import {
  ButtonGroup,
  Button,
  Callout,
  Checkbox,
  Classes,
  Dialog,
} from "@blueprintjs/core";
import React from "react";
import { useLocation } from "wouter";
import { updateCfg } from "~src/Pages/Sim";
import { useAppDispatch } from "~src/store";
import { AppToaster } from "./Share";

const LOCALSTORAGE_KEY = "gcsim-viewer-cpy-cfg-settings";

export default function Summary({ data }: { data: SimResults }) {
  let { t } = useTranslation();

  const [open, setOpen] = React.useState<boolean>(false);
  const dispatch = useAppDispatch();
  const [_, setLocation] = useLocation();
  const [keepExistingTeam, setKeepExistingTeam] = React.useState<boolean>(
    () => {
      const saved = localStorage.getItem(LOCALSTORAGE_KEY);
      if (saved === "true") {
        return true;
      }
      return false;
    }
  );

  function copyToClipboard() {
    navigator.clipboard.writeText(data.config_file).then(() => {
      AppToaster.show({
        message: t("viewer.copied_to_clipboard"),
        intent: "success",
      });
    });
    // TODO: Need to add a blueprintjs Toaster for ephemeral confirmation box
  }

  const openInSim = () => {
    setOpen(false);
    dispatch(updateCfg(data.config_file, keepExistingTeam));
    setLocation("/simulator");
  };

  const handleToggleSelected = () => {
    localStorage.setItem(LOCALSTORAGE_KEY, keepExistingTeam ? "false" : "true");
    setKeepExistingTeam(!keepExistingTeam);
  };

  //calculate per target damage
  let trgs: JSX.Element[] = [];

  for (const key in data.dps_by_target) {
    trgs.push(
      <div className="w-full flex flex-row" key={key}>
        <span className="w-24">
          <span className="pl-2" />
          {`${t("viewer.target")} ${key}:`}
        </span>
        <div className="grid grid-cols-4 grow">
          <span className="text-right">-</span>
          <span className="text-right">
            {data.dps_by_target[key].mean.toLocaleString(undefined, {
              maximumFractionDigits: 0,
              minimumFractionDigits: 0,
            })}
          </span>
          <span className="text-right">
            {(
              (100 * data.dps_by_target[key].mean) /
              data.dps.mean
            ).toLocaleString(undefined, {
              maximumFractionDigits: 0,
              minimumFractionDigits: 0,
            })}
            {"%"}
          </span>
          <span className="text-right">
            {data.dps_by_target[key].sd
              ? data.dps_by_target[key].sd!.toLocaleString(undefined, {
                  maximumFractionDigits: 0,
                  minimumFractionDigits: 0,
                })
              : "-"}
          </span>
        </div>
      </div>
    );
  }

  return (
    <div className="wide:w-[70rem] ml-auto mr-auto">
      <div className="flex flex-row justify-end m-2">
        <ButtonGroup>
          <Button onClick={copyToClipboard} icon="clipboard">
            <Trans>viewer.copy</Trans>
          </Button>
          <Button onClick={() => setOpen(true)} icon="send-to">
            <Trans>viewer.send_to_simulator</Trans>
          </Button>
        </ButtonGroup>
      </div>
      <TeamView team={data.char_details} />
      <div className="bg-gray-600 relative rounded-md p-2 m-2 pt-10">
        <DPSOverTime data={data} />
        <div className="w-full text-center">
          <Trans>viewer.sec_pre</Trans>
          {data.sim_duration.mean.toLocaleString(undefined, {
            maximumFractionDigits: 2,
          })}
          <Trans>viewer.sec_post</Trans>
          {data.iter}
          <Trans>viewer.time_pre</Trans>
          {(data.runtime / 1000000000).toFixed(3)}
          <Trans>viewer.time_post</Trans>
          <br />
          <Trans>viewer.git_hash</Trans>
          {data.version ? (
            <a
              href={
                "https://github.com/genshinsim/gcsim/commits/" + data.version
              }
            >
              {data.version.substring(0, 8)}
            </a>
          ) : (
            t("viewer.unknown")
          )}
          <Trans>viewer.built_on</Trans>
          {data.build_date ? data.build_date : t("viewer.unknown")}
        </div>
        <div className=" pl-4 pt-2 flex flex-row place-content-center">
          <div className="max-w-4xl w-full flex flex-col gap-1">
            <div className="flex flex-row border-solid border-b-2 font-bold">
              <span className="w-24">
                <Trans>viewer.target</Trans>
              </span>
              <div className="grid grid-cols-4 grow">
                <span className="text-right">
                  <Trans>viewer.level</Trans>
                </span>
                <span className="text-right">
                  <Trans>viewer.avg_dps</Trans>
                </span>
                <span className="text-right">%</span>
                <span className="text-right">Std. Dev.</span>
              </div>
            </div>
            {trgs}
            <div className="w-full flex flex-row border-solid border-t-2 font-bold">
              <span className="w-24">
                <Trans>viewer.combined</Trans>
              </span>
              <div className="grid grid-cols-4 grow">
                <span className="text-right"></span>
                <span className="text-right">
                  {" "}
                  {data.dps.mean.toLocaleString(undefined, {
                    maximumFractionDigits: 0,
                    minimumFractionDigits: 0,
                  })}
                </span>
                <span className="text-right"></span>
                <span className="text-right">
                  {data.dps.sd?.toLocaleString(undefined, {
                    maximumFractionDigits: 0,
                    minimumFractionDigits: 0,
                  })}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <Graphs data={data} />

      <Dialog isOpen={open} onClose={() => setOpen(false)}>
        <div className={Classes.DIALOG_BODY}>
          <Trans>viewer.load_this_configuration</Trans>
          <Callout intent="warning" className="mt-2">
            <Trans>viewer.this_will_overwrite</Trans>
          </Callout>
          <Checkbox
            label="Copy action list only (ignore character stats)"
            className="mt-2"
            checked={keepExistingTeam}
            onClick={handleToggleSelected}
          />
        </div>

        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button onClick={openInSim} intent="primary">
              <Trans>viewer.continue</Trans>
            </Button>
            <Button onClick={() => setOpen(false)}>
              <Trans>viewer.cancel</Trans>
            </Button>
          </div>
        </div>
      </Dialog>
    </div>
  );
}
