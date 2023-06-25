import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { Commit } from "./Commit";
import { DateItem } from "./Date";
import { DevBuild } from "./DevBuild";
import { Dirty } from "./Dirty";
import { Energy } from "./Energy";
import { Iterations } from "./Iterations";
import { ModeItem } from "./Mode";
import { Standard } from "./Standard";
import { Swap } from "./Swap";
import { WarningItem } from "./Warning";

type Props = {
  data: SimResults | null;
}

export default ({ data }: Props) => {
  return (
    <Card className="flex flex-row flex-wrap col-span-full p-2 gap-2 justify-center">
      <DevBuild signKey={data?.key_type} />
      {/* <Dirty modified={data?.modified} /> */}
      <WarningItem warnings={data?.statistics?.warnings} />
      <Standard standard={data?.standard} />
      <Iterations itr={data?.statistics?.iterations} />
      <ModeItem mode={data?.mode} />
      <DateItem date={data?.build_date} />
      <Commit commit={data?.sim_version} />
      <Swap swap={data?.simulator_settings?.delays?.swap} />
      <Energy energy={data?.energy_settings} />
    </Card>
  );
};
