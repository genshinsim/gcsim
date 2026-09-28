import type { UserInfo } from "@gcsim/types";
import { describe, expect, it } from "vitest";
import { fakeStorage } from "./fakeStorage";
import { defaultUser, userStore } from "./user";

const loggedIn: UserInfo = {
	uid: "123",
	name: "someone",
	role: 1,
	permalinks: [],
	data: {
		settings: { showTips: false, showBuilder: true, showNameSearch: true },
	},
};

describe("userStore", () => {
	it("starts logged out with default settings", () => {
		expect(userStore(fakeStorage()).get()).toEqual(defaultUser);
	});

	it("merges a logged in user", () => {
		const user = userStore(fakeStorage());
		user.merge(loggedIn);
		expect(user.get()).toEqual(loggedIn);
		expect(defaultUser.uid).toBe("");
	});

	it("persists settings but not the login", () => {
		const storage = fakeStorage();
		userStore(storage).merge(loggedIn);
		expect(userStore(storage).get()).toEqual({
			...defaultUser,
			data: loggedIn.data,
		});
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
		expect(userStore(storage).get().data.settings).toEqual(settings);
	});

	it("resets to logged out", () => {
		const user = userStore(fakeStorage());
		user.merge(loggedIn);
		user.reset();
		expect(user.get()).toEqual(defaultUser);
	});
});
