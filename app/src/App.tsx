import React from "react";
import DiscordLogo from "./discord-icon.svg";
import GithubLogo from "./github-icon.svg";
import SearchButton from "./search_white_24dp.svg";
import { Link, Route, Switch } from "wouter";
import "./App.css";

function Nav() {
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

function Footer() {
  return (
    <div className="flex flex-row w-full justify-end gap-2 lg:mb-10 md:m-3 sm:m-1">
      <div className=" hover:bg-gray-600 p-2 rounded-md h-12">
        <a href="https://github.com/genshinsim/gsim" target="_blank">
          <img
            src={GithubLogo}
            alt="Github Logo"
            className="object-contain h-full"
          />
        </a>
      </div>
      <div className=" hover:bg-gray-600 p-2 rounded-md h-12 ">
        <a href="https://discord.gg/m7jvjdxx7q" target="_blank">
          <img
            src={DiscordLogo}
            alt="Discord Logo"
            className="object-contain h-full"
          />
        </a>
      </div>
    </div>
  );
}

function Search() {
  return (
    <div className="flex-grow flex items-center justify-center">
      <div
        className=" mb-10 md:w-1/2 lg:w-2/3 sm:w-full bg-gray-500 rounded-lg p-3 m-1"
        style={{ position: "relative" }}
      >
        <form
          onSubmit={(e) => {
            e.preventDefault();
            console.log("hi");
          }}
        >
          <input
            type="text"
            className="p-1 w-full text-2xl md:text-xl bg-gray-500 outline-none"
            placeholder="Search for an action list..."
          />
        </form>
        <img
          src={SearchButton}
          alt="search"
          className="rounded-md hover:bg-gray-400"
          style={{
            objectFit: "contain",
            position: "absolute",
            right: 0,
            top: 0,
            height: "100%",
            padding: "0.5rem",
          }}
        />
      </div>
    </div>
  );
}

function About() {
  return (
    <div className="flex-grow flex p-10">
      Hi! Welcome to gcsim's website. gsim is a monte carlo simulation tool used
      to model team dps in Genshin Impact. For more information visit us on
      GitHub.
    </div>
  );
}

function App() {
  return (
    <div className="bg-gradient-to-br from-gray-700 via-gray-800 to-gray-800 text-white h-screen w-screen parent-bg">
      <div className="cover-image" />
      <div className="container mx-auto flex flex-col h-full">
        <Nav />
        <Switch>
          <Route path="/about">
            <About />
          </Route>
          <Route>
            <Search />
          </Route>
        </Switch>
        <Footer />
      </div>
    </div>
  );
}

export default App;
