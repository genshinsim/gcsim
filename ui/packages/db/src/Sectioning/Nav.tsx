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
import { useTranslation } from "react-i18next";
import { Link } from "wouter";
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

export default function Nav() {
	const { t, i18n } = useTranslation();

	return (
		<Navbar className="h-[50px]">
			<div className="flex w-full 2xl:mx-auto 2xl:container">
				<NavbarHeading className="!mr-[10px]">
					<Link href="/" className="flex h-[50px] items-center">
						<img
							src={logo}
							alt=""
							className="object-scale-down max-h-[75%] m-auto mr-2"
						/>
						<span className="font-medium font-mono">simpact</span>
					</Link>
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
	);
}
