import React from "react";
import { Link, Route, Switch } from "wouter";
import Fuse from "fuse.js";

import "./App.css";
import Store from "./Store";
import Nav from "./components/Nav";
import Footer from "./components/Footer";
import About from "./components/About";
import SearchBar from "./components/SearchBar";
import Browse from "./components/Browse";
import Search from "./components/Search";

import fuseIndex from "./data/fuse-index.json";
import data from "./data/configs.json";

const index = Fuse.parseIndex(fuseIndex);
export const fuse = new Fuse(
  data,
  { keys: ["title", "author", "description", "characters"], threshold: 0.4 },
  index
);

export default function AppWrapper(): JSX.Element {
  // Store, renders the provider, so the context will be accessible from App.
  return (
    <Store>
      <App />
    </Store>
  );
}

function App() {
  return (
    <div className="h-screen">
      <div className="container mx-auto flex flex-col h-full">
        <Nav />
        <Switch>
          <Route path="/">
            <div className="flex-grow flex flex-col items-center justify-center">
              <SearchBar />
              <div className="p-2">
                <Link href="/browse">
                  <a href="/browse">Browse All</a>
                </Link>
              </div>
            </div>
          </Route>
          <Route path="/browse" component={Browse} />
          <Route path="/about" component={About} />
          <Route path="/search" component={Search} />
        </Switch>
        <Footer />
      </div>
    </div>
  );
}
