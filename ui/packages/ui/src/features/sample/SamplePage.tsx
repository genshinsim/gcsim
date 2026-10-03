import { SampleLog } from "@gcsim/components";
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
import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { CopyToClipboard, SendToSimulator } from "../../components/buttons";
import { useSendToSimulator } from "../../components/buttons/useSendToSimulator";
import { downloadSample } from "./downloadSample";

type Props = {
	sample: Sample | null;
	error: string | null;
	retry?: () => void;
};

export default ({ sample, error, retry }: Props) => {
	const { t } = useTranslation();
	const onSendToSimulator = useSendToSimulator();

	if (sample?.initial_character == null || sample.character_details == null) {
		return (
			<>
				<NonIdealState loading />
				<ErrorAlert msg={error} retry={retry} />
			</>
		);
	}

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
			<div className="flex flex-grow flex-col gap-[15px] px-4">
				<SampleLog sample={sample} onDownload={downloadSample} />
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
						onClick={() => navigate({ to: "/" })}
					>
						{t("viewer.close")}
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	);
};
