import {
	Button,
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	toast,
} from "@gcsim/primitives";
import { Search } from "lucide-react";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import {
	buildShortcutEntries,
	type ShortcutEntry,
	type ShortcutKind,
	type ShortcutMatch,
	searchShortcuts,
} from "./shortcutIndex";

export function ShortcutSearch() {
	const { t } = useTranslation();
	const [open, setOpen] = useState(false);
	const [query, setQuery] = useState("");
	const [entries, setEntries] = useState<ShortcutEntry[]>([]);

	const groups = useMemo(() => {
		const out = new Map<ShortcutKind, ShortcutMatch[]>();
		for (const match of searchShortcuts(entries, query)) {
			const list = out.get(match.entry.kind);
			if (list) list.push(match);
			else out.set(match.entry.kind, [match]);
		}
		return out;
	}, [entries, query]);

	const headings: Record<ShortcutKind, string> = {
		character: t("db.characters"),
		weapon: t("simple.weapons"),
		artifact: t("simple.artifacts"),
		enemy: t("simple.enemies"),
		action: t("simple.actions"),
		stat: t("simple.stats"),
	};

	const copy = (item: string) => {
		setOpen(false);
		navigator.clipboard.writeText(item).then(
			() => {
				toast.success(t("simple.copied_to_clipboard", { item }), {
					duration: 2000,
				});
			},
			() => toast.error(t("simple.copy_failed", { item })),
		);
	};

	return (
		<>
			<Button
				variant="ghost"
				size="sm"
				aria-label={t("simple.shortcut_search")}
				onClick={() => {
					setEntries(buildShortcutEntries());
					setQuery("");
					setOpen(true);
				}}
			>
				<Search />
				<span className="hidden sm:inline">{t("simple.shortcut_search")}</span>
			</Button>
			<Dialog open={open} onOpenChange={setOpen}>
				<DialogContent
					className="gap-0 overflow-hidden p-0"
					showCloseButton={false}
				>
					<DialogHeader className="sr-only">
						<DialogTitle>{t("simple.shortcut_search")}</DialogTitle>
						<DialogDescription>{t("db.type_to_search")}</DialogDescription>
					</DialogHeader>
					<Command shouldFilter={false}>
						<CommandInput
							value={query}
							onValueChange={setQuery}
							placeholder={t("db.type_to_search")}
						/>
						<CommandList className="scrollbar-surface">
							<CommandEmpty>{t("common.data_not_found")}</CommandEmpty>
							{[...groups].map(([kind, matches]) => (
								<CommandGroup key={kind} heading={headings[kind]}>
									{matches.map(({ entry, alias }) => (
										<CommandItem
											key={entry.key}
											value={`${kind}:${entry.key}`}
											onSelect={() => copy(entry.key)}
										>
											<span className="min-w-0 flex-1 truncate">
												{entry.label}
												{alias ? (
													<span className="ml-2 text-g-xs text-g-ink-mute">
														{alias}
													</span>
												) : null}
											</span>
											<code className="shrink-0 font-g-mono text-g-xs text-g-ink-mute">
												{entry.key}
											</code>
										</CommandItem>
									))}
								</CommandGroup>
							))}
						</CommandList>
					</Command>
				</DialogContent>
			</Dialog>
		</>
	);
}
