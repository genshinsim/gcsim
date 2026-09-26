import {
	Button,
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@gcsim/primitives";
import { Download, Upload } from "lucide-react";
import React from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import ExecutorSettingsButton from "../../Components/Buttons/ExecutorSettingsButton";
import { ImportFromEnkaDialog, ImportFromGOODDialog } from "./Components";

export function EditorSettings() {
	const { t } = useTranslation();
	const navigate = useNavigate();
	const [openGOODImport, setOpenGOODImport] = React.useState(false);
	const [openEnkaImport, setOpenEnkaImport] = React.useState(false);

	return (
		<div className="flex flex-row flex-wrap items-center gap-1">
			<ExecutorSettingsButton />
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button variant="secondary">
						<Download />
						{t("simple.import")}
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent side="top">
					<DropdownMenuItem onClick={() => navigate("/sample/upload")}>
						<Upload />
						{t("simple.tools_sample_upload")}
					</DropdownMenuItem>
					<DropdownMenuSeparator />
					<DropdownMenuItem onClick={() => setOpenGOODImport(true)}>
						<Download />
						{t("simple.tools_import", { src: "GO" })}
					</DropdownMenuItem>
					<DropdownMenuItem onClick={() => setOpenEnkaImport(true)}>
						<Download />
						{t("simple.tools_import", { src: "Enka" })}
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
			<ImportFromGOODDialog
				isOpen={openGOODImport}
				onClose={() => setOpenGOODImport(false)}
			/>
			<ImportFromEnkaDialog
				isOpen={openEnkaImport}
				onClose={() => setOpenEnkaImport(false)}
			/>
		</div>
	);
}
