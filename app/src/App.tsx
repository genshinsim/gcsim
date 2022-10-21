import React, { useRef } from 'react';
import { Redirect, Route, Switch, useLocation } from 'wouter';
import Footer from '/src/Components/Footer/Footer';
import Nav from '/src/Components/Nav/Nav';
import { Dash } from '/src/Pages/Dash';
import { Simple } from '/src/Pages/Sim';
import { SimWrapper } from './Pages/Sim/SimWrapper';
import { ViewerDash } from './Pages/ViewerDashboard';
import { DB } from './Pages/DB';
import './i18n';
import { Trans, useTranslation } from 'react-i18next';
import { DiscordCallback } from './PageUserAccount/DiscordCallback';
import { PageUserAccount } from './PageUserAccount';
import { ViewerLoader, ViewTypes } from './Pages/Viewer';
import { RedirectDB } from './Pages/DB/RedirectDB';

export default function App() {
  useTranslation();
  const content = useRef<HTMLDivElement>(null);
  const [location,] = useLocation();

  React.useEffect(() => {
    let loc = window.location.href;

    if (loc.includes('www.gcsim.app')) {
      loc = loc.replace('www.gcsim.app', 'gcsim.app');
      window.location.href = loc;
    }
  }, []);

  // every time you change location, scroll to top of page. This is necessary since the outer
  // content div will never rerender through the entire lifespan of the app and will always retain
  // its scroll position.
  React.useEffect(() => {
    content.current?.scrollTo(0, 0);
  }, [location]);

  return (
    <div className="bp4-dark h-screen flex flex-col">
      <Nav />
      <div ref={content} className="flex flex-col flex-auto overflow-y-scroll overflow-x-clip">
        <Switch>
          <Route path="/" component={Dash} />
          <Route path="/simple">
            <Redirect to="/simulator" />
          </Route>
          <Route path="/advanced">
            <Redirect to="/simulator" />
          </Route>
          <Route path="/simulator">
            <SimWrapper>
              <Simple />
            </SimWrapper>
          </Route>

          {/* Viewer Routes */}
          <Route path="/viewer">
            <ViewerLoader type={ViewTypes.Landing} />
          </Route>
          <Route path="/viewer/upload">
            <ViewerLoader type={ViewTypes.Upload} />
          </Route>
          <Route path="/viewer/web">
            <ViewerLoader type={ViewTypes.Web} />
          </Route>
          <Route path="/viewer/local">
            <ViewerLoader type={ViewTypes.Local} />
          </Route>
          <Route path="/viewer/share/:id">
            {(params) => <ViewerLoader type={ViewTypes.Share} id={params.id} />}
          </Route>

          {/* reroute v3 -> new viewer */}
          <Route path="/v3/viewer/share/:id">
            {(params) => <Redirect to={"/viewer/share/" + params.id} />}
          </Route>

          {/* Legacy Viewer Routes */}
          <Route path="/legacy/viewer">
            <ViewerDash path="/" />
          </Route>
          <Route path="/legacy/viewer/local">
            <ViewerDash path="local" />
          </Route>
          <Route path="/legacy/viewer/share/:id">
            {(params) => <ViewerDash path={params.id} next={true} />}
          </Route>

          {/* <Route path="/db">
            <Database />
          </Route>
          <Route path="/db/:avatar">
            {(params) => <CharacterView char={params.avatar} />}
          </Route>
          <Route path="/db/:avatar/:team">
            {(params) => <TeamsList char={params.avatar} team={params.team} />}
          </Route> */}
          <Route path="/db">
            <RedirectDB />
          </Route>
          <Route path="/old-db">
            <DB />
          </Route>
          {/* <Route path="/db/:char">{({ char }) => <DbChar char={char} />}</Route> */}
          <Route path="/account">
            <PageUserAccount />
          </Route>
          <Route path="/auth/discord">
            <DiscordCallback />
          </Route>
          <Route>
            <div className="m-2 text-center">
              <Trans>src.this_page_is</Trans>
            </div>
          </Route>
        </Switch>
        <div className="w-full pt-4 pb-4 md:pl-4">
          <Footer />
        </div>
      </div>
    </div>
  );
}
