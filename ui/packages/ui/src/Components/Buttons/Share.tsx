import {
	Button,
	Classes,
	Dialog,
	Icon,
	InputGroup,
	Intent,
	Label,
	NonIdealState,
	Spinner,
	SpinnerSize,
	type Toaster,
} from "@blueprintjs/core";
import type { SimResults } from "@gcsim/types";
import classNames from "classnames";
import { type RefObject, useState } from "react";
import { useTranslation } from "react-i18next";

type ShareProps = {
	running: boolean;
	copyToast: RefObject<Toaster>;
	data: SimResults | null;
	hash: string | null;
	shareState: [string | null, (link: string | null) => void];
	onShare?: (data: SimResults, hash: string | null) => Promise<string>;
	className?: string;
};

export default ({
	running,
	copyToast,
	data,
	hash,
	className,
	shareState,
	onShare,
}: ShareProps) => {
	const { t } = useTranslation();

	const [isOpen, setOpen] = useState(false);
	const [shareLink, setShareLink] = shareState;

	if (onShare == null) {
		return null;
	}

	const handleShare = () => {
		if (data === null || shareLink != null) {
			return;
		}

		onShare(data, hash)
			.then((url) => {
				setShareLink(url);
			})
			.catch((err) => {
				console.log(err);
			});
	};

	const copy = () => {
		navigator.clipboard.writeText(shareLink ?? "").then(() => {
			copyToast.current?.show({
				message: "Link copied to clipboard!",
				intent: Intent.SUCCESS,
				timeout: 2000,
			});
		});
	};

	return (
		<>
			<Button
				icon={<Icon icon="link" className="!mr-0" />}
				intent={Intent.PRIMARY}
				disabled={running || data == null}
				onClick={() => {
					handleShare();
					setOpen(true);
				}}
			>
				<div className={className}>{t("viewer.share")}</div>
			</Button>
			<Dialog
				isOpen={isOpen}
				onClose={() => setOpen(false)}
				title={t("viewer.create_a_shareable")}
				icon="link"
				className="!pb-0"
			>
				<div
					className={classNames(
						Classes.DIALOG_BODY,
						"flex flex-col justify-center gap-2",
					)}
				>
					<DialogBody shareLink={shareLink} copy={copy} />
				</div>
			</Dialog>
		</>
	);
};

type DialogProps = {
	shareLink: string | null;
	copy: () => void;
};

const DialogBody = ({ shareLink, copy }: DialogProps) => {
	const { t } = useTranslation();
	if (shareLink == null) {
		return <NonIdealState icon={<Spinner size={SpinnerSize.LARGE} />} />;
	}

	return (
		<Label>
			{t("viewer.share_link")}
			<InputGroup
				readOnly={true}
				fill={true}
				onFocus={(e) => {
					e.target.select();
					copy();
				}}
				value={shareLink ?? ""}
				className={classNames({ "bp4-skeleton": shareLink == null })}
				large={true}
				rightElement={<Button icon="duplicate" onClick={() => copy()} />}
			/>
		</Label>
	);
};
