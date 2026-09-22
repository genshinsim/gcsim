import { Home } from "Pages/Home";
import { useEffect } from "react";
import { Route, Switch, useLocation } from "wouter";
import { Database } from "./Pages/Database";
import Layout from "./Sectioning/layout";
export default function App() {
	const [location] = useLocation();

	// wouter's client-side nav preserves window scroll, so the listing can mount at the bottom and
	// make react-infinite-scroll fetch pages forever. Reset scroll on every route change.
	// biome-ignore lint/correctness/useExhaustiveDependencies: location is the trigger, not a value used in the effect body
	useEffect(() => {
		window.scrollTo(0, 0);
	}, [location]);

	return (
		<Layout>
			<Switch>
				<Route path="/">
					<Home />
				</Route>
				<Route path="/database">
					<Database />
				</Route>
			</Switch>
		</Layout>
	);
}
