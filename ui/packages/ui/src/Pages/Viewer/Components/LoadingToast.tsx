import {
  Button,
  Classes,
  Intent,
  Position,
  ProgressBar,
  Toaster,
} from "@blueprintjs/core";
import classNames from "classnames";
import { MutableRefObject, RefObject, useEffect, useRef } from "react";
import { ResultSource } from "..";
import { useTranslation } from "react-i18next";

type Props = {
  running: boolean;
  src: ResultSource;
  error: string | null;
  current?: number;
  total?: number;
  cancel: () => void;
};

// TODO: Add translations + number format
export default ({ running, src, error, current, total, cancel }: Props) => {
  const { t } = useTranslation();
  const loadingToast = useRef<Toaster>(null);
  const key = useRef<string | undefined>(undefined);

  useEffect(() => {
    if (error != null) {
      loadingToast.current?.clear();
      return;
    }

    if (current == undefined || total == undefined) {
      key.current = loadingToast.current?.show(
        {
          message: t<string>("sim.loading"),
          icon: "refresh",
          intent: Intent.PRIMARY,
          isCloseButtonShown: false,
          timeout: 0,
        },
        key.current
      );
      return;
    }

    if (current >= total && src == ResultSource.Loaded) {
      key.current = loadingToast.current?.show(
        {
          message: t<string>("sim.loaded", { i: current }),
          icon: "tick",
          intent: Intent.SUCCESS,
          isCloseButtonShown: true,
          timeout: 2000,
        },
        key.current
      );
      return;
    }

    // TODO: bug with loading toast where it'll immediately reappear after cancel due to delayed
    //    flush from the throttled set calls. Need to find a way to have it ignore these cases
    //    or disappear on its own. This check "fixes" it but makes success timeout not correct.
    if (!running) {
      loadingToast.current?.clear();
      return;
    }

    key.current = loadingToast.current?.show(
      {
        message: (
          <ProgressToast
            cancel={cancel}
            current={current}
            total={total}
            toastKey={key}
            loadingToast={loadingToast}
          />
        ),
        className: "w-full !max-w-2xl",
        intent: Intent.NONE,
        isCloseButtonShown: current >= total,
        timeout: current < total ? 0 : 2000,
      },
      key.current
    );
  }, [current, total, src, error, running, cancel]);

  return <Toaster ref={loadingToast} position={Position.TOP} className="z-50" />;
};

const ProgressToast = ({
      cancel,
      current,
      total,
      toastKey,
      loadingToast,
    }: {
      cancel: () => void;
      current: number;
      total: number;
      toastKey: MutableRefObject<string | undefined>;
      loadingToast: RefObject<Toaster>;
    }) => {
  const { t } = useTranslation();
  const val = current / total;
  return (
    <div className="flex flex-row items-center justify-between gap-2">
      <div className="min-w-fit">
        {t<string>("sim.running")} ({current}/{total})
      </div>
      <ProgressBar
        className={classNames("basis-1/2 flex-auto sm:min-w-", {
          [Classes.PROGRESS_NO_STRIPES]: val >= 1,
        })}
        intent={val < 1 ? Intent.PRIMARY : Intent.SUCCESS}
        value={val}
      />
      {action(val, toastKey, loadingToast, cancel, t<string>("db.cancel"))}
    </div>
  );
};

function action(
  val: number,
  key: MutableRefObject<string | undefined>,
  loadingToast: RefObject<Toaster>,
  cancel: () => void,
  cancelText: string
) {
  if (val >= 1) {
    return null;
  }
  return (
    <Button
      className="!min-w-fit"
      text={cancelText}
      intent={Intent.DANGER}
      onClick={() => {
        cancel();
        loadingToast.current?.clear();
        key.current = undefined;
      }}
    />
  );
}
