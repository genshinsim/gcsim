import { Card } from "@blueprintjs/core";
import { Coord, Enemy } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { memo } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle, NoData } from "../../Util";
import { EnemyCard } from "./EnemyCard";
import { PositionGraph } from "./PositionGraph";

type Props = {
  enemies?: Enemy[];
  player?: Coord;
};

const TargetInfo = (props: Props) => {
  const { t } = useTranslation();
  return (
    <Card className="flex flex-col col-span-3">
      <CardTitle title={t<string>("result.target_info")} tooltip="x" />
      <CardData {...props} />
    </Card>
  );
};

export default memo(TargetInfo);

const CardData = ({ enemies, player }: Props) => {
  const { t } = useTranslation();
  if (enemies == null || enemies.length == 0) {
    return <NoData />;
  }

  return (
    <div className="flex flex-col-reverse lg:flex-row gap-2 pt-2 h-64">
      <div className="flex flex-col gap-2 grow basis-2/3 overflow-y-scroll h-full min-w-[250px]">
        {enemies.map((enemy, i) => (
          <EnemyCard key={`enemy-${i}`} id={i} enemy={enemy} />
        ))}
      </div>
      <div className="flex flex-col grow w-[236px] min-h-[100px] lg:self-auto self-center">
        <div className="lg:flex flex-row justify-center text-gray-400 font-mono hidden">
          {t<string>("result.target_pos")}
        </div>
        <ParentSize>
          {({ width, height }) => (
            <PositionGraph
              width={width}
              height={height}
              enemies={enemies}
              player={player}
            />
          )}
        </ParentSize>
      </div>
    </div>
  );
};
