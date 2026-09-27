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
	data: model.SimulationResult | null;
};

export default ({ data }: Props) => {
	return (
		<Card className="flex flex-row flex-wrap col-span-full !p-2 gap-2 justify-center">
			<DevBuild signKey={data?.key_type} />
			<Dirty modified={data?.modified} />
			<WarningItem warnings={data?.statistics?.warnings} />
			<Iterations itr={data?.statistics?.iterations} />
			<Mode mode={data?.mode} />
			<DateItem date={data?.build_date} />
			<Commit commit={data?.sim_version} />
			<Swap swap={data?.simulator_settings?.delays?.swap} />
			<Energy energy={data?.energy_settings} />
		</Card>
	);
};
