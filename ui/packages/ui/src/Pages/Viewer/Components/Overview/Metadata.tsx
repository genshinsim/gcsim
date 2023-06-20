import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { CardTitle, NoData } from "../Util";

type Props = {
  data: SimResults | null;
}

export default ({}: Props) => {
  return (
    <Card className="flex flex-col col-span-full h-12 min-h-full">
    </Card>
  );
};