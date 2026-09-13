import { Router } from "itty-router";
import { handleAssets } from "./assets";
import type { Env } from "./bindings";
import { handleEnka } from "./enka";
import { handleInjectHead, handleInjectHeadDB, handlePreview } from "./preview";
import { proxyRequest } from "./proxy";
import { handleLegacy, handleShare, handleView } from "./share";
import { handleWasm } from "./wasm";

const router = Router();

//passthrough
router.get("/api/login", proxyRequest);
router.post("/api/user/save", proxyRequest);
router.get("/api/share/random", proxyRequest);
router.get("/api/db/compute/work", proxyRequest);
router.post("/api/db/compute/work", proxyRequest);
router.post("/api/db/submit", proxyRequest);
router.get("/api/db", proxyRequest);
router.get("/api/db/*", proxyRequest);
// viewer files
router.post("/api/share", handleShare);
router.get("/api/share/:key", handleView);
router.get("/api/share/db/:key", handleView);
router.get("/api/legacy-share/:key", handleLegacy); //TODO: this endpoint should be deleted once we convert over to new
router.get("/api/preview/:key", handlePreview);
router.get("/api/preview/db/:key", handlePreview);

//enka
router.get("/api/enka/:key", handleEnka);

// rewrite doc head
router.get("/sh/:key", handleInjectHead);
router.get("/db/:key", handleInjectHeadDB);

router.get("/api/assets/*", handleAssets);
router.get("/api/wasm/*", handleWasm);

// SPA fallthrough: anything not matched above is served from the static-assets
// binding. `not_found_handling: "single-page-application"` makes deep links
// resolve to index.html. Kept explicit so the routing intent lives in code too.
router.all("*", (request: Request, env: Env) => env.ASSETS.fetch(request));

export default {
	async fetch(
		request: Request,
		env: Env,
		ctx: ExecutionContext,
	): Promise<Response> {
		return router.fetch(request, env, ctx);
	},
};
