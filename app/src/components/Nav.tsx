import React from "react";
import { Link } from "wouter";
import { AppContext } from "../Store";

export default function Nav() {
  const { dispatch } = React.useContext(AppContext);

  return (
    <header className="text-xl">
      <div className="mx-auto flex flex-wrap p-2 flex-col md:flex-row mb-1 items-center">
        <Link
          href="/"
          onClick={() => {
            dispatch({ type: "set", str: "" });
          }}
        >
          <a
            className="flex title-font font-bond items-center mb-4 md:mb-0"
            href="/"
          >
            <span className="ml-3 text-3xl border-b-2 border-transparent">
              gcsim
            </span>
          </a>
        </Link>
        <div className="flex-grow justify-end flex flex-row">
          <nav className="md:ml-auto md:mr-1 md:py-1 md:pl-4">
            <Link href="/getting-started">
              <a className="mr-5" href="/getting-started">
                Get Started
              </a>
            </Link>
          </nav>
          <nav className="md:ml-1 md:mr-1 md:py-1 md:pl-4 ">
            <Link href="/db">
              <a className="mr-5" href="/db">
                Action Lists
              </a>
            </Link>
          </nav>
        </div>
      </div>
    </header>
  );
}
