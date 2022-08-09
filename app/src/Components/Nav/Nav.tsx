import {
  Alignment,
  AnchorButton,
  Button,
  Classes,
  HTMLSelect,
  Icon,
  Navbar,
  NavbarDivider,
  NavbarGroup,
  NavbarHeading,
} from "@blueprintjs/core";
import { Link, useLocation } from "wouter";
import { Trans, useTranslation } from "react-i18next";
import { RootState, useAppSelector } from "~src/store";

export default function Nav() {
  let { t, i18n } = useTranslation();
  let language = i18n.language;

  const { user } = useAppSelector((state: RootState) => {
    return {
      user: state.user,
    };
  });

  const [location, _] = useLocation();
  return (
    <Navbar>
      <NavbarGroup align={Alignment.LEFT} className="w-full">
        <img
          src="/images/logo.png"
          className=" object-contain max-h-[75%] mt-auto mb-auto mr-1"
        />
        <Link href="/">
          <NavbarHeading>
            <a>gcsim web (beta)</a>
          </NavbarHeading>
        </Link>
        {location !== "/" ? (
          <>
            <NavbarDivider />
            <Link href="/simulator">
              <AnchorButton className={Classes.MINIMAL} icon="calculator">
                <span className="hidden md:block">
                  <Trans>nav.simulator</Trans>
                </span>
              </AnchorButton>
            </Link>
            <Link href="/viewer">
              <AnchorButton className={Classes.MINIMAL} icon="chart">
                <span className="hidden md:block">
                  <Trans>nav.viewer</Trans>
                </span>
              </AnchorButton>
            </Link>
            <Link href="/db">
              <AnchorButton className={Classes.MINIMAL} icon="database">
                <span className="hidden md:block">
                  <Trans>nav.teams_db</Trans>
                </span>
              </AnchorButton>
            </Link>
            <Link href="/about">
              <AnchorButton className={Classes.MINIMAL} icon="info-sign">
                <span className="hidden md:block">
                  <Trans>nav.about</Trans>
                </span>
              </AnchorButton>
            </Link>
          </>
        ) : null}
        <div className="ml-auto">
          <Link href="/account">
            <AnchorButton className={Classes.MINIMAL} icon="user">
              {user.user_name}
            </AnchorButton>
          </Link>
        </div>
        <div className="flex flex-row items-center ml-2">
          <HTMLSelect
            className="ml-2"
            value={language}
            onChange={(e) => {
              console.log(e.target.value);
              i18n.changeLanguage(e.target.value);
            }}
          >
            <option value="English">{t("nav.english")}</option>
            <option value="Chinese">{t("nav.chinese")}</option>
            <option value="German">{t("nav.german")}</option>
            <option value="Japanese">{t("nav.japanese")}</option>
            <option value="Spanish">{t("nav.spanish")}</option>
            <option value="Russian">{t("nav.russian")}</option>
          </HTMLSelect>
        </div>
      </NavbarGroup>
    </Navbar>
  );
}
