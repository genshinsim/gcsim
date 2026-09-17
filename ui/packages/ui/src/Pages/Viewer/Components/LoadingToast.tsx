import { Button, Progress, toast } from "@gcsim/primitives";
import { useEffect, useRef } from "react";
import { useTranslation } from "react-i18next";
import { ResultSource } from "..";

type Props = {
	running: boolean;
	src: ResultSource;
	error: string | null;
	current?: number;
	total?: number;
	cancel: () => void;
};

// TODO: Add translations + number format
export default ({ running, src, error, current, total, cancel }: Props) => {
	const { t } = useTranslation();
	const toastId = useRef<string | number | undefined>(undefined);

	useEffect(() => {
		const dismiss = () => {
			if (toastId.current !== undefined) {
				toast.dismiss(toastId.current);
				toastId.current = undefined;
			}
		};

		if (error != null) {
			dismiss();
			return;
		}

		if (current === undefined || total === undefined) {
			toastId.current = toast.loading(t("sim.loading"), {
				id: toastId.current,
				duration: Number.POSITIVE_INFINITY,
			});
			return;
		}

		if (current >= total && src === ResultSource.Loaded) {
			toastId.current = toast.success(t("sim.loaded", { i: current }), {
				id: toastId.current,
				duration: 2000,
				closeButton: true,
			});
			return;
		}

		// TODO: bug with loading toast where it'll immediately reappear after cancel due to delayed
		//    flush from the throttled set calls. Need to find a way to have it ignore these cases
		//    or disappear on its own. This check "fixes" it but makes success timeout not correct.
		if (!running) {
			dismiss();
			return;
		}

		toastId.current = toast(
			<ProgressToast
				cancel={() => {
					cancel();
					dismiss();
				}}
				current={current}
				total={total}
			/>,
			{
				id: toastId.current,
				duration: current < total ? Number.POSITIVE_INFINITY : 2000,
				closeButton: current >= total,
				className: "w-full",
				style: { width: "min(90vw, 42rem)" },
			},
		);
	}, [current, total, src, error, running, cancel, t]);

	return null;
};

const ProgressToast = ({
	cancel,
	current,
	total,
}: {
	cancel: () => void;
	current: number;
	total: number;
}) => {
	const { t } = useTranslation();
	const val = current / total;
	return (
		<div className="flex flex-row items-center justify-between gap-2 w-full">
			<div className="min-w-fit">
				{t("sim.running")} ({current}/{total})
			</div>
			<Progress
				className="basis-1/2 flex-auto"
				value={Math.min(val, 1) * 100}
			/>
			{val < 1 ? (
				<Button className="!min-w-fit" variant="destructive" onClick={cancel}>
					{t("db.cancel")}
				</Button>
			) : null}
		</div>
	);
};
