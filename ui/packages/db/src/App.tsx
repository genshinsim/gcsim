import { Database } from "./Pages/Database/Database";
import Layout from "./Sectioning/layout";
import { Route, Switch } from "wouter";
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
