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
import type { Character } from "@gcsim/types";
import React from "react";
import { Trans, useTranslation } from "react-i18next";
import { useAppDispatch } from "../../../../Stores/store";
import { userDataActions } from "../../../../Stores/userDataSlice";
import FetchCharsFromEnka from "./FetchCharsFromEnka";

type Props = {
	isOpen: boolean;
	onClose: () => void;
};

const lsKey = "Enka-UID";

interface CustDivProps
	extends React.DetailedHTMLProps<
		React.HTMLAttributes<HTMLDivElement>,
		HTMLDivElement
	> {
	i18nIsDynamicList?: boolean;
}

const CustDiv: React.FC<CustDivProps> = ({ children, ...props }) => {
	return <div {...props}>{children}</div>;
};

export function ImportFromEnkaDialog(props: Props) {
	const { t } = useTranslation();
	const [message, setMessage] = React.useState<string>("");
	const [errors, setErrors] = React.useState<string[]>([]);
	const [characters, setCharacters] = React.useState<Character[]>([]);
	const [uid, setUid] = React.useState<string>("");
	const dispatch = useAppDispatch();

	async function handleClick() {
		localStorage.setItem(lsKey, uid);
		if (uid && validateUid(uid)) {
			try {
				setCharacters([]);
				const result = await FetchCharsFromEnka(uid);
				setErrors(result.errors ? result.errors : []);
				console.log(result);
				dispatch(
					userDataActions.loadFromGOOD({
						data: result.characters,
						source: "enka",
					}),
				);
				setMessage("success");
				setCharacters(result.characters);
			} catch (e) {
				setMessage(`Error importing chars: ${e}`);
			}
		} else {
			setMessage("Invalid UID");
		}
	}

	return (
		<Dialog
			open={props.isOpen}
			onOpenChange={(open) => {
				if (!open) {
					props.onClose();
					setMessage("");
				}
			}}
		>
			<DialogContent className="sm:max-w-[85%]">
				<DialogHeader>
					<DialogTitle>
						{t("simple.tools_import", { src: "Enka.Network" })}
					</DialogTitle>
				</DialogHeader>
				<div>
					<p className="!pb-2">
						<Trans i18nKey="simple.tools_import_pre_enka">
							{/* biome-ignore lint/a11y/useAnchorContent: text injected at runtime by <Trans> */}
							<a
								href="https://enka.network/"
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
					<input
						value={uid}
						onChange={(e) => {
							setUid(e.target.value.trim());
						}}
						className="w-full p-2 bg-gray-600 rounded-md mt-2"
						placeholder={t("simple.tools_paste_uid")}
					/>

					{message === "success" ? (
						<>
							<Alert variant="success" className="mt-2">
								<AlertDescription>
									{characters.length > 0 ? (
										<Trans i18nKey="simple.tools_import_post_enka">
											<CustDiv i18nIsDynamicList>
												{characters.map((e, i) => {
													return (
														<div
															// biome-ignore lint/suspicious/noArrayIndexKey: character names may duplicate so index completes the composite; list is set wholesale after import, never reordered
															key={e.name + "-" + i}
															className="ml-2"
														>
															{e.name}{" "}
															{e.enka_build_name
																? "(" + e.enka_build_name + ")"
																: ""}
														</div>
													);
												})}
											</CustDiv>
										</Trans>
									) : null}
								</AlertDescription>
							</Alert>
							{errors.length > 0 ? (
								<Alert variant="warning" className="mt-2">
									<AlertDescription>
										Encountered the following issue(s) importing data:
										{errors.map((e, i) => {
											return (
												// biome-ignore lint/suspicious/noArrayIndexKey: string[] error messages, may duplicate, append-only
												<div key={i} className="ml-2">
													{e}
												</div>
											);
										})}
									</AlertDescription>
								</Alert>
							) : null}
						</>
					) : (
						<div>
							{message && (
								<Alert variant="warning" className="mt-2">
									<AlertDescription>{message}</AlertDescription>
								</Alert>
							)}
						</div>
					)}

					<p className="font-bold !pt-2">{t("simple.tools_import_after")}</p>
				</div>
				<DialogFooter>
					<ButtonGroup>
						<Button onClick={handleClick}>{t("simple.import")}</Button>
					</ButtonGroup>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}

function validateUid(uid: string) {
	if (!/^(18|[1-35-9])\d{8}$/.test(uid)) {
		toast.error("Invalid UID");
		return false;
	}
	return true;
}
