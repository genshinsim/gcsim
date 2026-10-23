import {
	Button,
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@gcsim/primitives";
import { Download } from "lucide-react";
import React from "react";
import { useTranslation } from "react-i18next";
import { ImportFromEnkaDialog } from "./enka/ImportFromEnkaDialog";
import { ImportFromGOODDialog } from "./good/ImportFromGOODDialog";

export function EditorSettings() {
	const { t } = useTranslation();
	const [openGOODImport, setOpenGOODImport] = React.useState(false);
	const [openEnkaImport, setOpenEnkaImport] = React.useState(false);

	return (
		<div className="flex flex-row flex-wrap items-center gap-1">
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button variant="secondary">
						<Download />
						{t("simple.import")}
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent side="top">
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
