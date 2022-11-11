import { Router } from "itty-router";
import { handleShare, handleView } from "./share";

const router = Router();

// viewer files
router.post("/api/share", handleShare);
router.get("/api/share/:key", handleView);

addEventListener("fetch", (event) => {
  event.respondWith(router.handle(event.request, event));
});
