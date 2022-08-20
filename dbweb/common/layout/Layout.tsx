import React from "react";
import { NavBar } from "../nav/NavBar";

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <NavBar />
      <main className="relative">{children}</main>
    </>
  );
}
