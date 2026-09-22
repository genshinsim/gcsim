import type { model, ParsedResult } from "@gcsim/types";
import { debounce } from "lodash-es";
import React from "react";
import { useExecutor } from "./ExecutorProvider";
import { toParsedTeam } from "./parsedTeam";

const VALIDATE_DEBOUNCE_MS = 200;

export interface Validation {
	isValid: boolean;
	error: string | null;
	parsedTeam: model.Character[];
}

function asError(err: unknown): string {
	return typeof err === "string" ? err : String(err);
}

export function useValidation(config: string): Validation {
	const { exec, isReady } = useExecutor();
	const [isValid, setValid] = React.useState(false);
	const [error, setError] = React.useState<string | null>(null);
	const [parsedTeam, setParsedTeam] = React.useState<model.Character[]>([]);

	const debouncedRef = React.useRef(
		debounce((fn: () => void) => fn(), VALIDATE_DEBOUNCE_MS),
	);

	const applyResult = React.useCallback((result: ParsedResult) => {
		setParsedTeam(toParsedTeam(result));
		if (result.errors && result.errors.length > 0) {
			setError(result.errors.join("\n"));
			setValid(false);
			return;
		}
		setError(null);
		setValid(true);
	}, []);

	const applyRejection = React.useCallback((err: unknown) => {
		setError(asError(err));
		setValid(false);
	}, []);

	React.useEffect(() => {
		const debounced = debouncedRef.current;
		return () => debounced.cancel();
	}, []);

	React.useEffect(() => {
		if (!isReady) {
			return;
		}
		if (config === "") {
			setParsedTeam([]);
			setError(null);
			setValid(false);
			return;
		}
		debouncedRef.current(() => {
			exec().validate(config).then(applyResult, applyRejection);
		});
	}, [exec, config, isReady, applyResult, applyRejection]);

	return { isValid, error, parsedTeam };
}
