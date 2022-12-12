import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { CardTitle } from "../Util";

type Props = {
  data: SimResults | null;
}

export default ({}: Props) => {
  return (
    <Card className="flex flex-col col-span-3 h-24 min-h-full">
      <CardTitle title="Target Information" tooltip="x" />
    </Card>
  );
};