import { Card, Toaster } from "@gcsim/primitives";
import "@gcsim/components/src/index.css";
import type { Entry } from "@gcsim/types/src/generated/index.db";
import axios from "axios";
import React, { type JSX } from "react";
import { Route, Switch } from "wouter";
import { DecisionButtons, DuplicateRow, ReviewHeadline } from "./shared";

function App({ id }: { id: string }) {
	const [main, setMain] = React.useState<Entry | null>(null);
	const [data, setData] = React.useState<Entry[]>([]);

	React.useEffect(() => {
		axios.get(`/api/db/id/${id}`).then((res) => {
			if (res.data) {
				setMain(res.data);
			}
		});
	}, [id]);

	React.useEffect(() => {
		if (main === null) {
			return;
		}
		const includedChars = main.summary?.char_names;
		if (includedChars === null || includedChars === undefined) {
			return;
		}
		const q = {
			query: {},
			limit: 100,
			sort: {
				created_data: -1,
			},
		};
		if (includedChars.length > 0) {
			const and: unknown[] = [];
			const trav: { [key in string]: boolean } = {};
			includedChars.forEach((char) => {
				if (char.includes("aether") || char.includes("lumine")) {
					const ele = char.replace(/(aether|lumine)(.+)/, "$2");
					trav[ele] = true;
					return;
				}
				and.push({
					"summary.char_names": char,
				});
			});
			Object.keys(trav).forEach((ele) => {
				and.push({
					$or: [
						{ "summary.char_names": `aether${ele}` },
						{ "summary.char_names": `lumine${ele}` },
					],
				});
			});
			if (and.length > 0) {
				q.query["$and"] = and;
			}
		}
		axios
			.get(`/api/db?q=${encodeURIComponent(JSON.stringify(q))}`)
			.then((res) => {
				if (res.data && res.data.data && res.data.data.length > 0) {
					setData(res.data.data);
				}
			});
	}, [main]);

	if (main === null) {
		return (
			<div className="p-g-page text-g-ink-dim">Loading, please wait...</div>
		);
	}

	const duplicates = data.filter((e) => e._id !== id);

	return (
		<div className="mx-auto flex max-w-4xl flex-col gap-g-section p-g-page">
			<section className="flex flex-col gap-g-base">
				<div className="flex items-center justify-between">
					<h2 className="font-g-display text-g-h3 font-semibold text-g-ink">
						Under review
					</h2>
					<span className="font-g-mono text-g-xs text-g-ink-mute">id {id}</span>
				</div>
				<Card className="flex flex-col gap-g-base border-g-accent/40 p-g-card ring-1 ring-g-accent/20">
					<ReviewHeadline entry={main} />
					<div className="flex flex-wrap gap-g-base border-t border-g-line-soft pt-g-base [&>*]:flex-1 sm:justify-end sm:[&>*]:flex-none">
						<DecisionButtons entry={main} />
					</div>
				</Card>
			</section>

			<section className="flex flex-col gap-g-base">
				<div className="flex items-center gap-g-base">
					<h2 className="font-g-display text-g-h3 font-semibold text-g-ink">
						Existing sims with the same team
					</h2>
					<span className="rounded-g-pill bg-g-surface-2 px-2 py-0.5 font-g-mono text-g-xs text-g-ink-dim">
						{duplicates.length}
					</span>
				</div>
				{duplicates.length === 0 ? (
					<div className="rounded-g-md border border-dashed border-g-line py-10 text-center text-g-ink-mute">
						No existing sims share this team.
					</div>
				) : (
					<div className="flex flex-col gap-g-base-lg">
						{duplicates.map((d) => (
							<DuplicateRow key={d._id} entry={d} main={main} />
						))}
					</div>
				)}
			</section>
			<Toaster />
		</div>
	);
}

const Routes = (): JSX.Element => {
	return (
		<>
			<Switch>
				<Route path="/">
					<div>nothing here</div>
				</Route>
				<Route path="/id/:id">{({ id }) => <App id={id} />}</Route>
			</Switch>
		</>
	);
};

export default Routes;
