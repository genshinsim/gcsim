import { describe, expect, it } from "vitest";
import { draftStore, mergeTeam } from "./draft";
import { fakeStorage } from "./fakeStorage";

const current = `bennett char lvl=90/90 cons=6 talent=9,9,9;
bennett add weapon="favoniussword" refine=5 lvl=90/90;
xiangling char lvl=90/90 cons=6 talent=9,9,9;
xiangling add weapon="thecatch" refine=5 lvl=90/90;
xiangling add set="emblemofseveredfate" count=4;

options iteration=100;
target lvl=100 resist=0.1;
bennett attack;
`;

const incoming = `raiden char lvl=90/90 cons=0 talent=9,9,9;
raiden add weapon="engulfinglightning" refine=1 lvl=90/90;
raiden add stats hp=4780 atk=311;

options iteration=1000;
active raiden;
raiden burst;
`;

describe("mergeTeam", () => {
	it("keeps the current draft's characters and the incoming config's rest", () => {
		expect(mergeTeam(current, incoming)).toBe(`
bennett char lvl=90/90 cons=6 talent=9,9,9;
bennett add weapon="favoniussword" refine=5 lvl=90/90;

xiangling char lvl=90/90 cons=6 talent=9,9,9;
xiangling add weapon="thecatch" refine=5 lvl=90/90;
xiangling add set="emblemofseveredfate" count=4;

options iteration=1000;
active raiden;
raiden burst;
`);
	});

	it("drops the incoming characters when the draft has none", () => {
		expect(mergeTeam("options iteration=100;\n", incoming)).toBe(`

options iteration=1000;
active raiden;
raiden burst;
`);
	});
});

describe("draftStore", () => {
	it("replaces the draft when keepTeam is false", () => {
		const draft = draftStore(fakeStorage());
		draft.set(current);
		draft.send(incoming, { keepTeam: false });
		expect(draft.get()).toBe(incoming);
	});

	it("merges the team when keepTeam is true", () => {
		const draft = draftStore(fakeStorage());
		draft.set(current);
		draft.send(incoming, { keepTeam: true });
		expect(draft.get()).toBe(mergeTeam(current, incoming));
	});

	it("persists the draft", () => {
		const storage = fakeStorage();
		draftStore(storage).set("options iteration=5;");
		expect(draftStore(storage).get()).toBe("options iteration=5;");
	});

	it("loads a draft saved by the redux store", () => {
		const storage = fakeStorage({
			"redux-app-data": JSON.stringify({ sampleOnLoad: true, cfg: current }),
		});
		expect(draftStore(storage).get()).toBe(current);
	});

	it("starts empty when nothing is stored", () => {
		expect(draftStore(fakeStorage()).get()).toBe("");
	});
});
