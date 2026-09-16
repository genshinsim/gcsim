import { Button, Icon } from "@blueprintjs/core";
import { toast } from "@gcsim/primitives";
import { memo } from "react";
import { useTranslation } from "react-i18next";

type Props = {
	config?: string;
	className?: string;
};

const CopyTo = ({ config, className }: Props) => {
	const { t } = useTranslation();

	const action = () => {
		navigator.clipboard.writeText(config ?? "").then(() => {
			toast.success(t("viewer.copied_to_clipboard"), { duration: 2000 });
		});
	};

	return (
		<>
			<Button
				icon={<Icon icon="clipboard" className="!mr-0" />}
				onClick={action}
				disabled={config == null}
			>
				<div className={className}>{t("viewer.copy")}</div>
			</Button>
		</>
	);
};

export default memo(CopyTo);
