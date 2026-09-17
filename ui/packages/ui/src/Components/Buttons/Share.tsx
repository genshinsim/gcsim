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
import type { SimResults } from "@gcsim/types";
import { Copy, Link } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

type ShareProps = {
	running: boolean;
	data: SimResults | null;
	hash: string | null;
	shareState: [string | null, (link: string | null) => void];
	onShare?: (data: SimResults, hash: string | null) => Promise<string>;
	className?: string;
};

export default ({
	running,
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
			toast.success("Link copied to clipboard!", { duration: 2000 });
		});
	};

	return (
		<>
			<Button
				disabled={running || data == null}
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
