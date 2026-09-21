import { CardTitle, NoData, PositionGraph } from "@gcsim/components";
import { Card } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { memo } from "react";
import { useTranslation } from "react-i18next";
import { EnemyCard } from "./EnemyCard";

type Props = {
	enemies?: model.Enemy[];
	player?: model.Coord;
};

const TargetInfo = (props: Props) => {
	const { t } = useTranslation();
	return (
		<Card className="flex flex-col col-span-3 gap-0 p-5">
			<CardTitle title={t("result.target_info")} />
			<CardData {...props} />
		</Card>
	);
};

export default memo(TargetInfo);

const CardData = ({ enemies, player }: Props) => {
	const { t } = useTranslation();
	if (enemies == null || enemies.length === 0) {
		return <NoData />;
	}

	return (
		<div className="flex flex-col-reverse lg:flex-row gap-2 pt-2 h-64">
			<div className="flex flex-col gap-2 grow basis-2/3 overflow-y-scroll scrollbar-surface h-full min-w-[250px]">
				{enemies.map((enemy, i) => (
					// biome-ignore lint/suspicious/noArrayIndexKey: name may be absent/duplicated so index completes the composite; enemy order is stable (built once in Go, never reordered)
					<EnemyCard key={`enemy-${enemy.name}-${i}`} id={i} enemy={enemy} />
				))}
			</div>
			<div className="flex flex-col grow w-[236px] min-h-[100px] lg:self-auto self-center">
				<div className="lg:flex flex-row justify-center text-g-ink-mute font-g-mono hidden">
					{t("result.target_pos")}
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
