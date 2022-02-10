import { Route, Switch } from "wouter";
import Dash from "./Dash";
import Shared from "./Shared";

export default function App() {
  return (
    <div className=".bp3-dark mx-auto h-full">
      <Switch>
        <Route path="/" component={Dash} />
        <Route path="/share/:id">
          {(params) => <Shared path={params.id} />}
        </Route>
      </Switch>
    </div>
  );
}
