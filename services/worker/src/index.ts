/**
 * Welcome to Cloudflare Workers! This is your first worker.
 *
 * - Run `wrangler dev src/index.ts` in your terminal to start a development server
 * - Open a browser tab at http://localhost:8787/ to see your worker in action
 * - Run `wrangler publish src/index.ts --name my-worker` to publish your worker
 *
 * Learn more at https://developers.cloudflare.com/workers/
 */

import { PostgrestClient } from '@supabase/postgrest-js';
import { Router } from 'itty-router';
import { handleAssets } from './assets';
import { handleAuth } from './auth';
import { handleListDBChars, handleListDBSims } from './db';
import { handleListAllDBSims } from './db/listAll';
import { handleEnka } from './enka';
import { handlePreview } from './preview';
import { handleShare } from './share';
import { handleListUserSims } from './sims/listByUser';
import { handleView } from './view';

export const dbClient = new PostgrestClient(POSTGREST_ENDPOINT); //secrets?
const router = Router();

router.get('/api/auth', handleAuth);
router.post('/api/share', handleShare);

//cached
router.get('/api/view/:key', handleView);
router.get('/api/preview/:key', handlePreview);
router.get('/api/assets/*', handleAssets);

//enka
router.get('/api/enka/:key', handleEnka);

// db routes
router.get('/api/db', handleListDBChars);
router.get('/api/db/all', handleListAllDBSims);
router.get('/api/db/:key', handleListDBSims);

// user sims
router.get('/api/:key/sims', handleListUserSims);

addEventListener('fetch', (event) => {
  event.respondWith(router.handle(event.request, event));
});
