import { CharacterCard } from "@gcsim/components";
import { dynamicKey } from "@gcsim/localization";
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	ButtonGroup,
	NonIdealState,
} from "@gcsim/primitives";
import type { Sample } from "@gcsim/types";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import { CopyToClipboard, SendToSimulator } from "../../Components/Buttons";
import { useSendToSimulator } from "../../Components/Buttons/useSendToSimulator";
import { characterCardsClassNames } from "../Viewer/Components/Overview/TeamHeader";
import {
	DefaultSampleOptions,
	parseLogV2,
	type SampleRow,
	Sampler,
} from "./Components";

const SAVED_SAMPLE_KEY = "gcsim-sample-settings";

type UseSampleData = {
	parsed: SampleRow[] | null;
	team?: string[];
	searchable: { [key: number]: string[] };
	settings: string[];
	setSettings: (val: string[]) => void;
};

type Props = {
	sample: Sample | null;
	error: string | null;
	retry?: () => void;
};

export default ({ sample, error, retry }: Props) => {
	const { t } = useTranslation();
	const data = useSample(sample);
	const onSendToSimulator = useSendToSimulator();

	if (sample == null || data.team == null || data.parsed == null) {
		return (
			<>
				<NonIdealState loading />
				<ErrorAlert msg={error} retry={retry} />
			</>
		);
	}

	const cardClass = characterCardsClassNames(
		sample.character_details?.length ?? 4,
	);
	return (
		<div className="flex flex-col gap-2 w-full 2xl:mx-auto 2xl:container py-6">
			<div className="flex flex-row justify-between pl-6 pr-4 pb-2">
				<span className="text-g-lg font-bold font-g-mono">
					{t("db.number_of_targets") + sample.target_details?.length}
				</span>
				<ButtonGroup>
					<CopyToClipboard
						config={sample.config}
						className="hidden ml-[7px] sm:flex"
					/>
					<SendToSimulator
						config={sample.config}
						onSendToSimulator={onSendToSimulator}
					/>
				</ButtonGroup>
			</div>
			<div className="flex flex-row gap-2 justify-center flex-wrap px-4 pb-2">
				{sample.character_details?.map((c) => (
					<CharacterCard
						key={c.name}
						char={c}
						showDetails={false}
						stats={[]}
						snapshot={[]}
						statsRows={0}
						name={t(dynamicKey(`game:character_names.${c.name}`))}
						constellationLabel={`${t("character.c_pre")}${c.cons ?? 0}${t("character.c_post")}`}
						levelLabel={t("character.lvl")}
						talentsLabel={t("character.talents")}
						artifactStatsLabel={t("character.artifact_stats")}
						totalStatsLabel={t("character.total_stats")}
						weaponName={
							c.weapon
								? t(dynamicKey(`game:weapon_names.${c.weapon.name}`))
								: ""
						}
						className={cardClass}
					/>
				))}
			</div>
			<div className="flex flex-grow flex-col gap-[15px] px-4">
				<Sampler
					sample={sample}
					data={data.parsed}
					team={data.team}
					searchable={data.searchable}
					settings={data.settings}
					setSettings={data.setSettings}
				/>
				<ErrorAlert msg={error} retry={retry} />
			</div>
		</div>
	);
};

type ErrorProps = {
	msg: string | null;
	retry?: () => void;
};

const ErrorAlert = ({ msg, retry }: ErrorProps) => {
	const { t } = useTranslation();
	const navigate = useNavigate();

	return (
		<AlertDialog open={msg != null}>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>{t("viewer.error_encountered")}</AlertDialogTitle>
					<AlertDialogDescription>{msg}</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter>
					{retry != null && (
						<AlertDialogCancel onClick={() => retry()}>
							{t("viewer.retry")}
						</AlertDialogCancel>
					)}
					<AlertDialogAction
						variant="destructive"
						onClick={() => navigate("/")}
					>
						{t("viewer.close")}
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	);
};

function useSample(sample: Sample | null): UseSampleData {
	const [selected, setSelected] = useState<string[]>(() => {
		const saved = localStorage.getItem(SAVED_SAMPLE_KEY);
		if (saved) {
			const initialValue = JSON.parse(saved);
			return initialValue || DefaultSampleOptions;
		}
		return DefaultSampleOptions;
	});

	const setAndStore = (val: string[]) => {
		setSelected(val);
		localStorage.setItem(SAVED_SAMPLE_KEY, JSON.stringify(val));
	};

	const parsed = useMemo(() => {
		if (
			sample?.initial_character == null ||
			sample?.character_details == null
		) {
			return null;
		}

		return parseLogV2(
			sample.initial_character,
			sample.character_details.map((c) => c.name),
			sample.logs,
			selected,
		);
	}, [
		sample?.character_details,
		sample?.initial_character,
		sample?.logs,
		selected,
	]);

	const searchable = useMemo(() => {
		const out: { [key: number]: string[] } = {};
		if (parsed == null) {
			return out;
		}

		parsed.forEach((row, i) => {
			const results: string[] = [];
			row.slots.forEach((slot) => {
				slot.forEach((e) => {
					results.push(e.msg);
				});
			});
			out[i] = results;
		});
		return out;
	}, [parsed]);

	return {
		parsed: parsed,
		team: sample?.character_details?.map((c) => c.name),
		searchable: searchable,
		settings: selected,
		setSettings: setAndStore,
	};
}
