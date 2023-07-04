import { Router } from "itty-router";
import { handleAssets } from "./assets";
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
// viewer files
router.post("/api/share", handleShare);
router.get("/api/share/:key", handleView);
router.get("/api/share/db/:key", handleView);
router.get("/api/legacy-share/:key", handleLegacy); //TODO: this endpoint should be deleted once we convert over to new
router.get("/api/preview/:key", handlePreview);
router.get("/api/preview/db/:key", handlePreview);

// rewrite doc head
router.get("/sh/:key", handleInjectHead);
router.get("/db/:key", handleInjectHeadDB);

router.get("/api/assets/*", handleAssets);
router.get("/api/wasm/*", handleWasm);

addEventListener("fetch", (event) => {
  event.respondWith(router.handle(event.request, event));
});
