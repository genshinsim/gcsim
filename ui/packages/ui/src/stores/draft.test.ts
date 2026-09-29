import { describe, expect, it } from "vitest";
import { loadDraft, mergeTeam, saveDraft, sendDraft } from "./draft";
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

describe("sendDraft", () => {
	it("replaces the draft when keepTeam is false", () => {
		expect(sendDraft(current, incoming, { keepTeam: false })).toBe(incoming);
	});

	it("merges the team when keepTeam is true", () => {
		expect(sendDraft(current, incoming, { keepTeam: true })).toBe(
			mergeTeam(current, incoming),
		);
	});
});

describe("loadDraft", () => {
	it("round-trips a saved draft", () => {
		const storage = fakeStorage();
		saveDraft(storage, "options iteration=5;");
		expect(loadDraft(storage)).toBe("options iteration=5;");
	});

	it("loads a draft saved by the redux store", () => {
		const storage = fakeStorage({
			"redux-app-data": JSON.stringify({ sampleOnLoad: true, cfg: current }),
		});
		expect(loadDraft(storage)).toBe(current);
	});

	it("starts empty when nothing is stored", () => {
		expect(loadDraft(fakeStorage())).toBe("");
	});
});
