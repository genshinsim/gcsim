import {
	Button,
	Callout,
	Checkbox,
	Classes,
	Dialog,
	Icon,
	Intent,
} from "@blueprintjs/core";
import classNames from "classnames";
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
				icon={<Icon icon="send-to" className="!mr-0" />}
				onClick={() => setOpen(true)}
				disabled={config == null}
			>
				<div className="hidden ml-[7px] sm:flex">
					{t("viewer.send_to_simulator")}
				</div>
			</Button>
			<Dialog
				isOpen={isOpen}
				onClose={() => setOpen(false)}
				title={t("viewer.load_this_configuration")}
				icon="bring-data"
			>
				<div className={Classes.DIALOG_BODY}>
					<Callout intent="warning" className="">
						{t("viewer.this_will_overwrite")}
					</Callout>
					<Checkbox
						label={t("viewer.copy_list_only")}
						className="my-3 mx-1"
						checked={keepTeam}
						onClick={toggleKeepTeam}
					/>
				</div>
				<div
					className={classNames(
						Classes.DIALOG_FOOTER,
						Classes.DIALOG_FOOTER_ACTIONS,
					)}
				>
					<Button
						onClick={toSimulator}
						intent={Intent.PRIMARY}
						text={t("viewer.continue")}
					/>
					<Button onClick={() => setOpen(false)} text={t("viewer.cancel")} />
				</div>
			</Dialog>
		</>
	);
};

export default memo(SendTo);
