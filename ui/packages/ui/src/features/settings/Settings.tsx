import {
	cn,
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@gcsim/primitives";
import type { ComponentProps, ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { usePrefs } from "../../stores/AppState";
import { THEMES } from "../../stores/themes";
import { HERO_IMAGES } from "../dashboard/heroImages";
import { SectionHead } from "../dashboard/SectionHead";

export function Settings() {
	const { t } = useTranslation();

	return (
		<main className="w-full flex-grow bg-g-canvas text-g-ink">
			<div className="mx-auto flex w-full max-w-[960px] flex-col gap-g-section px-4 py-8 sm:px-g-page">
				<div>
					<h1 className="mb-1 font-g-display text-g-h1 font-bold">
						{t("simple.settings")}
					</h1>
					<p className="text-g-body text-g-ink-dim">{t("settings.subtitle")}</p>
				</div>

				<Section
					title={t("settings.language")}
					subtitle={t("settings.language_sub")}
				>
					<LanguageSelect />
				</Section>

				<Section title={t("settings.theme")} subtitle={t("settings.theme_sub")}>
					<ThemePicker />
				</Section>

				<Section title={t("settings.hero")} subtitle={t("settings.hero_sub")}>
					<HeroPicker />
				</Section>
			</div>
		</main>
	);
}

function Section({
	title,
	subtitle,
	children,
}: {
	title: string;
	subtitle: string;
	children: ReactNode;
}) {
	return (
		<section className="flex flex-col gap-4">
			<SectionHead title={title} subtitle={subtitle} />
			{children}
		</section>
	);
}

const LANGUAGES = [
	{ value: "en", key: "nav.english" },
	{ value: "zh", key: "nav.chinese" },
	{ value: "ja", key: "nav.japanese" },
	{ value: "ko", key: "nav.korean" },
	{ value: "es", key: "nav.spanish" },
	{ value: "ru", key: "nav.russian" },
	{ value: "de", key: "nav.german" },
] as const;

function LanguageSelect() {
	const { t, i18n } = useTranslation();
	return (
		<Select
			value={i18n.resolvedLanguage}
			onValueChange={(value) => i18n.changeLanguage(value)}
		>
			<SelectTrigger aria-label={t("settings.language")} className="w-56">
				<SelectValue />
			</SelectTrigger>
			<SelectContent>
				{LANGUAGES.map(({ value, key }) => (
					<SelectItem key={value} value={value}>
						{t(key)}
					</SelectItem>
				))}
			</SelectContent>
		</Select>
	);
}

function Choice({
	name,
	checked,
	onSelect,
	className,
	children,
	...rest
}: ComponentProps<"label"> & {
	name: string;
	checked: boolean;
	onSelect: () => void;
}) {
	return (
		<label
			{...rest}
			className={cn(
				"cursor-pointer overflow-hidden rounded-g-lg border-2 transition-colors has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-g-accent",
				checked ? "border-g-accent" : "border-g-line hover:border-g-ink-mute",
				className,
			)}
		>
			<input
				type="radio"
				name={name}
				checked={checked}
				onChange={onSelect}
				className="sr-only"
			/>
			{children}
		</label>
	);
}

function ThemePicker() {
	const { t } = useTranslation();
	const { theme, setTheme } = usePrefs();

	return (
		<fieldset className="grid grid-cols-2 gap-g-base-lg sm:grid-cols-3 md:grid-cols-5">
			<legend className="sr-only">{t("settings.theme")}</legend>
			{THEMES.map((option) => (
				<Choice
					key={option.id}
					name="theme"
					checked={theme === option.id}
					onSelect={() => setTheme(option.id)}
					data-theme={option.id}
					className="bg-g-canvas text-g-ink"
				>
					<div className="flex gap-1 p-2">
						<span className="h-6 flex-1 rounded-g-sm bg-g-surface-2" />
						<span className="size-6 rounded-g-sm bg-g-accent" />
						<span className="size-6 rounded-g-sm bg-g-pyro" />
						<span className="size-6 rounded-g-sm bg-g-hydro" />
					</div>
					<div className="flex items-baseline justify-between gap-2 bg-g-surface px-2 py-1.5">
						<span className="font-medium text-g-sm">{option.name}</span>
						<span className="text-g-xs text-g-ink-mute">
							{t(option.light ? "settings.light" : "settings.dark")}
						</span>
					</div>
				</Choice>
			))}
		</fieldset>
	);
}

function HeroPicker() {
	const { t } = useTranslation();
	const { hero, setHero } = usePrefs();

	return (
		<fieldset className="grid grid-cols-2 gap-g-base-lg sm:grid-cols-3">
			<legend className="sr-only">{t("settings.hero")}</legend>
			{HERO_IMAGES.map((option) => (
				<Choice
					key={option.id}
					name="hero"
					checked={hero === option.id}
					onSelect={() => setHero(option.id)}
					className="bg-g-surface"
				>
					<img
						src={option.src}
						alt=""
						className="aspect-video w-full object-cover"
					/>
					<div className="px-2 py-1.5 font-medium text-g-sm">{option.name}</div>
				</Choice>
			))}
		</fieldset>
	);
}
