import Layout from "./Sectioning/layout";
import { Route, Switch } from "wouter";
import Database from "./Pages/Database/database";
// import { Dash } from "@gcsim/ui/src/Pages";
export default function App() {
  return (
    <Layout>
      <Switch>
        <Route path="/">
          <Database />
        </Route>
      </Switch>
    </Layout>
  );
}
