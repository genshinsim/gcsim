import React from "react";
import { Database } from "Pages/Database";
import { Route, Switch } from "wouter";
import { Teams, SimByTeam } from "Pages/Teams";
import { Navbar } from "Components/Nav/Navbar";
import { AllSims } from "Pages/Database/All";

function App() {
  return (
    <div className="App">
      <Navbar />
      <Switch>
        <Route path="/" component={Database} />
        <Route path="/all" component={AllSims} />
        <Route path="/db/:avatar">
          {(params) => <Teams char={params.avatar} />}
        </Route>
        <Route path="/db/:avatar/:team">
          {(params) => <SimByTeam char={params.avatar} team={params.team} />}
        </Route>
      </Switch>
    </div>
  );
}

export default App;
