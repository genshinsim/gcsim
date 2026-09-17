import { Colors } from "@gcsim/components";
import { Card } from "@gcsim/primitives";
import classNames from "classnames";
import { ChevronRight } from "lucide-react";
import { memo, type ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router";
import { CardTitle } from "../../Util";

type AuxStat = {
	title: string;
	value?: string;
};

type CardProps = {
	title: string;
	color: string;
	value?: string;
	label?: string;
	auxStats?: Array<AuxStat>;
	tooltip?: string | JSX.Element;
	hashLink?: string;
};

const CardTemplate = ({
	title,
	color,
	value,
	label,
	auxStats,
	tooltip,
	hashLink,
}: CardProps) => {
	const history = useHistory();
	// const interactable = hashLink != null;
	const interactable = false;
	const click = () => {
		history.replace({ hash: hashLink });
	};

	return (
		<div
			className="flex basis-1/4 flex-auto pl-1 min-w-fit"
			style={{ background: color }}
		>
			<Card
				className="flex flex-auto flex-row items-stretch justify-between gap-0 p-5"
				onClick={() => interactable && value !== undefined && click()}
			>
				<div className="flex flex-col justify-start">
					<CardTitle title={title} tooltip={tooltip} />
					<CardValue value={value} label={label} />
					<CardAux aux={auxStats} />
				</div>
				<CardChevron interactable={interactable} />
			</Card>
		</div>
	);
};

export const RollupCard = memo(CardTemplate);

const CardValue = ({
	value,
	label,
}: {
	value?: number | string | null;
	label?: string;
}) => {
	const { i18n } = useTranslation();

	const out = value == null ? 1234 : value;
	const valueClass = classNames("text-5xl font-bold tabular-nums", {
		"animate-pulse rounded bg-muted text-transparent": value == null,
	});

	let lbl: ReactNode;
	if (label != null) {
		lbl = (
			<div className="flex items-start text-base text-gray-400">{label}</div>
		);
	}

	return (
		<div className="flex flex-row py-2 gap-1 justify-start">
			<div className={valueClass}>{out.toLocaleString(i18n.language)}</div>
			{lbl}
		</div>
	);
};

const CardChevron = ({ interactable }: { interactable: boolean }) => {
	if (!interactable) {
		return null;
	}
	return (
		<div className="flex flex-grow justify-end self-stretch justify-self-end">
			<ChevronRight size={36} color={Colors.GRAY1} className="self-center" />
		</div>
	);
};

const CardAux = ({ aux }: { aux?: Array<AuxStat> }) => {
	if (aux == null) {
		return null;
	}

	return (
		<div className="grid grid-cols-3 gap-x-5 pt-1 justify-start text-sm font-mono min-w-fit">
			{aux.map((e) => (
				<AuxItem key={e.title} stat={e} />
			))}
		</div>
	);
};

const AuxItem = ({ stat }: { stat: AuxStat }) => {
	const { i18n } = useTranslation();

	const cls = classNames("font-black text-current text-sm text-gray-100", {
		"animate-pulse rounded bg-muted text-transparent": stat.value == null,
	});
	const val = stat.value == null ? 123.45 : stat.value;

	return (
		<div className="flex flex-row items-start gap-3">
			<div className="text-gray-400">{stat.title}</div>
			<div className={cls}>{val.toLocaleString(i18n.language)}</div>
		</div>
	);
};
