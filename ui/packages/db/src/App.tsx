import Layout from "./Sectioning/layout";
import { Route, Switch } from "wouter";
import Leaderboard from "./Pages/Leaderboard/leaderboard";
import Database from "./Pages/Database/database";
// import { Dash } from "@gcsim/ui/src/Pages";
export default function App() {
  return (
    <Layout>
      <Switch>
        <Route path="/" component={Dash} />
        <Route path="/leaderboard">
          <Leaderboard />
        </Route>
        <Route path="/database">
          <Database />
        </Route>
      </Switch>
    </Layout>
  );
}

function Dash() {
  return (
    <div className="flex justify-center h-4/6 my-40">
      ... empty, like my heart
    </div>
  );
}
