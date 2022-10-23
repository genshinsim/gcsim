import { Button, ButtonGroup, Classes, Dialog } from "@blueprintjs/core";
import React from "react";
import { NumberInput } from "../../../Components";
import { RootState, useAppDispatch, useAppSelector } from "../../../Stores/store";
import { Trans, useTranslation } from "react-i18next";
import { setTotalWorkers } from "../../../Stores/appSlice";
import { ExecutorSupplier } from "@gcsim/executors";

type Props = {
  exec: ExecutorSupplier;
  isOpen: boolean;
  onClose: () => void;
};

export function SimWorkerOptions(props: Props) {
  const { t } = useTranslation();

  const { workers } = useAppSelector((state: RootState) => {
    return {
      workers: state.app.workers,
    };
  });
  const dispatch = useAppDispatch();
  const [next, setNext] = React.useState<number>(workers);

  const updateWorkers = () => {
    dispatch(setTotalWorkers(props.exec, next));
    props.onClose();
  };

  return (
    <Dialog isOpen={props.isOpen} onClose={props.onClose}>
      <div className="w-full flex flex-col p-4">
        <NumberInput
          label={`${t("components.currently_loaded_workers")}${workers}`}
          onChange={(v) => setNext(v)}
          value={next}
          min={1}
          max={30}
          integerOnly
        />
      </div>
      <div className={Classes.DIALOG_FOOTER_ACTIONS}>
        <ButtonGroup fill>
          <Button onClick={updateWorkers} className="mt-4" intent="primary">
            <Trans>components.set</Trans>
          </Button>
          <Button onClick={props.onClose} className="mt-4" intent="danger">
            Cancel
          </Button>
        </ButtonGroup>
      </div>
    </Dialog>
  );
}
