import type { Character } from "@gcsim/types";

export interface TeamViewProps {
	parsedTeam: Character[];
	isValid: boolean;
	error: string | null;
}

export function TeamView({ parsedTeam, isValid, error }: TeamViewProps) {
	return (
		<div data-testid="editor-team-view" data-valid={isValid}>
			{error ?? `${parsedTeam.length} characters`}
		</div>
	);
}
