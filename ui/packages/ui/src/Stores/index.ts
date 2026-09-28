import { draftStore } from "./draft";
import { prefsStore } from "./prefs";
import { userStore } from "./user";

export { useStore } from "./externalStore";

export const draft = draftStore(localStorage);
export const prefs = prefsStore(localStorage);
export const user = userStore(localStorage);
