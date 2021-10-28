import React from "react";
import { Route, Switch } from "wouter";
import Nav from "./components/Nav";
import Footer from "./components/Footer";
import About from "./components/About";
import SearchBar from "./components/SearchBar";
import Search from "./components/Search";
import Store from "./Store";
import "./App.css";

import fuseIndex from "./data/fuse-index.json";
import data from "./data/configs.json";
import Fuse from "fuse.js";

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
    <div className=" h-screen w-screen parent-bg">
      <div className="cover-image" />
      <div className="container mx-auto flex flex-col h-full">
        <Nav />
        <Switch>
          <Route path="/">
            <div className="flex-grow flex items-center justify-center">
              <SearchBar />
            </div>
          </Route>
          <Route path="/about" component={About} />
          <Route path="/search" component={Search} />
        </Switch>
        <Footer />
      </div>
    </div>
  );
}
