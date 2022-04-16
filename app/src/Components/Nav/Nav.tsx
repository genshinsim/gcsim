import {
  Alignment,
  AnchorButton,
  Classes,
  HTMLSelect,
  Navbar,
  NavbarDivider,
  NavbarGroup,
  NavbarHeading,
} from "@blueprintjs/core";
import { Link, useLocation } from "wouter";
import { Trans, useTranslation } from "react-i18next";

export default function Nav() {
  let { t, i18n } = useTranslation();
  let language = i18n.language;

  const [location, _] = useLocation();
  return (
    <Navbar>
      <NavbarGroup align={Alignment.LEFT} className="w-full">
        <Link href="/">
          <NavbarHeading>
            <a>
              gcsim web (beta)
            </a>
          </NavbarHeading>
        </Link>
        {location !== "/" ? (
          <>
            <NavbarDivider />
            <Link href="/simulator">
              <AnchorButton
                className={Classes.MINIMAL}
                icon="calculator"
              >
                <span className="hidden md:block">
                  <Trans>nav.simulator</Trans>
                </span>
              </AnchorButton>
            </Link>
            <Link href="/viewer">
              <AnchorButton
                className={Classes.MINIMAL}
                icon="chart"
              >
                <span className="hidden md:block">
                  <Trans>nav.viewer</Trans>
                </span>
              </AnchorButton>
            </Link>
            <Link href="/db">
              <AnchorButton
                className={Classes.MINIMAL}
                icon="database"
              >
                <span className="hidden md:block">
                  <Trans>nav.teams_db</Trans>
                </span>
              </AnchorButton>
            </Link>
            <Link href="/about">
              <AnchorButton
                className={Classes.MINIMAL}
                icon="info-sign"
              >
                <span className="hidden md:block">
                  <Trans>nav.about</Trans>
                </span>
              </AnchorButton>
            </Link>
          </>
        ) : null}
        <div className="flex flex-row items-center ml-auto">
          <span className="hidden lg:block">
            <Trans>nav.language</Trans>
          </span>
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
