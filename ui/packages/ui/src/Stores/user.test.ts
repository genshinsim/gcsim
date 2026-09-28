import type { UserInfo } from "@gcsim/types";
import { describe, expect, it } from "vitest";
import { fakeStorage } from "./fakeStorage";
import { defaultUser, loadUser, mergeUser, saveUserSettings } from "./user";

const loggedIn: UserInfo = {
	uid: "123",
	name: "someone",
	role: 1,
	permalinks: [],
	data: {
		settings: { showTips: false, showBuilder: true, showNameSearch: true },
	},
};

describe("mergeUser", () => {
	it("merges a logged in user without touching the previous one", () => {
		expect(mergeUser(defaultUser, loggedIn)).toEqual(loggedIn);
		expect(defaultUser.uid).toBe("");
	});
});

describe("loadUser", () => {
	it("starts logged out with default settings", () => {
		expect(loadUser(fakeStorage())).toEqual(defaultUser);
	});

	it("restores saved settings but not the login", () => {
		const storage = fakeStorage();
		saveUserSettings(storage, loggedIn);
		expect(loadUser(storage)).toEqual({ ...defaultUser, data: loggedIn.data });
	});

	it("loads settings saved by the redux store", () => {
		const settings = {
			showTips: false,
			showBuilder: false,
			showNameSearch: true,
		};
		const storage = fakeStorage({
			"redux-user-local-settings": JSON.stringify(settings),
		});
		expect(loadUser(storage).data.settings).toEqual(settings);
	});
});
