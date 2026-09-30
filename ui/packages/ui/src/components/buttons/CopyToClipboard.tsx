import {
	Button,
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	toast,
} from "@gcsim/primitives";
import { Clipboard, FileText, Link } from "lucide-react";
import { memo, useState } from "react";
import { useTranslation } from "react-i18next";
import { sharedConfigUrl } from "../../features/simulator/sharedConfig";

type Props = {
	config?: string;
	className?: string;
};

const CopyTo = ({ config, className }: Props) => {
	const { t } = useTranslation();
	const [isOpen, setOpen] = useState(false);

	const copy = (text: string, message: string) => {
		navigator.clipboard.writeText(text).then(() => {
			toast.success(message, { duration: 2000 });
			setOpen(false);
		});
	};

	const copyText = () => {
		if (config == null) return;
		copy(config, t("viewer.copied_to_clipboard"));
	};

	const copyUrl = () => {
		if (config == null) return;
		copy(
			sharedConfigUrl(window.location.origin, { config, source: "gcsim" }),
			t("viewer.copied_url_to_clipboard"),
		);
	};

	return (
		<>
			<Button
				variant="secondary"
				onClick={() => setOpen(true)}
				disabled={config == null}
			>
				<Clipboard />
				<div className={className}>{t("viewer.copy")}</div>
			</Button>
			<Dialog open={isOpen} onOpenChange={setOpen}>
				<DialogContent className="max-h-[calc(100dvh-2rem)] overflow-y-auto">
					<DialogHeader>
						<DialogTitle className="pr-6">{t("viewer.copy")}</DialogTitle>
					</DialogHeader>
					<div className="flex flex-col gap-2 sm:flex-row">
						<Button
							variant="secondary"
							className="w-full sm:flex-1"
							onClick={copyText}
						>
							<FileText />
							{t("viewer.copy_as_text")}
						</Button>
						<Button className="w-full sm:flex-1" onClick={copyUrl}>
							<Link />
							{t("viewer.copy_as_url")}
						</Button>
					</div>
				</DialogContent>
			</Dialog>
		</>
	);
};

export default memo(CopyTo);
