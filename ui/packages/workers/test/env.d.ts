import type { Env as WorkerEnv } from "../src/bindings";

declare global {
	namespace Cloudflare {
		interface GlobalProps {
			mainModule: typeof import("../src/index");
		}
		interface Env extends WorkerEnv {}
	}
}
