import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
	Label,
	Switch,
} from "@gcsim/primitives";
import { useTranslation } from "react-i18next";
import { usePrefs } from "../../stores/AppState";
import { type ExecutorKind, useHost } from "../../stores/host";

const EXECUTOR_PILL = {
	wasm: { fixedBg: "bg-violet-700", labelKey: "settings.executor_wasm" },
	server: { fixedBg: "bg-amber-700", labelKey: "settings.executor_server" },
} as const satisfies Record<ExecutorKind, unknown>;

export function ExecutorDialog() {
	const { t } = useTranslation();
	const { executorKind, executorSettings } = useHost();
	const { sampleOnLoad, setSampleOnLoad } = usePrefs();
	const { fixedBg, labelKey } = EXECUTOR_PILL[executorKind];

	return (
		<Dialog>
			<DialogTrigger
				className={`${fixedBg} cursor-pointer rounded-full px-2 py-0.5 font-g-mono text-g-xs font-medium text-white outline-none hover:brightness-110 focus-visible:ring-2 focus-visible:ring-g-accent`}
			>
				{t(labelKey)}
			</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>{t("settings.executor")}</DialogTitle>
					<DialogDescription>{t("settings.executor_sub")}</DialogDescription>
				</DialogHeader>
				<div className="flex flex-col gap-6">
					{executorSettings}
					<div className="flex items-center gap-2">
						<Switch
							id="sample-on-load"
							checked={sampleOnLoad}
							onCheckedChange={() => setSampleOnLoad(!sampleOnLoad)}
						/>
						<Label htmlFor="sample-on-load">
							{t("simple.generate_sample")}
						</Label>
					</div>
				</div>
			</DialogContent>
		</Dialog>
	);
}
