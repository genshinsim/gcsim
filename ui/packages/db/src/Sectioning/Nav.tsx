import { Alignment, HTMLSelect, Navbar } from "@blueprintjs/core";
import { useTranslation } from "react-i18next";
import { Link } from "wouter";
import logo from "./logo.png";

export default function Nav() {
  const { t, i18n } = useTranslation();

  const PageNavs = [
    // <NavLink key="database" href="/database" icon="database" text="" />,
    // <NavLink
    //   key="management"
    //   href="/management"
    //   icon="clipboard"
    //   text={t("nav.management")}
    // />,
  ];

  return (
    <Navbar>
      <div className="w-full 2xl:mx-auto 2xl:container">
        <Navbar.Group align={Alignment.LEFT}>
          <Navbar.Heading className="!mr-[10px]">
            <Link href="/" className="flex h-[50px] items-center">
              <img
                src={logo}
                className="object-scale-down max-h-[75%] m-auto mr-2"
              />
              <span className="font-medium font-mono">gcdatabase</span>
            </Link>
          </Navbar.Heading>
        </Navbar.Group>
        <Navbar.Group
          align={Alignment.LEFT}
          className=" min-[550px]:!flex !items-stretch"
        >
          <Navbar.Divider className="self-center" />
          {PageNavs}
        </Navbar.Group>
        <Navbar.Group align={Alignment.RIGHT}>
          {/* <Link href="/account">
              <Button minimal={true} icon="user" text="Guest" />
            </Link> */}
          <HTMLSelect
            value={i18n.resolvedLanguage}
            onChange={(e) => i18n.changeLanguage(e.target.value)}
          >
            <option value="en">{t<string>("nav.english")}</option>
            <option value="zh">{t<string>("nav.chinese")}</option>
            <option value="de">{t<string>("nav.german")}</option>
            <option value="ja">{t<string>("nav.japanese")}</option>
            <option value="es">{t<string>("nav.spanish")}</option>
            <option value="ru">{t<string>("nav.russian")}</option>
          </HTMLSelect>
        </Navbar.Group>
      </div>
    </Navbar>
  );
}

// type NavLinkProps = {
//   href: string;
//   icon: IconName | MaybeElement;
//   text: string;
// };

// const NavLink = ({ href, icon, text }: NavLinkProps) => {
//   return (
//     <Link href={href}>
//       <AnchorButton minimal={true} icon={icon}>
//         <span className="hidden min-[798px]:block">{text}</span>
//       </AnchorButton>
//     </Link>
//   );
// };
