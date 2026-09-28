import {
	CharacterCard,
	characterCardsClassNames,
	SampleLog,
} from "@gcsim/components";
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
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import { CopyToClipboard, SendToSimulator } from "../../Components/Buttons";
import { useSendToSimulator } from "../../Components/Buttons/useSendToSimulator";
import { downloadSample } from "./downloadSample";
import { useSample } from "./useSample";

type Props = {
	sample: Sample | null;
	error: string | null;
	retry?: () => void;
};

export default ({ sample, error, retry }: Props) => {
	const { t } = useTranslation();
	const { settings, setSettings } = useSample();
	const onSendToSimulator = useSendToSimulator();

	if (sample?.initial_character == null || sample.character_details == null) {
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
				<SampleLog
					sample={sample}
					settings={settings}
					onSettingsChange={setSettings}
					onDownload={downloadSample}
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
