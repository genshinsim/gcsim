import {
	Button,
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	Input,
	Label,
	NonIdealState,
	toast,
} from "@gcsim/primitives";
import { Copy, Link } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

export type SignedResult = {
	raw: string;
	hash: string | null;
};

type ShareProps = {
	running: boolean;
	signed: SignedResult | null;
	shareState: [string | null, (link: string | null) => void];
	onShare?: (signed: SignedResult) => Promise<string>;
	className?: string;
};

export default ({
	running,
	signed,
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
		if (signed == null || shareLink != null) {
			return;
		}

		onShare(signed)
			.then((url) => {
				setShareLink(url);
			})
			.catch((err) => {
				console.log(err);
			});
	};

	const copy = () => {
		navigator.clipboard.writeText(shareLink ?? "").then(() => {
			toast.success("Link copied to clipboard!", { duration: 2000 });
		});
	};

	return (
		<>
			<Button
				disabled={running || signed == null}
				onClick={() => {
					handleShare();
					setOpen(true);
				}}
			>
				<Link />
				<div className={className}>{t("viewer.share")}</div>
			</Button>
			<Dialog open={isOpen} onOpenChange={(open) => !open && setOpen(false)}>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>{t("viewer.create_a_shareable")}</DialogTitle>
					</DialogHeader>
					<div className="flex flex-col justify-center gap-2">
						<DialogBody shareLink={shareLink} copy={copy} />
					</div>
				</DialogContent>
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
		return <NonIdealState loading />;
	}

	return (
		<div className="flex flex-col gap-1">
			<Label>{t("viewer.share_link")}</Label>
			<div className="flex gap-1">
				<Input
					readOnly
					value={shareLink}
					onFocus={(e) => {
						e.currentTarget.select();
						copy();
					}}
				/>
				<Button variant="secondary" size="icon" onClick={copy}>
					<Copy />
				</Button>
			</div>
		</div>
	);
};
