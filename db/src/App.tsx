import React from "react";
import { Database } from "Pages/Database";
import { Route, Switch } from "wouter";
import { Teams, SimByTeam } from "Pages/Teams";
import { Navbar } from "Components/Nav/Navbar";

function App() {
  return (
    <div className="App">
      <Navbar />
      <main className="flex flex-col h-full m-2 w-full xs:w-full sm:w-[640px] hd:w-full wide:w-[1160px] ml-auto mr-auto ">
        <div className="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative">
          <strong className="font-bold">Warning: </strong>
          <span className="block sm:inline">
            The purpose of this database is to aid users in writing their own configs by providing samples.
            By no means does it claim to be a DPS leaderboard or follow any standards.
            The UI is also a work in progress, expect things to be broken.
          </span>
        </div>
        <Switch>
          <Route path="/" component={Database} />
          <Route path="/db/:avatar">
            {(params) => <Teams char={params.avatar} />}
          </Route>
          <Route path="/db/:avatar/:team">
            {(params) => <SimByTeam char={params.avatar} team={params.team} />}
          </Route>
        </Switch>
      </main>
    </div>
  );
}

export default App;
