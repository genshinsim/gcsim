import { Route, Switch } from "wouter";
import Footer from "/src/Components/Footer/Footer";
import Nav from "/src/Components/Nav/Nav";
import { Dash } from "/src/Pages/Dash";
import { Simple, Advanced } from "/src/Pages/Sim";
import { SimWrapper } from "./Pages/Sim/SimWrapper";
import { ViewerDash } from "./Pages/ViewerDashboard";
import { DB } from "./Pages/DB";

export default function App() {
  return (
    <div className=".bp3-dark h-screen flex flex-col">
      <Nav />
      <Switch>
        <Route path="/" component={Dash} />
        <Route path="/simple">
          <SimWrapper>
            <Simple />
          </SimWrapper>
        </Route>
        <Route path="/advanced">
          <SimWrapper>
            <Advanced />
          </SimWrapper>
        </Route>
        <Route path="/viewer/share/:id">
          {(params) => <ViewerDash path={params.id} />}
        </Route>
        <Route path="/viewer">
          <ViewerDash path="/" />
        </Route>
        <Route path="/db">
          <DB />
        </Route>
        <Route>
          <div className="m-2 text-center">
            This page is not implemented yet. Stay tuned!
          </div>
        </Route>
      </Switch>
      <div className="w-full pt-4 pb-4 md:pl-4">
        <Footer />
      </div>
    </div>
  );
}
