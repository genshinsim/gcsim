import { namedSeeds, SampleLog, SeedPicker } from "@gcsim/components";
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	NonIdealState,
} from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { downloadSample } from "../../sample/downloadSample";
import type { SampleState } from "../../sample/useSample";

type Props = {
	data: model.SimulationResult | null;
	sample: SampleState;
	running: boolean;
};

export default ({ data, sample, running }: Props) => {
	const { t } = useTranslation();
	const [picking, setPicking] = useState(false);

	if (
		data?.character_details == null ||
		data.config_file == null ||
		sample.generating
	) {
		return <NonIdealState loading />;
	}

	return (
		<div className="w-full 2xl:mx-auto 2xl:container flex flex-grow flex-col px-2">
			<SampleLog
				sample={sample.sample}
				onDownload={downloadSample}
				onGenerate={() => setPicking(true)}
			/>
			<Dialog open={picking} onOpenChange={setPicking}>
				<DialogContent className="sm:max-w-sm" aria-describedby={undefined}>
					<DialogHeader>
						<DialogTitle>{t("viewer.generate")}</DialogTitle>
					</DialogHeader>
					<SeedPicker
						seeds={namedSeeds(data)}
						value={sample.seed}
						running={running}
						onPick={(seed) => {
							setPicking(false);
							sample.generate(seed);
						}}
					/>
				</DialogContent>
			</Dialog>
		</div>
	);
};
