import {
	Button,
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
	Separator,
} from "@gcsim/primitives";
import { Link, type LinkProps } from "@tanstack/react-router";
import { Settings } from "lucide-react";
import type { ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { FaCalculator, FaDiscord } from "react-icons/fa";
import { IoIosDocument, IoIosMenu } from "react-icons/io";
import { MdOutlineUpdate } from "react-icons/md";
import { ExecutorDialog } from "./ExecutorDialog";
import logo from "./logo.png";

type NavLink = {
	key: string;
	icon: ReactNode;
	text: string;
} & ({ to: LinkProps["to"] } | { href: string });

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

const MobileMenu = () => {
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
			</DropdownMenuContent>
		</DropdownMenu>
	);
};

const SettingsButton = () => {
	const { t } = useTranslation();
	return (
		<Button asChild variant="ghost" size="icon">
			<Link to="/settings" aria-label={t("simple.settings")}>
				<Settings />
			</Link>
		</Button>
	);
};

export default () => {
	const links = useNavLinks();
	return (
		<nav className="bg-g-canvas text-g-ink border-b border-g-line-soft">
			<div className="flex w-full 2xl:mx-auto 2xl:container">
				<div className="flex h-[50px] w-full items-center px-2">
					<Link to="/" className="mr-2.5 flex h-[50px] items-center">
						<img
							src={logo}
							alt=""
							className="m-auto mr-2 max-h-[75%] object-scale-down"
						/>
						<span className="font-medium font-g-mono">gcsim</span>
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

					<div className="ml-auto flex items-center gap-1">
						<ExecutorDialog />
						<SettingsButton />
						<div className="min-[902px]:hidden">
							<MobileMenu />
						</div>
					</div>
				</div>
			</div>
		</nav>
	);
};
