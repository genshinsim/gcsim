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
import axios from "axios";
import { Copy, Link, TriangleAlert } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

export type SignedResult = {
	raw: string;
	hash: string | null;
};

type ShareErrorKey = "viewer.share_rate_limited" | "viewer.share_failed";

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
	const [error, setError] = useState<ShareErrorKey | null>(null);

	if (onShare == null) {
		return null;
	}

	const handleShare = () => {
		if (signed == null || shareLink != null) {
			return;
		}

		setError(null);
		onShare(signed)
			.then((url) => {
				setShareLink(url);
			})
			.catch((err) => {
				setError(
					axios.isAxiosError(err) && err.response?.status === 429
						? "viewer.share_rate_limited"
						: "viewer.share_failed",
				);
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
						<DialogBody
							shareLink={shareLink}
							error={error}
							copy={copy}
							retry={handleShare}
						/>
					</div>
				</DialogContent>
			</Dialog>
		</>
	);
};

type DialogProps = {
	shareLink: string | null;
	error: ShareErrorKey | null;
	copy: () => void;
	retry: () => void;
};

const DialogBody = ({ shareLink, error, copy, retry }: DialogProps) => {
	const { t } = useTranslation();
	if (error != null) {
		return (
			<NonIdealState
				icon={<TriangleAlert />}
				description={t(error)}
				action={
					<Button variant="secondary" onClick={retry}>
						{t("viewer.share_retry")}
					</Button>
				}
			/>
		);
	}
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
