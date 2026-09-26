import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import { debounce } from "lodash-es";
import { useEffect, useRef, useState } from "react";

// Detects changes in the config and validates them with the executor.
// validated == true means a validation check ran to completion, not that the
// config is valid. Still used by the viewer's Config tab.
export function useConfigValidateListener(
	exec: ExecutorSupplier<Executor>,
	cfg: string,
	isReady: boolean | null,
	setErr: (str: string) => void,
): boolean {
	const [validated, setValidated] = useState(false);
	const debounced = useRef(debounce((x: () => void) => x(), 200));

	useEffect(() => {
		if (!isReady || cfg === "") {
			return;
		}

		setValidated(false);
		debounced.current(() => {
			exec()
				.validate(cfg)
				.then(
					(res) => {
						setErr("");
						if (res.errors) {
							let msg = "";
							res.errors.forEach((err) => {
								msg += err + "\n";
							});
							setErr(msg);
						}
						setValidated(true);
					},
					(err) => {
						setErr(err);
						setValidated(false);
					},
				);
		});
	}, [exec, cfg, setErr, isReady]);

	return validated;
}
