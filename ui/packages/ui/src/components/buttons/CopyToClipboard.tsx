import { Button, toast } from "@gcsim/primitives";
import { Clipboard } from "lucide-react";
import { memo } from "react";
import { useTranslation } from "react-i18next";

type Props = {
	config?: string;
	className?: string;
};

const CopyTo = ({ config, className }: Props) => {
	const { t } = useTranslation();

	const copy = () => {
		if (config == null) return;
		navigator.clipboard.writeText(config).then(() => {
			toast.success(t("viewer.copied_to_clipboard"), { duration: 2000 });
		});
	};

	return (
		<Button variant="secondary" onClick={copy} disabled={config == null}>
			<Clipboard />
			<div className={className}>{t("viewer.copy")}</div>
		</Button>
	);
};

export default memo(CopyTo);
