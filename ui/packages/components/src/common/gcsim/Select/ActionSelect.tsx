import type { IAction } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { actionLabel, actionPredicate, actions } from "./actions";
import { OmniSelect } from "./OmniSelect";
import { renderLabeledItem } from "./utils";

export type ActionSelectProps = {
	isOpen: boolean;
	onClose: () => void;
	onSelect: (action: IAction) => void;
	value?: IAction;
};

export function ActionSelect({
	isOpen,
	onClose,
	onSelect,
	value,
}: ActionSelectProps) {
	const { t } = useTranslation();
	return (
		<OmniSelect<IAction>
			isOpen={isOpen}
			onClose={onClose}
			items={actions}
			itemKey={(action) => action}
			itemPredicate={actionPredicate}
			itemRenderer={(action, state) =>
				renderLabeledItem(action, actionLabel(action), state)
			}
			onSelect={onSelect}
			value={value}
			title={t("simple.actions")}
			placeholder={t("db.type_to_search")}
		/>
	);
}
