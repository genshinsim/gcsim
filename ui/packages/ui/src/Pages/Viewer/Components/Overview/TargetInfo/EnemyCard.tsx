import { DataColorsConst } from "@gcsim/components";
import { dynamicKey } from "@gcsim/localization";
import { Card } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { isNumber } from "lodash-es";
import { useTranslation } from "react-i18next";
import {
	IconAnemo,
	IconCryo,
	IconDendro,
	IconElectro,
	IconGeo,
	IconHydro,
	IconPhysical,
	IconPyro,
} from "../../../../../Components/Icons";

type Props = {
	id: number;
	enemy?: model.Enemy;
};

export const EnemyCard = (props: Props) => {
	const bgColor = DataColorsConst.enemy(props.id);

	return (
		<div className="flex min-w-fit">
			<Card
				className="flex flex-auto flex-col gap-1 border-l-4 p-5"
				style={{ borderLeftColor: bgColor }}
			>
				<EnemyTitle {...props} />
				<EnemyInfo {...props} />
				<EnemyResistances {...props} />
			</Card>
		</div>
	);
};

const EnemyTitle = ({ id, enemy }: Props) => {
	const { t } = useTranslation();
	let name = enemy?.name;
	if (name) {
		name = `(${t(dynamicKey("game:monster_names." + name))})`;
	}

	return (
		<div className="flex flex-row items-end gap-3">
			<div
				className="text-g-ink-mute text-g-lg"
				style={{ color: DataColorsConst.qualitative5(id) }}
			>
				{t("viewer.target")} {id + 1} {name}
			</div>
		</div>
	);
};

const EnemyInfo = ({ enemy }: Props) => {
	const { t } = useTranslation();
	const modified = enemy?.modified ?? false;
	return (
		<div className="flex flex-row font-g-mono gap-3 h-full items-center">
			<InfoItem name={t("character.lvl")} value={enemy?.level} />
			<InfoItem name={t("stats.hp")} value={enemy?.hp} />
			<InfoItem
				name={t("stats.modified")}
				value={t(dynamicKey("states." + modified.toString()))}
			/>
		</div>
	);
};

const InfoItem = ({
	name,
	value,
}: {
	name: string;
	value?: number | string;
}) => {
	const { i18n } = useTranslation();

	if (value == null) {
		return null;
	}
	if (isNumber(value)) {
		value = value.toLocaleString(i18n.language);
	}

	return (
		<div className="flex flex-row gap-1 text-g-xs items-center">
			<div className="text-g-ink-mute">{name}</div>
			<div className="font-black text-g-sm text-g-ink">{value}</div>
		</div>
	);
};

const EnemyResistances = ({ enemy }: Props) => {
	return (
		<div className="grid grid-cols-4 gap-y-1 text-g-sm font-g-mono">
			<Resistance type="anemo" num={enemy?.resist?.["anemo"]} />
			<Resistance type="geo" num={enemy?.resist?.["geo"]} />
			<Resistance type="electro" num={enemy?.resist?.["electro"]} />
			<Resistance type="hydro" num={enemy?.resist?.["hydro"]} />
			<Resistance type="pyro" num={enemy?.resist?.["pyro"]} />
			<Resistance type="cryo" num={enemy?.resist?.["cryo"]} />
			<Resistance type="dendro" num={enemy?.resist?.["dendro"]} />
			<Resistance type="physical" num={enemy?.resist?.["physical"]} />
		</div>
	);
};

const Resistance = ({ type, num }: { type: string; num?: number }) => {
	const { i18n } = useTranslation();
	const format = (val?: number) =>
		val?.toLocaleString(i18n.language, {
			maximumFractionDigits: 2,
			style: "percent",
		});

	return (
		<div className="flex flex-row gap-2 items-center">
			<Icon type={type} />
			<div>{format(num ?? 0)}</div>
		</div>
	);
};

const Icon = ({ type }: { type: string }) => {
	const size = "w-[16px] h-[16px] min-w-[16px] min-h-[16px]";
	switch (type) {
		case "electro":
			return <IconElectro className={`${size} text-g-electro`} />;
		case "pyro":
			return <IconPyro className={`${size} text-g-pyro`} />;
		case "cryo":
			return <IconCryo className={`${size} text-g-cryo`} />;
		case "hydro":
			return <IconHydro className={`${size} text-g-hydro`} />;
		case "geo":
			return <IconGeo className={`${size} text-g-geo`} />;
		case "anemo":
			return <IconAnemo className={`${size} text-g-anemo`} />;
		case "physical":
			return <IconPhysical className={`${size}`} />;
		case "dendro":
			return <IconDendro className={`${size} text-g-dendro`} />;
		default:
			return <IconHydro className={`${size} text-g-hydro`} />;
	}
};
