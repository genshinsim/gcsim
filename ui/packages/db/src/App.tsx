import { Route, Switch } from "wouter";
import { Database } from "./Pages/Database";
import Layout from "./Sectioning/layout";
// import { Dash } from "@gcsim/ui/src/Pages";
export default function App() {
  return (
    <Layout>
      <Switch>
        <Route path="/">
          <Database />
        </Route>
        <Route path="/database">
          <Database />
        </Route>
        {/* <Route path="/tag/:tag">{({ tag }) => <TagDatabase tag={tag} />}</Route> */}

        {/* <Route path="/management">
          <Management />
        </Route> */}
      </Switch>
    </Layout>
  );
}
