import {
  Alignment,
  AnchorButton,
  Button,
  Classes,
  Collapse,
  HTMLSelect,
  Icon,
  IconName,
  Navbar,
} from "@blueprintjs/core";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { FaDiscord } from "react-icons/fa";
import { IoIosDocument, IoIosMenu } from "react-icons/io";
import { MdOutlineUpdate } from "react-icons/md";
import { Link } from "react-router-dom";
import { useAppSelector } from "../Stores/store";
import logo from "./logo.png";

export default ({}) => {
  const { t, i18n } = useTranslation();
  const user = useAppSelector((state) => state.user);
  const [isOpen, setIsOpen] = useState(false);

  const PageNavs = [
    <NavButton
      key="sim"
      href="/simulator"
      icon="calculator"
      text={t("nav.simulator")}
    />,
    <AnchorButton
      key="db"
      className={Classes.MINIMAL}
      icon="database"
      href="https://simpact.app/"
      target="_blank"
    >
      <span>{t<string>("nav.teams_db")}</span>
    </AnchorButton>,
    <AnchorButton
      key="doc"
      className={Classes.MINIMAL}
      icon={<IoIosDocument size="24px" color="#abb3bf" />}
      href="https://docs.gcsim.app"
      target="_blank"
    >
      <span>Documentation</span>
    </AnchorButton>,
    <AnchorButton
      key="update"
      className={Classes.MINIMAL}
      icon={<MdOutlineUpdate size="24px" color="#abb3bf" />}
      href="https://github.com/genshinsim/gcsim/releases"
      target="_blank"
    >
      <span>Releases</span>
    </AnchorButton>,
    <AnchorButton
      key="discord"
      className={Classes.MINIMAL}
      href="https://discord.gg/m7jvjdxx7q"
      target="_blank"
      rel="noreferrer"
      icon={<FaDiscord size="24px" color="#abb3bf" />}
    >
      <span>{"Discord"}</span>
    </AnchorButton>,
  ];

  return (
    <>
      <Navbar>
        <div className="flex flex-col w-full 2xl:mx-auto 2xl:container">
          <div>
            <Navbar.Group align={Alignment.LEFT}>
              <Navbar.Heading className="!mr-[10px]">
                <Link to="/" className="flex h-[50px] items-center">
                  <img
                    src={logo}
                    className="object-scale-down max-h-[75%] m-auto mr-2"
                  />
                  <span className="font-medium font-mono">gcsim</span>
                </Link>
              </Navbar.Heading>
            </Navbar.Group>
            <Navbar.Group
              align={Alignment.LEFT}
              className="!hidden min-[902px]:!flex !items-stretch"
            >
              <Navbar.Divider className="!hidden min-[550px]:!flex self-center" />
              {PageNavs}
            </Navbar.Group>
            <Navbar.Group
              align={Alignment.RIGHT}
              className="!hidden min-[902px]:!flex !items-stretch"
            >
              {/* <NavButton href="/account" icon="user" text={user.uid === "" ? "Guest" : user.name} /> */}
              <HTMLSelect
                className="ml-2 self-center"
                value={i18n.resolvedLanguage}
                onChange={(e) => i18n.changeLanguage(e.target.value)}
              >
                <option value="en">{t<string>("nav.english")}</option>
                <option value="zh">{t<string>("nav.chinese")}</option>
                <option value="ja">{t<string>("nav.japanese")}</option>
                <option value="ko">{t<string>("nav.korean")}</option>
                <option value="es">{t<string>("nav.spanish")}</option>
                <option value="ru">{t<string>("nav.russian")}</option>
              </HTMLSelect>
            </Navbar.Group>
            <Navbar.Group
              align={Alignment.RIGHT}
              className="!flex min-[902px]:!hidden !items-stretch"
            >
              <Button
                className={Classes.MINIMAL}
                icon={<IoIosMenu size="24px" color="#abb3bf" />}
                onClick={() => setIsOpen(!isOpen)}
                outlined={false}
              />
            </Navbar.Group>
          </div>
        </div>
      </Navbar>
      <Collapse isOpen={isOpen}>
        {/* Navbar's CSS will overwrite height and padding unless marked as important */}
        <Navbar className="!h-full !py-2 flex flex-col gap-0.5">
          {PageNavs}
          <HTMLSelect
            className="my-1 ml-2 self-center"
            value={i18n.resolvedLanguage}
            onChange={(e) => i18n.changeLanguage(e.target.value)}
          >
            <option value="en">{t<string>("nav.english")}</option>
            <option value="zh">{t<string>("nav.chinese")}</option>
            <option value="ja">{t<string>("nav.japanese")}</option>
            <option value="ko">{t<string>("nav.korean")}</option>
            <option value="es">{t<string>("nav.spanish")}</option>
            <option value="ru">{t<string>("nav.russian")}</option>
          </HTMLSelect>
        </Navbar>
      </Collapse>
    </>
  );
};

type NavButtonProps = {
  href: string;
  icon: IconName;
  text: string;
};

const NavButton = ({ href, icon, text }: NavButtonProps) => {
  return (
    <Link
      to={href}
      role="button"
      className="bp4-button bp4-minimal"
      tabIndex={0}
    >
      <Icon icon={icon} />
      <span className="bp4-button-text">{text}</span>
    </Link>
  );
};
