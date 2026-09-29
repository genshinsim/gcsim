import {
	Alert,
	AlertDescription,
	Button,
	Checkbox,
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	Label,
} from "@gcsim/primitives";
import { Send } from "lucide-react";
import { memo, useState } from "react";
import { useTranslation } from "react-i18next";

type SendToProps = {
	config?: string;
	onSendToSimulator?: (cfg: string, opts: { keepTeam: boolean }) => void;
};

const SendTo = ({ config, onSendToSimulator }: SendToProps) => {
	const LOCALSTORAGE_KEY = "gcsim-viewer-cpy-cfg-settings";
	const { t } = useTranslation();

	const [isOpen, setOpen] = useState(false);
	const [keepTeam, setKeep] = useState<boolean>(() => {
		return localStorage.getItem(LOCALSTORAGE_KEY) === "true";
	});

	if (onSendToSimulator == null) {
		return null;
	}

	const toggleKeepTeam = () => {
		localStorage.setItem(LOCALSTORAGE_KEY, String(!keepTeam));
		setKeep(!keepTeam);
	};

	const toSimulator = () => {
		if (config == null) {
			return;
		}
		onSendToSimulator(config, { keepTeam });
	};

	return (
		<>
			<Button
				variant="secondary"
				onClick={() => setOpen(true)}
				disabled={config == null}
			>
				<Send />
				<div className="hidden ml-[7px] sm:flex">
					{t("viewer.send_to_simulator")}
				</div>
			</Button>
			<Dialog open={isOpen} onOpenChange={(open) => !open && setOpen(false)}>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>{t("viewer.load_this_configuration")}</DialogTitle>
					</DialogHeader>
					<Alert variant="warning">
						<AlertDescription>
							{t("viewer.this_will_overwrite")}
						</AlertDescription>
					</Alert>
					<div className="flex items-center gap-2">
						<Checkbox
							id="keep-team"
							checked={keepTeam}
							onCheckedChange={toggleKeepTeam}
						/>
						<Label htmlFor="keep-team">{t("viewer.copy_list_only")}</Label>
					</div>
					<DialogFooter>
						<Button variant="outline" onClick={() => setOpen(false)}>
							{t("viewer.cancel")}
						</Button>
						<Button onClick={toSimulator}>{t("viewer.continue")}</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</>
	);
};

export default memo(SendTo);
