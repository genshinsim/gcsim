import React from "react";
import { Link, Route, Switch } from "wouter";
import Fuse from "fuse.js";

import "./App.css";
import Store from "./Store";
import Nav from "./components/Nav";
import Footer from "./components/Footer";
import About from "./components/About";
import Home from "./components/Home";
import Browse from "./components/Browse";
import SearchResults from "./components/SearchResults";
import ActionDB from "./components/ActionDB";
import GetStarted from "./components/GetStarted";

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
          <Route path="/" component={Home} />
          <Route path="/getting-started" component={GetStarted} />
          <Route path="/browse" component={Browse} />
          <Route path="/about" component={About} />
          <Route path="/db" component={ActionDB} />
          <Route path="/db/results" component={SearchResults} />
        </Switch>
        <Footer />
      </div>
    </div>
  );
}
