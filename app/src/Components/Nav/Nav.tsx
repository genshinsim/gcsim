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
              text="Simulator"
              onClick={() => setLocation("/simple")}
            />
            <Button
              className={Classes.MINIMAL}
              icon="rocket-slant"
              text="Advanced"
              onClick={() => setLocation("/simple")}
            />
            <Button
              className={Classes.MINIMAL}
              icon="chart"
              text="Viewer"
              onClick={() => setLocation("/viewer")}
            />
            <Button
              className={Classes.MINIMAL}
              icon="database"
              text="Action List DB"
              onClick={() => setLocation("/db")}
            />
            <Button
              className={Classes.MINIMAL}
              icon="lightbulb"
              text="About"
              onClick={() => setLocation("/about")}
            />
          </>
        ) : null}
      </NavbarGroup>
      <NavbarGroup align={Alignment.RIGHT}></NavbarGroup>
    </Navbar>
  );
}
