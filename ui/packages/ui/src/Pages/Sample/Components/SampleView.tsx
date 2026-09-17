import { Button, ButtonGroup, Card, Input, Label } from "@gcsim/primitives";
import type { Sample } from "@gcsim/types";
import { useVirtualizer } from "@tanstack/react-virtual";
import { saveAs } from "file-saver";
import { ArrowDown, Download, RotateCcw, Settings } from "lucide-react";
import Pako from "pako";
import React, { useState } from "react";
import { useTranslation } from "react-i18next";
import AutoSizer from "react-virtualized-auto-sizer";
import { Options } from "./Options";
import type { SampleItem, SampleRow } from "./parse";
import { SampleItemView } from "./SampleItemView";
import {
	AdvancedPreset,
	AllSampleOptions,
	DebugPreset,
	DefaultSampleOptions,
	SimplePreset,
	VerbosePreset,
} from "./SampleOptions";

type buffSetting = {
	start: number;
	end: number;
	show: boolean;
};

const Row = ({
	row,
	highlight,
	showBuffDuration,
}: {
	row: SampleRow;
	highlight: buffSetting;
	showBuffDuration: (e: SampleItem) => void;
}) => {
	const cols = row.slots.map((slot, ci) => {
		const events = slot.map((e, ei) => {
			return (
				// biome-ignore lint/suspicious/noArrayIndexKey: parsed sample row cells, static once parsed, no unique field
				<SampleItemView item={e} key={ei} showBuffDuration={showBuffDuration} />
			);
		});

		return (
			<div
				// biome-ignore lint/suspicious/noArrayIndexKey: fixed positional team slot (0–4)
				key={ci}
				className={
					row.active === ci
						? "border-l-2 border-gray-500 bg-gray-400	"
						: "border-l-2 border-gray-500"
				}
			>
				{events}
			</div>
		);
	});

	const hl =
		highlight.show && row.f >= highlight.start && row.f <= highlight.end;

	//map out each col
	return (
		<div className="flex flex-row" key={row.key}>
			<div
				className={
					hl
						? "text-right text-gray-100 border-b-2 border-gray-500 bg-blue-500"
						: "text-right text-gray-100 border-b-2 border-gray-500"
				}
				style={{ minWidth: "100px" }}
			>
				<div>{`${row.f} | ${(row.f / 60).toFixed(2)}s`}</div>
			</div>
			<div className="grid grid-cols-5 flex-grow border-b-2 border-gray-500">
				{cols}
			</div>
			<div style={{ width: "20px", minWidth: "20px" }} />
		</div>
	);
};

type SampleOptionsProps = {
	settings: string[];
	setSettings: (val: string[]) => void;
};

const SampleOptions = ({ settings, setSettings }: SampleOptionsProps) => {
	const { t } = useTranslation();
	const [isOpen, setOpen] = useState(false);

	const toggle = (t: string) => {
		const i = settings.indexOf(t);
		const next = [...settings];
		if (i === -1) {
			next.push(t);
		} else {
			next.splice(i, 1);
		}
		setSettings(next);
	};

	const presets = (opt: "simple" | "advanced" | "verbose" | "debug") => {
		switch (opt) {
			case "simple":
				setSettings(SimplePreset);
				return;
			case "advanced":
				setSettings(AdvancedPreset);
				return;
			case "verbose":
				setSettings(VerbosePreset);
				return;
			case "debug":
				setSettings(DebugPreset);
				return;
		}
	};

	return (
		<>
			<Button variant="secondary" onClick={() => setOpen(true)}>
				<Settings />
				{t("simple.settings")}
			</Button>
			<Options
				isOpen={isOpen}
				handleClose={() => setOpen(false)}
				handleClear={() => setSettings([])}
				handleResetDefault={() => setSettings(DefaultSampleOptions)}
				handleToggle={toggle}
				handleSetPresets={presets}
				selected={settings}
				options={AllSampleOptions}
			/>
		</>
	);
};

let lastSearchIndex = 0;

type SamplerProps = {
	sample: Sample;
	data: SampleRow[];
	team: string[];
	searchable: { [key: number]: string[] };
	settings: string[];
	setSettings: (val: string[]) => void;
};

function SamplerUI({
	sample,
	data,
	team,
	searchable,
	settings,
	setSettings,
}: SamplerProps) {
	const { t } = useTranslation();
	// State-backed ref (not useRef) so the virtualizer re-measures once AutoSizer
	// mounts the scroll element — AutoSizer defers rendering its child until it has
	// a non-zero size, so a plain ref is still null on the virtualizer's first pass.
	const [scrollParent, setScrollParent] = React.useState<HTMLDivElement | null>(
		null,
	);
	const searchRef = React.useRef<HTMLInputElement>(null);
	const [hl, sethl] = React.useState<buffSetting>({
		start: 0,
		end: 0,
		show: false,
	});

	const handleShowBuffDuration = (e: SampleItem) => {
		// const show = hl.show;
		const next = {
			show: true,
			start: e.added,
			end: e.ended,
		};
		sethl(next);
	};

	const rowVirtualizer = useVirtualizer({
		count: data.length,
		getScrollElement: () => scrollParent,
		estimateSize: () => 30,
		getItemKey: React.useCallback(
			(index: number) => {
				return data[index].f;
			},
			[data],
		),
	});

	const char = team.map((c) => {
		return (
			<div
				key={c}
				className="capitalize text-lg font-medium text-gray-100 border-l-2 border-b-2 border-gray-500"
			>
				{c}
			</div>
		);
	});

	const searchAndScroll = (val: string) => {
		const total = Object.keys(searchable).length;
		for (let index = lastSearchIndex; index < total; index++) {
			for (const msg of searchable[index]) {
				const lowerMsg = msg.toLowerCase();
				if (lowerMsg.indexOf(val.toLowerCase()) > -1) {
					console.log(index, lastSearchIndex);
					lastSearchIndex = index + 1;
					rowVirtualizer.scrollToIndex(index, { align: "start" });
					return;
				}
			}
		}
	};

	return (
		<>
			<div className="flex flex-col sm:flex-row justify-between">
				<div className="flex flex-row items-center gap-2">
					<Label>{t("viewer.search")}</Label>
					<div className="flex flex-row gap-1">
						<Input type="text" ref={searchRef} />
						<Button
							variant="secondary"
							size="icon"
							onClick={() => {
								if (searchRef.current != null) {
									searchAndScroll(searchRef.current.value);
								}
							}}
						>
							<ArrowDown />
						</Button>
						<Button
							variant="secondary"
							size="icon"
							onClick={() => {
								if (searchRef.current != null) {
									searchRef.current.value = "";
								}
								lastSearchIndex = 0;
								rowVirtualizer.scrollToIndex(0);
							}}
						>
							<RotateCcw />
						</Button>
					</div>
				</div>
				<ButtonGroup className="mb-[15px]">
					<SampleOptions settings={settings} setSettings={setSettings} />
					<Button
						variant="secondary"
						onClick={() => {
							const out = Pako.deflate(JSON.stringify(sample));
							const blob = new Blob([out], { type: "application/base64" });
							saveAs(blob, "sample.gz");
						}}
					>
						<Download />
						{t("viewer.download")}
					</Button>
				</ButtonGroup>
			</div>
			<div className="flex flex-col overflow-x-auto h-[80vh]">
				<Card className="flex-auto gap-0 p-2 !bg-gray-600 !text-xs min-w-[60rem] ">
					<AutoSizer disableWidth={true}>
						{({ height }) => (
							<div
								ref={setScrollParent}
								style={{
									minHeight: "100px",
									height: height,
									overflow: "auto",
									position: "relative",
								}}
								id="resize-inner"
							>
								<div className="flex flex-row sample-header">
									<div
										className={
											"font-medium text-lg text-gray-100 border-b-2 border-gray-500 text-right "
										}
										style={{ minWidth: "100px" }}
									>
										F | Sec
									</div>
									<div className="grid grid-cols-5 flex-grow">
										<div className="font-medium text-lg text-gray-100 border-l-2 border-b-2 border-gray-500">
											Sim
										</div>
										{char}
									</div>
									<div style={{ width: "20px", minWidth: "20px" }} />
								</div>
								<div
									className="ListInner"
									style={{
										height: `${rowVirtualizer.getTotalSize()}px`,
										width: "100%",
										position: "relative",
									}}
								>
									{rowVirtualizer.getVirtualItems().map((virtualRow) => (
										<div
											key={virtualRow.index}
											ref={rowVirtualizer.measureElement}
											data-index={virtualRow.index}
											style={{
												position: "absolute",
												top: 0,
												left: 0,
												width: "100%",
												// Positions the virtual elements at the right place in container.
												// minHeight: `${virtualRow.size - 10}px`,
												transform: `translateY(${virtualRow.start}px)`,
											}}
											// id={"virtual-row-"+virtualRow.key}
										>
											<Row
												row={data[virtualRow.index]}
												highlight={hl}
												showBuffDuration={handleShowBuffDuration}
											/>
										</div>
									))}
								</div>
							</div>
						)}
					</AutoSizer>
				</Card>
			</div>
		</>
	);
}

export const Sampler = React.memo(SamplerUI);
