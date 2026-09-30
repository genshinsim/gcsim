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

type Pending = { ok: true; shared: SharedConfig } | { ok: false };

function decode(encoded: string): Pending {
	try {
		return { ok: true, shared: decodeSharedConfig(encoded) };
	} catch {
		return { ok: false };
	}
}

// Reads `?cfg=` once, strips it from the URL so a refresh can't re-trigger it,
// and asks before overwriting the draft.
export function SharedConfigDialog({
	onLoad,
}: {
	onLoad: (cfg: string) => void;
}) {
	const { t } = useTranslation();
	const { cfg } = routeApi.useSearch();
	const navigate = routeApi.useNavigate();
	const [pending, setPending] = React.useState<Pending | null>(null);

	React.useEffect(() => {
		if (!cfg) return;
		setPending(decode(cfg));
		void navigate({ to: "/simulator", search: {}, replace: true });
	}, [cfg, navigate]);

	const close = () => setPending(null);

	if (pending == null) return null;

	return (
		<Dialog open onOpenChange={(open) => !open && close()}>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>{t("simple.shared_config_heading")}</DialogTitle>
				</DialogHeader>
				{pending.ok ? (
					<>
						<dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1">
							<dt className="text-g-ink-mute">
								{t("simple.shared_config_title")}
							</dt>
							<dd className="break-words">
								{pending.shared.title ?? t("simple.shared_config_unknown")}
							</dd>
							<dt className="text-g-ink-mute">
								{t("simple.shared_config_source")}
							</dt>
							<dd className="break-words">
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
