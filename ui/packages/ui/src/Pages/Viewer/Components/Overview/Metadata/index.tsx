import {
	Commit,
	DateItem,
	DevBuild,
	Dirty,
	Energy,
	Iterations,
	Mode,
	Swap,
	WarningItem,
} from "@gcsim/components";
import { Card } from "@gcsim/primitives";
import type { model } from "@gcsim/types";

type Props = {
	modelData: model.SimulationResult | null;
};

export default ({ modelData }: Props) => {
	return (
		<Card className="flex flex-row flex-wrap col-span-full !p-2 gap-2 justify-center">
			<DevBuild signKey={modelData?.key_type} />
			<Dirty modified={modelData?.modified} />
			<WarningItem warnings={modelData?.statistics?.warnings} />
			<Iterations itr={modelData?.statistics?.iterations} />
			<Mode mode={modelData?.mode} />
			<DateItem date={modelData?.build_date} />
			<Commit commit={modelData?.sim_version} />
			<Swap swap={modelData?.simulator_settings?.delays?.swap} />
			<Energy energy={modelData?.energy_settings} />
		</Card>
	);
};
