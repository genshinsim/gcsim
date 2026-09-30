import {
	Alert,
	AlertDescription,
	Button,
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@gcsim/primitives";
import { getRouteApi } from "@tanstack/react-router";
import React from "react";
import { useTranslation } from "react-i18next";
import { decodeSharedConfig, type SharedConfig } from "./sharedConfig";

const routeApi = getRouteApi("/simulator");

type DecodeResult = { ok: true; shared: SharedConfig } | { ok: false };

function tryDecodeSharedConfig(encoded: string): DecodeResult {
	try {
		return { ok: true, shared: decodeSharedConfig(encoded) };
	} catch {
		return { ok: false };
	}
}

function useConsumeCfgParamOnce() {
	const { cfg: encoded } = routeApi.useSearch();
	const navigate = routeApi.useNavigate();
	const [result, setResult] = React.useState<DecodeResult | null>(null);

	React.useEffect(() => {
		if (!encoded) return;
		setResult(tryDecodeSharedConfig(encoded));
		void navigate({ to: "/simulator", search: {}, replace: true });
	}, [encoded, navigate]);

	return [result, () => setResult(null)] as const;
}

export function SharedConfigDialog({
	onLoad,
}: {
	onLoad: (cfg: string) => void;
}) {
	const { t } = useTranslation();
	const [pending, close] = useConsumeCfgParamOnce();

	if (pending == null) return null;

	return (
		<Dialog open onOpenChange={(open) => !open && close()}>
			<DialogContent className="max-h-[calc(100dvh-2rem)] overflow-y-auto">
				<DialogHeader>
					<DialogTitle className="pr-6">
						{t("simple.shared_config_heading")}
					</DialogTitle>
				</DialogHeader>
				{pending.ok ? (
					<>
						<dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1">
							<dt className="text-g-ink-mute">
								{t("simple.shared_config_title")}
							</dt>
							<dd className="min-w-0 [overflow-wrap:anywhere]">
								{pending.shared.title ?? t("simple.shared_config_unknown")}
							</dd>
							<dt className="text-g-ink-mute">
								{t("simple.shared_config_source")}
							</dt>
							<dd className="min-w-0 [overflow-wrap:anywhere]">
								{pending.shared.source ?? t("simple.shared_config_unknown")}
							</dd>
						</dl>
						<Alert variant="warning">
							<AlertDescription>
								{t("viewer.this_will_overwrite")}
							</AlertDescription>
						</Alert>
					</>
				) : (
					<Alert variant="destructive">
						<AlertDescription>
							{t("simple.shared_config_invalid")}
						</AlertDescription>
					</Alert>
				)}
				<DialogFooter>
					<Button variant="outline" onClick={close}>
						{t("viewer.cancel")}
					</Button>
					{pending.ok && (
						<Button
							onClick={() => {
								onLoad(pending.shared.config);
								close();
							}}
						>
							{t("viewer.continue")}
						</Button>
					)}
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
