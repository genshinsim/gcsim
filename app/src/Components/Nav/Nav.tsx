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
  const [_, setLocation] = useLocation();
  return (
    <Navbar>
      <NavbarGroup align={Alignment.LEFT}>
        <NavbarHeading>gcsim web (beta)</NavbarHeading>
      </NavbarGroup>
      <NavbarGroup align={Alignment.RIGHT}>
        <Button
          className={Classes.MINIMAL}
          icon="home"
          text="Home"
          onClick={() => setLocation("/")}
        />
        <Button
          className={Classes.MINIMAL}
          icon="lightbulb"
          text="About"
          onClick={() => setLocation("/about")}
        />
      </NavbarGroup>
    </Navbar>
  );
}
