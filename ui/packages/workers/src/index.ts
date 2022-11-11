import { Router } from "itty-router";
import { handleAssets } from "./assets";
import { handleShare, handleView } from "./share";

const router = Router();

// viewer files
router.post("/api/share", handleShare);
router.get("/api/share/:key", handleView);

router.get("/api/assets/*", handleAssets);

addEventListener("fetch", (event) => {
  event.respondWith(router.handle(event.request, event));
});
