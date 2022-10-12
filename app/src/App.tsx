import React from 'react';
import { Redirect, Route, Switch } from 'wouter';
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
import { CharacterView, Database, TeamsList } from '~src/PageDatabase';
import { PageUserAccount } from './PageUserAccount';
import { RedirectDB } from './Pages/DB/RedirectDB';

export default function App() {
  useTranslation();

  React.useEffect(() => {
    let loc = window.location.href;

    if (loc.includes('www.gcsim.app')) {
      loc = loc.replace('www.gcsim.app', 'gcsim.app');
      window.location.href = loc;
    }
  }, [window.location.href]);

  return (
    <div className="bp4-dark h-screen flex flex-col">
      <Nav />
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
        <Route path="/viewer/share/:id">
          {(params) => <ViewerDash path={params.id} />}
        </Route>
        <Route path="/v3/viewer/share/:id">
          {(params) => <ViewerDash path={params.id} next={true} />}
        </Route>
        <Route path="/viewer/local">
          <ViewerDash path="local" />
        </Route>
        <Route path="/viewer">
          <ViewerDash path="/" />
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
  );
}
