import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	Checkbox,
	Label,
} from "@gcsim/primitives";
import {
	createContext,
	type MouseEvent,
	type ReactNode,
	useContext,
	useState,
} from "react";
import { useTranslation } from "react-i18next";
import { loadSkipKqmWarning, saveSkipKqmWarning } from "../../stores/prefs";

type KqmLinkClick = (e: MouseEvent<HTMLAnchorElement>) => void;

const KqmLinkClickContext = createContext<KqmLinkClick | null>(null);

export function useKqmLinkClick(): KqmLinkClick {
	const onClick = useContext(KqmLinkClickContext);
	if (onClick == null) {
		throw new Error("useKqmLinkClick must be used within KqmLeaveWarning");
	}
	return onClick;
}

const isPlainClick = (e: MouseEvent) =>
	e.button === 0 && !e.metaKey && !e.ctrlKey && !e.shiftKey && !e.altKey;

export function KqmLeaveWarning({ children }: { children: ReactNode }) {
	const { t } = useTranslation();
	const [pendingUrl, setPendingUrl] = useState<string | null>(null);
	const [dontShow, setDontShow] = useState(false);

	const onLinkClick: KqmLinkClick = (e) => {
		if (!isPlainClick(e) || loadSkipKqmWarning(localStorage)) return;
		e.preventDefault();
		setDontShow(false);
		setPendingUrl(e.currentTarget.href);
	};

	const onContinue = () => {
		if (pendingUrl == null) return;
		if (dontShow) saveSkipKqmWarning(localStorage, true);
		window.open(pendingUrl, "_blank", "noopener,noreferrer");
		setPendingUrl(null);
	};

	return (
		<KqmLinkClickContext.Provider value={onLinkClick}>
			{children}
			<AlertDialog
				open={pendingUrl != null}
				onOpenChange={(open) => !open && setPendingUrl(null)}
			>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle>{t("dash.kqm_leave_title")}</AlertDialogTitle>
						<AlertDialogDescription>
							{t("dash.kqm_leave_body")}
						</AlertDialogDescription>
					</AlertDialogHeader>
					<div className="flex items-center gap-2">
						<Checkbox
							id="kqm-leave-dont-show"
							checked={dontShow}
							onCheckedChange={(checked) => setDontShow(checked === true)}
						/>
						<Label htmlFor="kqm-leave-dont-show">
							{t("dash.kqm_leave_dont_show")}
						</Label>
					</div>
					<AlertDialogFooter>
						<AlertDialogCancel>{t("dash.kqm_leave_cancel")}</AlertDialogCancel>
						<AlertDialogAction onClick={onContinue}>
							{t("dash.kqm_leave_continue")}
						</AlertDialogAction>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>
		</KqmLinkClickContext.Provider>
	);
}
