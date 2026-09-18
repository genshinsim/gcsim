import {
	Button,
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuRadioGroup,
	DropdownMenuRadioItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
	Separator,
} from "@gcsim/primitives";
import type { ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { FaCalculator, FaDatabase, FaDiscord } from "react-icons/fa";
import { IoIosDocument, IoIosMenu } from "react-icons/io";
import { MdOutlineUpdate } from "react-icons/md";
import { Link } from "react-router-dom";
import logo from "./logo.png";

type NavLink = {
	key: string;
	icon: ReactNode;
	text: string;
} & ({ to: string } | { href: string });

function useNavLinks(): NavLink[] {
	const { t } = useTranslation();
	return [
		{
			key: "sim",
			to: "/simulator",
			icon: <FaCalculator />,
			text: t("nav.simulator"),
		},
		{
			key: "db",
			href: "https://simpact.app/",
			icon: <FaDatabase />,
			text: t("nav.teams_db"),
		},
		{
			key: "doc",
			href: "https://docs.gcsim.app",
			icon: <IoIosDocument size="24px" />,
			text: t("nav.documentation"),
		},
		{
			key: "update",
			href: "https://github.com/genshinsim/gcsim/releases",
			icon: <MdOutlineUpdate size="24px" />,
			text: t("nav.releases"),
		},
		{
			key: "discord",
			href: "https://discord.gg/m7jvjdxx7q",
			icon: <FaDiscord size="24px" />,
			text: "Discord",
		},
	];
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

// Returns the bare <Link>/<a> element (internal vs external) so it can be the
// single child of any radix `asChild` slot. Kept as a function, not a
// component, so the slot's ref/props reach the real anchor.
function renderNavLink(link: NavLink) {
	const inner = (
		<>
			{link.icon}
			<span>{link.text}</span>
		</>
	);
	return "to" in link ? (
		<Link to={link.to}>{inner}</Link>
	) : (
		<a href={link.href} target="_blank" rel="noreferrer">
			{inner}
		</a>
	);
}

const NavItem = ({ link }: { link: NavLink }) => (
	<Button asChild variant="ghost">
		{renderNavLink(link)}
	</Button>
);

const LanguageSelect = ({ className }: { className?: string }) => {
	const { t, i18n } = useTranslation();
	return (
		<Select
			value={i18n.resolvedLanguage}
			onValueChange={(value) => i18n.changeLanguage(value)}
		>
			<SelectTrigger className={className}>
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
};

const MobileMenu = () => {
	const { t, i18n } = useTranslation();
	const links = useNavLinks();
	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant="ghost" size="icon">
					<IoIosMenu size="24px" />
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="end" className="w-56">
				{links.map((link) => (
					<DropdownMenuItem key={link.key} asChild>
						{renderNavLink(link)}
					</DropdownMenuItem>
				))}
				<DropdownMenuSeparator />
				<DropdownMenuRadioGroup
					value={i18n.resolvedLanguage}
					onValueChange={(value) => i18n.changeLanguage(value)}
				>
					{LANGUAGES.map(({ value, key }) => (
						<DropdownMenuRadioItem key={value} value={value}>
							{t(key)}
						</DropdownMenuRadioItem>
					))}
				</DropdownMenuRadioGroup>
			</DropdownMenuContent>
		</DropdownMenu>
	);
};

export default () => {
	const links = useNavLinks();
	return (
		<nav className="bg-bp-header-color text-deprecated-foreground shadow-md">
			<div className="flex w-full 2xl:mx-auto 2xl:container">
				<div className="flex h-[50px] w-full items-center px-2">
					<Link to="/" className="mr-2.5 flex h-[50px] items-center">
						<img
							src={logo}
							alt=""
							className="m-auto mr-2 max-h-[75%] object-scale-down"
						/>
						<span className="font-medium font-mono">gcsim</span>
					</Link>

					<Separator
						orientation="vertical"
						className="mr-1 hidden !h-6 self-center min-[550px]:block"
					/>

					<div className="hidden items-center gap-0.5 min-[902px]:flex">
						{links.map((link) => (
							<NavItem key={link.key} link={link} />
						))}
					</div>

					<div className="ml-auto flex items-center">
						<LanguageSelect className="hidden min-[902px]:flex" />
						<div className="min-[902px]:hidden">
							<MobileMenu />
						</div>
					</div>
				</div>
			</div>
		</nav>
	);
};
