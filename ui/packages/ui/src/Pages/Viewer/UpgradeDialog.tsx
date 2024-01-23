import {
  Button,
  Callout,
  Classes,
  Dialog,
  Divider,
  Intent,
} from "@blueprintjs/core";
import { Executor, ExecutorSupplier } from "@gcsim/executors";
import { SimResults, Version } from "@gcsim/types";
import classNames from "classnames";
import { useEffect, useState } from "react";
import { useHistory, useLocation } from "react-router";
import ExecutorSettingsButton from "../../Components/Buttons/ExecutorSettingsButton";

// THIS MUST ALWAYS BE IN SYNC WITH THE GCSIM BINARY
const MAJOR = "4"; // Make sure the gcsim binary has also been updated
const MINOR = "2"; // Make sure the gcsim binary has also been updated

enum MismatchType {
  MajorVersionMismatch,
  MinorVersionMismatch,
  CommitMismatch,
  NoMismatch,
}

type Props = {
  exec: ExecutorSupplier<Executor>;
  data: SimResults | null;
  redirect: string;
  mode: string;
  commit: string;
  setResult: (result: SimResults | null, hash: string | null) => void;
  setError: (recoveryConfig: string | null, err: string | null) => void;
};

// TODO: translations
export default ({ exec, data, redirect, mode, commit, setResult, setError }: Props) => {
  const mismatch = useMismatch(data?.sim_version, commit, data?.schema_version);
  const [isOpen, setOpen] = useState(true);
  const location = useLocation();

  if (data == null || mismatch == MismatchType.NoMismatch) {
    return null;
  }

  // only show hash mismatch on share links to reduce noise (for now)
  if (mismatch == MismatchType.CommitMismatch
      && !location.pathname.startsWith("/sh/") && !location.pathname.startsWith("/db/")) {
    return null;
  }

  // only show major version errors in development
  if (mismatch != MismatchType.MajorVersionMismatch && mode === "development") {
    return null;
  }

  const minor = mismatch == MismatchType.CommitMismatch
      || mismatch == MismatchType.MinorVersionMismatch;

  return (
    <Dialog
        isOpen={isOpen}
        title="Results Outdated"
        icon="outdated"
        usePortal={false}
        canEscapeKeyClose={minor}
        canOutsideClickClose={minor}
        isCloseButtonShown={minor}
        onClose={() => setOpen(false)}>
      <div className={Classes.DIALOG_BODY}>
        <DialogBody mismatch={mismatch} data={data} latestCommit={commit} />
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

function useMismatch(
  resultCommit?: string, latestCommit?: string, schema_version?: Version): MismatchType | null {
  const [mismatch, setMismatch] = useState<MismatchType | null>(null);

  useEffect(() => {
    if (schema_version == null) {
      setMismatch(MismatchType.MajorVersionMismatch);
    } else if (schema_version.major != MAJOR) {
      setMismatch(MismatchType.MajorVersionMismatch);
    } else if (schema_version.minor != MINOR) {
      setMismatch(MismatchType.MinorVersionMismatch);
    } else if (resultCommit != latestCommit) {
      setMismatch(MismatchType.CommitMismatch);
    } else {
      setMismatch(MismatchType.NoMismatch);
    }
  }, [schema_version, resultCommit, latestCommit]);

  return mismatch;
}

type BodyProps = {
  mismatch: MismatchType | null;
  data: SimResults | null;
  latestCommit?: string;
};

const DialogBody = ({ mismatch, data, latestCommit }: BodyProps) => {
  const simCommit = data?.sim_version ?? data?.version;

  const shortResultCommit = simCommit?.substring(0, 7);
  const shortLatestCommit = latestCommit?.substring(0, 7);
  const resultCommitUrl = "https://github.com/genshinsim/gcsim/commits/" + simCommit;
  const latestCommitUrl = "https://github.com/genshinsim/gcsim/commits/" + latestCommit;
  const diffUrl = (
    "https://github.com/genshinsim/gcsim/compare/" + simCommit + "..." + latestCommit
  );

  const dirty = data?.modified || simCommit === "";

  const major = data?.schema_version?.major;
  const minor = data?.schema_version?.minor;

  const versionClass = classNames(
    "inline-grid grid-cols-[repeat(6,_max-content)] justify-start gap-y-0 gap-x-3",
    "text-xs pt-2 font-mono text-gray-400"
  );

  const VersionInfo = ({}) => (
    <div className={versionClass}>
      {/* version line */}
      <div>version</div>
      <div>{major == null || minor == null ? "legacy" : `${major}.${minor}`}</div>
      <Divider className="h-full" />
      <div>latest</div>
      <div>{MAJOR}.{MINOR}{" "}</div>
      <div></div>
      
      {/* commit line */}
      <div className="justify-self-end">commit</div>
      <a href={resultCommitUrl} target="_blank" rel="noreferrer">
        {shortResultCommit}
      </a>
      <Divider className={classNames({ ["h-full"]: dirty })} />
      <div>latest</div>
      <a href={latestCommitUrl} target="_blank" rel="noreferrer">
        {shortLatestCommit}
      </a>
      <a href={diffUrl} target="_blank" rel="noreferrer">
        (diff)
      </a>

      {/* dirty line */}
      {dirty && (
        <>
          <div className="justify-self-end">dirty?</div>
          <div className="text-red-500">true</div>
          <Divider />
        </>
      )}
    </div>
  );

  if (mismatch == MismatchType.CommitMismatch) {
    return (
      <Callout title="Commit Hash Mismatch" intent={Intent.WARNING}>
        <div>
          This simulation was generated with an outdated commit. Some data, graphs, or features
          may be missing or inaccurate. Upgrade to resolve compatibility issues.
        </div>
        <VersionInfo />
      </Callout>
    );
  }

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
      exec: ExecutorSupplier<Executor>,
      cfg?: string;
      setResult: (result: SimResults | null, hash: string | null) => void;
      setError: (config: string | null, err: string | null) => void;
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

    setResult(null, null);
    setError(null, null);
    exec().run(cfg, setResult).catch((err) => {
      setError(cfg, err);
    });
  };

  return <Button text="Upgrade" intent={Intent.SUCCESS} loading={!isReady} onClick={run} />;
};

const CancelButton = ({ mismatch, setOpen, redirect }: {
      mismatch: MismatchType | null;
      setOpen: (open: boolean) => void;
      redirect: string;
    }) => {
  const history = useHistory();

  if (mismatch == MismatchType.MajorVersionMismatch) {
    return <Button text="Cancel" intent={Intent.DANGER} onClick={() => history.push(redirect)} />;
  }
  return <Button text="Dismiss" onClick={() => setOpen(false)} />;
};
