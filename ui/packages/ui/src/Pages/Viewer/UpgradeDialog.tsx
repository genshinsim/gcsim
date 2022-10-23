import {
  Button,
  Callout,
  Classes,
  Dialog,
  Intent,
} from "@blueprintjs/core";
import { ExecutorSupplier } from "@gcsim/executors";
import { SimResults, Version } from "@gcsim/types";
import { useEffect, useState } from "react";
import { useLocation } from "wouter";
import ExecutorSettingsButton from "../../Components/ExecutorSettingsButton";

// THIS MUST ALWAYS BE IN SYNC WITH THE GCSIM BINARY
const MAJOR = 4; // Make sure the gcsim binary has also been updated
const MINOR = 0; // Make sure the gcsim binary has also been updated

enum MismatchType {
  MajorVersionMismatch,
  MinorVersionMismatch,
  NoMismatch,
}

type Props = {
  exec: ExecutorSupplier;
  data: SimResults | null;
  redirect: string;
  setResult: (result: SimResults | null) => void;
  setError: (err: string | null) => void;
};

// TODO: translations
export default ({ exec, data, redirect, setResult, setError }: Props) => {
  const mismatch = useMismatch(data?.schema_version);
  const [isOpen, setOpen] = useState(true);

  if (data == null || mismatch == MismatchType.NoMismatch) {
    return null;
  }

  return (
    <Dialog
      isOpen={isOpen}
      title="Results Outdated"
      icon="outdated"
      usePortal={false}
      canEscapeKeyClose={false}
      canOutsideClickClose={false}
      isCloseButtonShown={mismatch == MismatchType.MinorVersionMismatch}
      onClose={() => setOpen(false)}
    >
      <div className={Classes.DIALOG_BODY}>
        <DialogBody
          mismatch={mismatch}
          major={data.schema_version?.major}
          minor={data.schema_version?.minor}
        />
      </div>
      <div className="flex justify-between items-end gap-16 mx-4">
        <div className="max-w-[196px] min-w-[120px] flex-auto">
          <ExecutorSettingsButton />
        </div>
        <div className="flex justify-end gap-[10px]">
          <UpgradeButton
              exec={exec} cfg={data.config_file} setResult={setResult} setError={setError} />
          <CancelButton
              mismatch={mismatch}
              setOpen={setOpen}
              redirect={redirect} />
        </div>
      </div>
    </Dialog>
  );
};

function useMismatch(schema_version?: Version): MismatchType | null {
  const [mismatch, setMismatch] = useState<MismatchType | null>(null);

  useEffect(() => {
    if (schema_version == null) {
      setMismatch(MismatchType.MajorVersionMismatch);
    } else if (schema_version.major != MAJOR) {
      setMismatch(MismatchType.MajorVersionMismatch);
    } else if (schema_version.minor < MINOR) {
      setMismatch(MismatchType.MinorVersionMismatch);
    } else {
      setMismatch(MismatchType.NoMismatch);
    }
  }, [schema_version]);

  return mismatch;
}

const DialogBody = ({
  mismatch,
  major,
  minor,
}: {
  mismatch: MismatchType | null;
  major?: number;
  minor?: number;
}) => {
  const VersionInfo = ({}) => (
    <div className="flex justify-start gap-2 text-xs pt-2 font-mono text-gray-400">
      <div>version: {major == null || minor == null ? "legacy" : `${major}.${minor}`}</div>
      <div>|</div>
      <div>
        latest: {MAJOR}.{MINOR}{" "}
      </div>
    </div>
  );

  if (mismatch == MismatchType.MinorVersionMismatch) {
    return (
      <Callout title="Minor Version Mismatch" intent={Intent.WARNING}>
        <div>
          This simulation was generated with outdated results. Some data, graphs, or features may be
          missing or inaccurate. Upgrade to resolve compatibility issues.
        </div>
        <VersionInfo />
      </Callout>
    );
  }
  return (
    <Callout title="Major Version Mismatch" intent={Intent.DANGER}>
      <div>
        Simulation results are incompatible with latest version of gcsim. Upgrade will attempt to
        resimulate and generate new results.
      </div>
      <VersionInfo />
    </Callout>
  );
};

const UpgradeButton = ({
      exec,
      cfg,
      setResult,
      setError,
    }: {
      exec: ExecutorSupplier,
      cfg?: string;
      setResult: (result: SimResults | null) => void;
      setError: (err: string | null) => void;
    }) => {
  const [isReady, setReady] = useState(false);
  useEffect(() => {
    const interval = setInterval(() => {
      setReady(exec().ready());
    }, 250);
    return () => clearInterval(interval);
  }, [exec]);

  const run = () => {
    if (cfg == null) {
      return;
    }

    setResult(null);
    setError(null);
    exec().run(cfg, (result) => {
      setResult(result);
    }).catch((err) => {
      setError(err);
    });
  };

  return <Button text="Upgrade" intent={Intent.SUCCESS} loading={!isReady} onClick={run} />;
};

const CancelButton = ({ mismatch, setOpen, redirect }: {
      mismatch: MismatchType | null;
      setOpen: (open: boolean) => void;
      redirect: string;
    }) => {
  const [, setLocation] = useLocation();

  if (mismatch == MismatchType.MinorVersionMismatch) {
    return <Button text="Dismiss" onClick={() => setOpen(false)} />;
  }
  return <Button text="Cancel" intent={Intent.DANGER} onClick={() => setLocation(redirect)} />;
};
