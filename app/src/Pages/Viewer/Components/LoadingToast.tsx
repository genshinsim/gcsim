import { Button, Classes, Intent, ProgressBar, Toaster } from "@blueprintjs/core";
import classNames from "classnames";
import React, { useRef } from "react";


// TODO: Add translations
export default function useLoadingToast(current?: number, total?: number) {
  const loadingToast = useRef<Toaster>(null);

  React.useEffect(() => {
    const val = (current == undefined || total == undefined) ? 0 : current / total;
    const content = (
      <div className="flex flex-row items-center justify-between gap-2">
        <div className="min-w-fit">Running ({current ?? "??"}/{total ?? "??"})</div>
        <ProgressBar
            className={classNames("basis-1/2 flex-auto sm:min-w-", {
              [Classes.PROGRESS_NO_STRIPES]: val >= 1,
            })}
            intent={val < 1 ? Intent.PRIMARY : Intent.SUCCESS}
            value={val}/>
        {action(val, loadingToast)}
      </div>
    );

    loadingToast.current?.show({
      message: content,
      className: "w-full !max-w-2xl",
      intent: Intent.NONE,
      isCloseButtonShown: false,
      timeout: val < 1 ? 0 : 2000
    });
  }, [current, total]);

  return loadingToast;
}

// TODO: Abort runs callback to stop execution
function action(val: number, loadingToast: React.RefObject<Toaster>) {
  const cls = classNames("!min-w-fit");

  if (val < 1.0) {
    return (
      <Button
          className={cls}
          intent={Intent.DANGER}>
        Abort
      </Button>
    );
  } 
  return (
    <Button
        className={cls}
        onClick={() => loadingToast.current?.clear()}
        intent={Intent.SUCCESS}>
      Dismiss
    </Button>
  );
}