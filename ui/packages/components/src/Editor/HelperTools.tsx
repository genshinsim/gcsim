import {
	Button,
	DropdownMenu,
	DropdownMenuCheckboxItem,
	DropdownMenuContent,
	DropdownMenuTrigger,
} from "@gcsim/primitives";
import { HelpCircle, Search, Users, Wrench } from "lucide-react";
import { useTranslation } from "react-i18next";

export interface EditorToggles {
	team: boolean;
	nameSearch: boolean;
	tips: boolean;
}

export interface HelperToolsProps {
	toggles: EditorToggles;
	onToggle: (key: keyof EditorToggles) => void;
	className?: string;
}

export function HelperTools({
	toggles,
	onToggle,
	className,
}: HelperToolsProps) {
	const { t } = useTranslation();
	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button
					variant="secondary"
					className={className}
					data-testid="editor-helper-tools"
				>
					<Wrench />
					{t("simple.tools")}
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent side="top">
				<DropdownMenuCheckboxItem
					checked={toggles.team}
					onCheckedChange={() => onToggle("team")}
				>
					<Users />
					{t("simple.team")}
				</DropdownMenuCheckboxItem>
				<DropdownMenuCheckboxItem
					checked={toggles.nameSearch}
					onCheckedChange={() => onToggle("nameSearch")}
				>
					<Search />
					{t("simple.name_search")}
				</DropdownMenuCheckboxItem>
				<DropdownMenuCheckboxItem
					checked={toggles.tips}
					onCheckedChange={() => onToggle("tips")}
				>
					<HelpCircle />
					{t("simple.tips")}
				</DropdownMenuCheckboxItem>
			</DropdownMenuContent>
		</DropdownMenu>
	);
}
