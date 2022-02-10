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

export default function Nav() {
    return (
        <Navbar>
            <NavbarGroup align={Alignment.LEFT}>
                <NavbarHeading>gcsim web<span className="align-super text-xs">beta</span></NavbarHeading>
                <NavbarDivider />
                        <Button className={Classes.MINIMAL} icon="home" text="Home" />
                        <Button className={Classes.MINIMAL} icon="lightbulb" text="About" />
            </NavbarGroup>

        </Navbar>
    )
}