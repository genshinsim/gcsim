import React from "react";
import { Link } from "wouter";

export default function Nav() {
  return (
    <header className="text-xl">
      <div className="mx-auto flex flex-wrap p-2 flex-col md:flex-row mb-1 items-center">
        <Link href="/">
          <a className="flex title-font font-bond items-center mb-4 md:mb-0">
            <span className="ml-3 text-xl border-b-2 border-transparent">
              gcsim
            </span>
          </a>
        </Link>
        <div className="flex-grow justify-end flex flex-row">
          <nav className="md:ml-auto md:mr-1 md:py-1 md:pl-4   ">
            <a
              className="mr-5"
              href="https://github.com/genshinsim/gsimui/releases"
            >
              Install
            </a>
          </nav>
          <nav className="md:ml-1 md:mr-1 md:py-1 md:pl-4 ">
            <a className="mr-5" href="https://github.com/genshinsim/gsim/wiki">
              Docs
            </a>
          </nav>
          <nav className="md:ml-1 md:mr-1 md:py-1 md:pl-4  ">
            <Link href="/about">
              <a className="mr-5">About</a>
            </Link>
          </nav>
        </div>
      </div>
    </header>
  );
}
