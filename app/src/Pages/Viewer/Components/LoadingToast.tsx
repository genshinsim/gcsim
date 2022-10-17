import { Button, Classes, Intent, ProgressBar, Toaster } from "@blueprintjs/core";
import classNames from "classnames";
import React, { useRef } from "react";
import { ViewTypes } from "..";

// TODO: Add translations + number format
export default function useLoadingToast(
    type: ViewTypes, error: string | null, cancel?: () => void,
    current?: number, total?: number) {
  const loadingToast = useRef<Toaster>(null);
  const key = useRef<string | undefined>(undefined);

  React.useEffect(() => {
    if (error != null) {
      loadingToast.current?.clear();
      return;
    }

    if (current == undefined || total == undefined) {
      key.current = loadingToast.current?.show({
        message: "Loading...",
        icon: "refresh",
        intent: Intent.PRIMARY,
        isCloseButtonShown: false,
        timeout: 0,
      }, key.current);
      return;
    }

    if (current == total && type != ViewTypes.Web) {
      key.current = loadingToast.current?.show({
        message: "Loaded " + current + " iterations!",
        icon: "tick",
        intent: Intent.SUCCESS,
        isCloseButtonShown: true,
        timeout: 2000,
      }, key.current);
      return;
    }

    const val = current / total;
    const content = (
      <div className="flex flex-row items-center justify-between gap-2">
        <div className="min-w-fit">Running ({current}/{total})</div>
        <ProgressBar
            className={classNames("basis-1/2 flex-auto sm:min-w-", {
              [Classes.PROGRESS_NO_STRIPES]: val >= 1,
            })}
            intent={val < 1 ? Intent.PRIMARY : Intent.SUCCESS}
            value={val}/>
        {action(val, loadingToast, cancel)}
      </div>
    );

    key.current = loadingToast.current?.show({
      message: content,
      className: "w-full !max-w-2xl",
      intent: Intent.NONE,
      isCloseButtonShown: val >= 1,
      timeout: val < 1 ? 0 : 2000
    }, key.current);
  }, [current, total, type, cancel, error]);
  return loadingToast;
}

function action(val: number, loadingToast: React.RefObject<Toaster>, cancel?: () => void) {
  if (val >= 1 || cancel == null) {
    return null;
  }

  return (
    <Button
        className="!min-w-fit"
        intent={Intent.DANGER}
        onClick={() => {
          cancel();
          loadingToast.current?.clear();
        }}>
      Cancel
    </Button>
  );
}