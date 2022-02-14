import { Route, Switch } from "wouter";
import Footer from "/src/Components/Footer/Footer";
import Nav from "/src/Components/Nav/Nav";
import { Dash } from "/src/Pages/Dash";
import Shared from "/src/Shared";
import { Simple } from "/src/Pages/Sim";
import { SimWrapper } from "./Pages/Sim/SimWrapper";
import { Viewer } from "./Pages/Viewer/Viewer";

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
        <Route path="/share/:id">
          {(params) => <Shared path={params.id} />}
        </Route>
        <Route path="/viewer">
          <Viewer data="{}" handleClose={() => { return false }} />
        </Route>
      </Switch>
      <div className="w-full pt-4 pb-4 md:pl-4">
        <Footer />
      </div>
    </div>
  );
}
