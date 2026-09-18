import { type ExecutorSupplier, ServerExecutor } from "@gcsim/executors";
import {
	Field,
	FieldDescription,
	FieldLabel,
	FieldTitle,
	Input,
} from "@gcsim/primitives";
import { UI } from "@gcsim/ui";
import React, { type ReactNode, useRef } from "react";
import { useTranslation } from "react-i18next";

let exec: ServerExecutor | undefined;
const urlKey = "server-mode-url";
const defaultURL = "http://127.0.0.1:54321";

const ServerMode = ({ children }: { children: ReactNode }) => {
	const { t } = useTranslation();
	const [url, setURL] = React.useState<string>((): string => {
		const saved = localStorage.getItem(urlKey);
		if (saved === null) {
			localStorage.setItem(urlKey, defaultURL);
			return defaultURL;
		}
		return saved;
	});
	React.useEffect(() => {
		localStorage.setItem(urlKey, url);
	}, [url]);

	React.useEffect(() => {
		if (exec != null) {
			exec.set_url(url);
		}
	}, [url]);

	const supplier = useRef<ExecutorSupplier<ServerExecutor>>(() => {
		if (exec == null) {
			exec = new ServerExecutor(url);
		}
		return exec;
	});

	return (
		<UI
			exec={supplier.current}
			gitCommit={import.meta.env.VITE_GIT_COMMIT_HASH}
			mode={import.meta.env.MODE}
		>
			<Field>
				<FieldTitle>{t("simple.workers")}</FieldTitle>
				{children}
				<Field>
					<FieldLabel htmlFor="server-mode-url">
						{t("simple.server_mode_url")}
						<span className="text-deprecated-muted-foreground">
							{t("simple.server_mode_required")}
						</span>
					</FieldLabel>
					<Input
						id="server-mode-url"
						value={url}
						onChange={(e) => {
							setURL(e.target.value);
						}}
					/>
					<FieldDescription>
						{t("simple.server_mode_default") + defaultURL}
					</FieldDescription>
				</Field>
			</Field>
		</UI>
	);
};

export default ServerMode;
