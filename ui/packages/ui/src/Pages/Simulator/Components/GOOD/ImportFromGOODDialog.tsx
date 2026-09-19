import {
	Alert,
	AlertDescription,
	Button,
	ButtonGroup,
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	toast,
} from "@gcsim/primitives";
import React from "react";
import { Trans, useTranslation } from "react-i18next";
import { useAppDispatch } from "../../../../Stores/store";
import { userDataActions } from "../../../../Stores/userDataSlice";
import { type IGOODImport, parseFromGOOD } from "./parseFromGOOD";

type Props = {
	isOpen: boolean;
	onClose: () => void;
};

const lsKey = "GOOD-import";

export function ImportFromGOODDialog(props: Props) {
	const [data, setData] = React.useState<IGOODImport>();
	const dispatch = useAppDispatch();
	const { t } = useTranslation();

	const handleLoad = () => {
		if (data !== undefined) {
			dispatch(
				userDataActions.loadFromGOOD({ data: data.characters, source: "good" }),
			);
			props.onClose();
			toast.success(t("importer.import_success"));
		}
	};
	const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
		localStorage.setItem(lsKey, e.target.value);
		setData(parseFromGOOD(e.target.value));
	};
	return (
		<Dialog
			open={props.isOpen}
			onOpenChange={(open) => !open && props.onClose()}
		>
			<DialogContent className="sm:max-w-[85%]">
				<DialogHeader>
					<DialogTitle>
						{t("simple.tools_import", { src: "Genshin Optimizer/GOOD" })}
					</DialogTitle>
				</DialogHeader>
				<div>
					<p className="!pb-2">
						<Trans i18nKey="simple.tools_import_pre_go">
							{/* biome-ignore lint/a11y/useAnchorContent: text injected at runtime by <Trans> */}
							<a
								href="https://frzyc.github.io/genshin-optimizer/#/setting"
								target="_blank"
								rel="noreferrer"
							/>
						</Trans>
					</p>
					<Alert variant="warning">
						<AlertDescription>
							{t("simple.tools_import_warning", { src: "GOOD/Enka" })}
						</AlertDescription>
					</Alert>
					<textarea
						value={localStorage.getItem(lsKey) ?? ""}
						onChange={handleChange}
						className="w-full p-2 bg-g-surface-2 rounded-g-md mt-2"
						rows={7}
					/>
					<p className="font-bold !pt-2">{t("simple.tools_import_after")}</p>
					{data ? (
						data.err === "" ? (
							<Alert variant="success" className="mt-2">
								<AlertDescription>
									{t("simple.tools_import_post_go")}
								</AlertDescription>
							</Alert>
						) : (
							<Alert variant="warning" className="mt-2">
								<AlertDescription>{data.err}</AlertDescription>
							</Alert>
						)
					) : null}
				</div>
				<DialogFooter>
					<ButtonGroup>
						<Button onClick={handleLoad} disabled={!data || data.err !== ""}>
							{t("simple.import")}
						</Button>
						<Button onClick={props.onClose} variant="destructive">
							{t("db.cancel")}
						</Button>
					</ButtonGroup>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
