import {
	Navbar,
	NavbarDivider,
	NavbarGroup,
	NavbarHeading,
} from "@gcsim/components";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@gcsim/primitives";
import type { ReactNode } from "react";
import { Trans, useTranslation } from "react-i18next";
import logo from "./logo.png";

const LANGUAGES = [
	{ value: "en", key: "nav.english" },
	{ value: "zh", key: "nav.chinese" },
	{ value: "ja", key: "nav.japanese" },
	{ value: "ko", key: "nav.korean" },
	{ value: "es", key: "nav.spanish" },
	{ value: "ru", key: "nav.russian" },
	{ value: "de", key: "nav.german" },
] as const;

const TextLink = ({
	href,
	newTab,
	children,
}: {
	href: string;
	newTab?: boolean;
	children?: ReactNode;
}) => (
	<a
		href={href}
		target={newTab ? "_blank" : undefined}
		rel={newTab ? "noreferrer" : undefined}
		className="font-semibold text-g-accent underline hover:text-g-ink"
	>
		{children}
	</a>
);

export default function App() {
	const { t, i18n } = useTranslation();

	return (
		<>
			<Navbar className="h-[50px]">
				<div className="mx-auto flex w-full max-w-[1160px] px-8">
					<NavbarHeading className="!mr-[10px]">
						<a href="/" className="flex h-[50px] items-center">
							<img
								src={logo}
								alt=""
								className="object-scale-down max-h-[75%] m-auto mr-2"
							/>
							<span className="font-medium font-mono">simpact</span>
						</a>
					</NavbarHeading>
					<NavbarGroup className="min-[550px]:flex items-stretch">
						<NavbarDivider />
					</NavbarGroup>
					<NavbarGroup align="end">
						<Select
							value={i18n.resolvedLanguage}
							onValueChange={(value) => i18n.changeLanguage(value)}
						>
							<SelectTrigger>
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
					</NavbarGroup>
				</div>
			</Navbar>
			<div className="mx-auto flex max-w-[1160px] flex-col p-g-page">
				<section className="flex flex-col gap-g-base-lg rounded-g-xl border border-g-line-soft bg-g-surface p-8">
					<h1 className="font-g-display text-g-h1 font-bold text-g-ink md:text-g-hero">
						{t("db.home.archived_title")}
					</h1>
					<div className="flex flex-col gap-g-base text-g-lg text-g-ink-dim">
						<p>
							<Trans
								i18nKey="db.home.archived_body"
								components={{ url: <TextLink href="https://db.kqm.gg/" /> }}
							/>
						</p>
						<p>
							<Trans
								i18nKey="db.home.archived_contact"
								components={{
									url: <TextLink href="https://discord.gg/keqing" newTab />,
								}}
							/>
						</p>
					</div>
				</section>
			</div>
		</>
	);
}
