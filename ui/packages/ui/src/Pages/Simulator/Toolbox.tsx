import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import {
	Button,
	ButtonGroup,
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
	Spinner,
} from "@gcsim/primitives";
import type { SimResults } from "@gcsim/types";
import { throttle } from "lodash-es";
import {
	Download,
	HelpCircle,
	Play,
	Scissors,
	Search,
	Upload,
	Users,
	Wrench,
} from "lucide-react";
import React from "react";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router";
import ExecutorSettingsButton from "../../Components/Buttons/ExecutorSettingsButton";
import {
	type AppThunk,
	type RootState,
	useAppDispatch,
	useAppSelector,
} from "../../Stores/store";
import { userActions } from "../../Stores/userSlice";
import { viewerActions } from "../../Stores/viewerSlice";
import { VIEWER_THROTTLE } from "../Viewer";
import { ImportFromEnkaDialog, ImportFromGOODDialog } from "./Components";

type Props = {
	exec: ExecutorSupplier<Executor>;
	cfg: string;
	isReady: boolean;
	isValid: boolean;
};

export function runSim(pool: Executor, cfg: string): AppThunk {
	return (dispatch) => {
		console.log("starting run");
		dispatch(viewerActions.start());

		const updateResult = throttle(
			(res: SimResults, hash: string) => {
				dispatch(viewerActions.setResult({ data: res, hash: hash }));
			},
			VIEWER_THROTTLE,
			{ leading: true, trailing: true },
		);

		pool.run(cfg, updateResult).catch((err) => {
			dispatch(viewerActions.setError({ recoveryConfig: cfg, error: err }));
		});
	};
}

export const Toolbox = ({ exec, cfg, isReady, isValid }: Props) => {
	const { t } = useTranslation();
	const history = useHistory();

	const [openImport, setOpenGOODImport] = React.useState<boolean>(false);
	const [openImportFromEnka, setOpenImportFromEnka] =
		React.useState<boolean>(false);
	const { settings } = useAppSelector((state: RootState) => {
		return {
			settings: state.user.data.settings,
		};
	});

	const dispatch = useAppDispatch();
	const toggleTips = () => {
		dispatch(
			userActions.setUserSettings({
				showTips: !settings.showTips,
				showBuilder: settings.showBuilder,
				showNameSearch: settings.showNameSearch,
			}),
		);
	};

	const run = () => {
		dispatch(runSim(exec(), cfg));
		history.push("/web");
	};

	const toggleBuilder = () => {
		dispatch(
			userActions.setUserSettings({
				showTips: settings.showTips,
				showBuilder: !settings.showBuilder,
				showNameSearch: settings.showNameSearch,
			}),
		);
	};

	const toggleNameSearch = () => {
		dispatch(
			userActions.setUserSettings({
				showTips: settings.showTips,
				showBuilder: settings.showBuilder,
				showNameSearch: !settings.showNameSearch,
			}),
		);
	};

	return (
		<div className="p-2 wide:ml-2 wide:mr-2 flex flex-row flex-wrap place-items-center gap-x-1 gap-y-1">
			<div className="basis-full wide:basis-0 flex-grow p-1 flex flex-row items-center">
				<ExecutorSettingsButton />
			</div>
			<ButtonGroup className="basis-full wide:basis-2/3 p-1 w-full flex-wrap">
				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<Button variant="secondary" className="basis-full md:basis-1/2">
							<Wrench />
							{t("simple.tools")}
						</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent side="top">
						<DropdownMenuItem onClick={toggleTips}>
							<HelpCircle />
							{settings.showTips
								? t("simple.tools_hide_tooltips")
								: t("simple.tools_show_tooltips")}
						</DropdownMenuItem>
						<DropdownMenuItem onClick={toggleBuilder}>
							<Users />
							{settings.showBuilder
								? t("simple.tools_hide_builder")
								: t("simple.tools_show_builder")}
						</DropdownMenuItem>
						<DropdownMenuItem onClick={toggleNameSearch}>
							<Search />
							{settings.showNameSearch
								? t("simple.tools_hide_name_search")
								: t("simple.tools_show_name_search")}
						</DropdownMenuItem>
						<DropdownMenuSeparator />
						<DropdownMenuItem onClick={() => history.push("/sample/upload")}>
							<Upload />
							{t("simple.tools_sample_upload")}
						</DropdownMenuItem>
						<DropdownMenuItem disabled>
							<Scissors />
							{t("simple.tools_substat_snippets")}
						</DropdownMenuItem>
						<DropdownMenuSeparator />
						<DropdownMenuItem onClick={() => setOpenGOODImport(true)}>
							<Download />
							{t("simple.tools_import", { src: "GO" })}
						</DropdownMenuItem>
						<DropdownMenuItem onClick={() => setOpenImportFromEnka(true)}>
							<Download />
							{t("simple.tools_import", { src: "Enka" })}
						</DropdownMenuItem>
					</DropdownMenuContent>
				</DropdownMenu>
				<Button
					className="basis-full md:basis-1/2"
					onClick={run}
					disabled={!isReady || !isValid}
				>
					{isReady ? <Play /> : <Spinner />}
					{t("simple.run")}
				</Button>
			</ButtonGroup>
			<ImportFromGOODDialog
				isOpen={openImport}
				onClose={() => setOpenGOODImport(false)}
			/>
			<ImportFromEnkaDialog
				isOpen={openImportFromEnka}
				onClose={() => setOpenImportFromEnka(false)}
			/>
		</div>
	);
};
