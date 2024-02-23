import { FormGroup } from "@blueprintjs/core";
import { ExecutorSupplier, ServerExecutor } from "@gcsim/executors";
import { UI } from "@gcsim/ui";
import { ReactNode, useRef } from "react";
import { useTranslation } from "react-i18next";

let exec: ServerExecutor | undefined;

const ServerMode = ({ children }: { children: ReactNode }) => {
  const { t } = useTranslation();

  const supplier = useRef<ExecutorSupplier<ServerExecutor>>(() => {
    if (exec == null) {
      exec = new ServerExecutor("http://127.0.0.1:54321");
    }
    return exec;
  });

  return (
    <UI
      exec={supplier.current}
      gitCommit={import.meta.env.VITE_GIT_COMMIT_HASH}
      mode={import.meta.env.MODE}
    >
      <FormGroup className="!m-0" label={t<string>("simple.workers")}>
        {children}
        <span>hi</span>
      </FormGroup>
    </UI>
  );
};

export default ServerMode;
