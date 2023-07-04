import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { DevBuild } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/DevBuild";
import { Dirty } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Dirty";
import { Item } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Item";
import { Iterations } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Iterations";
import { ModeItem } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Mode";
import { Standard } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Standard";
import { WarningItem } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Warning";

type Props = {
  data: SimResults;
}

export const Metadata = ({ data }: Props) => {
  if (data.schema_version == null) {
    return (
      <Card className="flex flex-row flex-wrap !p-2 gap-2 justify-center">
        <Item value="legacy sim" intent="danger" bold bright />
      </Card>
    );
  }

  return (
    <Card className="flex flex-row flex-wrap !p-2 gap-2 justify-center">
      <Error signKey={data.key_type} modified={data.modified} />
      <WarningItem warnings={data?.statistics?.warnings} />
      <Standard standard={data?.standard} />
      <Iterations itr={data?.statistics?.iterations} />
      <ModeItem mode={data?.mode} />
    </Card>
  );
};

type ErrorProps = {
  signKey?: string;
  modified?: boolean;
}

export const Error = ({ signKey, modified }: ErrorProps) => {
  if (signKey == null || signKey == "prod") {
    return <Dirty modified={modified} />;
  }

  return <DevBuild signKey={signKey} />;
};