import {
  Alignment,
  Button,
  Classes,
  H5,
  HTMLSelect,
  Navbar,
  NavbarDivider,
  NavbarGroup,
  NavbarHeading,
  Switch,
} from "@blueprintjs/core";
import { useLocation } from "wouter";
import { Trans, useTranslation } from "react-i18next";

export default function Nav() {
  let { t, i18n } = useTranslation();
  let language = i18n.language;

  const [location, setLocation] = useLocation();
  return (
    <Navbar>
      <NavbarGroup align={Alignment.LEFT}>
        <NavbarHeading
          onClick={() => setLocation("/")}
          className="hover:cursor-pointer"
        >
          gcsim web (beta)
        </NavbarHeading>

        {location !== "/" ? (
          <>
            <NavbarDivider />
            <Button
              className={Classes.MINIMAL}
              icon="calculator"
              onClick={() => setLocation("/simple")}
            >
              <span className="xs:hidden md:block">
                <Trans>nav.simulator</Trans>
              </span>
            </Button>
            <Button
              className={Classes.MINIMAL}
              icon="rocket-slant"
              onClick={() => setLocation("/advanced")}
            >
              <span className="xs:hidden md:block">
                <Trans>nav.advanced</Trans>
              </span>
            </Button>
            <Button
              className={Classes.MINIMAL}
              icon="chart"
              onClick={() => setLocation("/viewer")}
            >
              <span className="xs:hidden md:block">
                <Trans>nav.viewer</Trans>
              </span>
            </Button>
            <Button
              className={Classes.MINIMAL}
              icon="database"
              onClick={() => setLocation("/db")}
            >
              <span className="xs:hidden md:block">
                <Trans>nav.teams_db</Trans>
              </span>
            </Button>
            <Button
              className={Classes.MINIMAL}
              icon="info-sign"
              onClick={() => setLocation("/about")}
            >
              <span className="xs:hidden md:block">
                <Trans>nav.about</Trans>
              </span>
            </Button>
          </>
        ) : null}
      </NavbarGroup>
      <NavbarGroup align={Alignment.RIGHT}>
        <Trans>nav.language</Trans>
        <HTMLSelect
          className="ml-2"
          value={language}
          onChange={(e) => i18n.changeLanguage(e.target.value)}
        >
          <option value="English">{t("nav.english")}</option>
          <option value="Chinese">{t("nav.chinese")}</option>
          <option value="Japanese">{t("nav.japanese")}</option>
          <option value="Spanish">{t("nav.spanish")}</option>
        </HTMLSelect>
      </NavbarGroup>
    </Navbar>
  );
}
