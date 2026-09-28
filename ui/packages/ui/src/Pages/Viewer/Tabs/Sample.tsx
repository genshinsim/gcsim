import { namedSeeds, SampleLog, SeedPicker } from "@gcsim/components";
import { NonIdealState } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { FlaskConical } from "lucide-react";
import { downloadSample } from "../../Sample/downloadSample";
import type { SampleState } from "../../Sample/useSample";

type Props = {
	data: model.SimulationResult | null;
	sample: SampleState;
	running: boolean;
};

export default ({ data, sample, running }: Props) => {
	if (
		data?.character_details == null ||
		data.config_file == null ||
		sample.generating
	) {
		return <NonIdealState loading />;
	}

	const picker = (
		<SeedPicker
			seeds={namedSeeds(data)}
			value={sample.seed}
			onPick={sample.generate}
			running={running}
		/>
	);

	if (sample.sample == null) {
		return (
			<NonIdealState
				icon={<FlaskConical />}
				action={picker}
				className="!px-2"
			/>
		);
	}

	return (
		<div className="w-full 2xl:mx-auto 2xl:container flex flex-grow flex-col gap-[15px] px-2">
			{picker}
			<SampleLog
				sample={sample.sample}
				settings={sample.settings}
				onSettingsChange={sample.setSettings}
				onDownload={downloadSample}
			/>
		</div>
	);
};
