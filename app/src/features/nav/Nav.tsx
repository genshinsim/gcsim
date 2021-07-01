import { Button, Navbar } from "@blueprintjs/core";
import { RootState } from "app/store";

import React from "react";
import { useSelector } from "react-redux";

import { useLocation } from "wouter";

function Nav() {
  const { haveResults, haveDebug } = useSelector((state: RootState) => {
    return {
      haveResults: state.results.haveResult,
      haveDebug: state.debug.haveDebug,
    };
  });

  const [, setLocation] = useLocation();

  const navigate = (url: string) => {
    return () => {
      setLocation(url);
    };
  };

  return (
    <Navbar style={{ marginBottom: "5px" }}>
      <Navbar.Group>
        <Navbar.Heading>gisim</Navbar.Heading>
        <Navbar.Divider />
        <Button
          className="bp3-minimal"
          icon="home"
          text="Sim"
          onClick={navigate("/")}
        />
        <Button
          className="bp3-minimal"
          icon="document"
          text="Results"
          onClick={navigate("/results")}
          disabled={!haveResults}
        />
        <Button
          className="bp3-minimal"
          icon="error"
          text="Debug"
          onClick={navigate("/debug")}
          disabled={!haveDebug}
        />
      </Navbar.Group>
    </Navbar>
  );
}

export default Nav;
