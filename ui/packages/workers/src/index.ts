import { Router } from "itty-router";
import { handleAssets } from "./assets";
import { proxyRequest } from "./proxy";
import { handleLegacy, handleShare, handleView } from "./share";

const router = Router();

// viewer files
router.post("/api/share", handleShare);
router.get("/api/share/:key", handleView);
router.get("/api/legacy-share/:key", handleLegacy); //TODO: this endpoint should be deleted once we convert over to new

router.get("/api/assets/*", handleAssets);
router.get("/api/login", proxyRequest);
router.post("/api/user/save", proxyRequest);

addEventListener("fetch", (event) => {
  event.respondWith(router.handle(event.request, event));
});
