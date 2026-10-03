import { Editor } from "@gcsim/components";
import { Button, NonIdealState, Spinner } from "@gcsim/primitives";
import { Play } from "lucide-react";
import { useTranslation } from "react-i18next";

export type ViewerEditorState = {
	config: string;
	setConfig: (v: string) => void;
	error: string | null;
	canRun: boolean;
	busy: boolean;
};

type Props = ViewerEditorState & {
	loading: boolean;
	onRerun: () => void;
};

export default ({
	config,
	setConfig,
	error,
	canRun,
	busy,
	loading,
	onRerun,
}: Props) => {
	const { t } = useTranslation();
	if (loading) {
		return <NonIdealState loading />;
	}
	return (
		<div className="w-full 2xl:mx-auto 2xl:container px-2 flex flex-col gap-2">
			<div className="flex flex-row flex-wrap items-center gap-1">
				<Button className="flex-1" onClick={onRerun} disabled={!canRun}>
					{busy ? <Spinner /> : <Play />}
					{t("viewer.rerun")}
				</Button>
			</div>
			<Editor
				value={config}
				onChange={setConfig}
				error={error}
				maxLines={Infinity}
			/>
		</div>
	);
};
