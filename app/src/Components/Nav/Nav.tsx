import {
  Alignment,
  Button,
  Classes,
  H5,
  Navbar,
  NavbarDivider,
  NavbarGroup,
  NavbarHeading,
  Switch,
} from "@blueprintjs/core";
import { useLocation } from "wouter";

export default function Nav() {
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
              <span className="xs:hidden md:block">Simulator</span>
            </Button>
            <Button
              className={Classes.MINIMAL}
              icon="rocket-slant"
              onClick={() => setLocation("/advanced")}
            >
              <span className="xs:hidden md:block">Advanced</span>
            </Button>
            <Button
              className={Classes.MINIMAL}
              icon="chart"
              onClick={() => setLocation("/viewer")}
            >
              <span className="xs:hidden md:block">Viewer</span>
            </Button>
            <Button
              className={Classes.MINIMAL}
              icon="database"
              onClick={() => setLocation("/db")}
            >
              <span className="xs:hidden md:block">Action List DB</span>
            </Button>
            <Button
              className={Classes.MINIMAL}
              icon="info-sign"
              onClick={() => setLocation("/about")}
            >
              <span className="xs:hidden md:block">About</span>
            </Button>
          </>
        ) : null}
      </NavbarGroup>
      <NavbarGroup align={Alignment.RIGHT}></NavbarGroup>
    </Navbar>
  );
}
