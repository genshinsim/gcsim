import {
	defaultWorkerCount,
	type ExecutorSupplier,
	WasmExecutor,
} from "@gcsim/executors";
import { Field, FieldLabel, NumberInput } from "@gcsim/primitives";
import { UI } from "@gcsim/ui";
import { useLocalStorage } from "@gcsim/utils";
import { type ReactNode, useRef } from "react";
import { useTranslation } from "react-i18next";

const minWorkers = 1;
const maxWorkers = 30;

let exec: WasmExecutor | undefined;

function wasmLocation() {
	if (import.meta.env.PROD) {
		return (
			"/api/wasm/" +
			import.meta.env.VITE_BRANCH +
			"/" +
			import.meta.env.VITE_GIT_COMMIT_HASH +
			"/" +
			"main.wasm"
		);
	}
	return "/main.wasm";
}

const WasmMode = ({ children }: { children: ReactNode }) => {
	const { t } = useTranslation();
	const [workers, setWorkers] = useLocalStorage<number>(
		"wasm-num-workers",
		defaultWorkerCount(),
	);

	const supplier = useRef<ExecutorSupplier<WasmExecutor>>(() => {
		if (exec == null) {
			exec = new WasmExecutor(wasmLocation());
			exec.setWorkerCount(workers);
		}
		return exec;
	});

	const updateWorkers = (num: number) => {
		num = Math.min(Math.max(num, minWorkers), maxWorkers);
		setWorkers(num);
		supplier.current().setWorkerCount(num);
	};

	return (
		<UI
			exec={supplier.current}
			gitCommit={import.meta.env.VITE_GIT_COMMIT_HASH}
			mode={import.meta.env.MODE}
			executorKind="wasm"
			executorSettings={
				<>
					{children}
					<Field>
						<FieldLabel htmlFor="wasm-workers">
							{t("simple.workers")}
						</FieldLabel>
						<NumberInput
							id="wasm-workers"
							value={workers}
							onValueChange={updateWorkers}
							min={minWorkers}
							max={maxWorkers}
						/>
					</Field>
				</>
			}
		/>
	);
};

export default WasmMode;
