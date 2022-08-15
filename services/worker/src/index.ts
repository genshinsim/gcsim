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
import { handleAuth } from './auth';
import { handleListDBChars, handleListDBSims } from './db';
import { handleOptions } from './options';
import { handlePreview } from './preview';
import { handleShare } from './share';
import { handleView } from './view';

export const dbClient = new PostgrestClient(POSTGREST_ENDPOINT); //secrets?
const router = Router();

router.options('*', handleOptions);

router.get('/api/auth', handleAuth);

router.post('/api/share', handleShare);

router.get('/api/view/:key', handleView);

router.get('/api/preview/:key', handlePreview);

router.get('/api/db', handleListDBChars);

router.get('/api/db/:key', handleListDBSims);

router.get('/api/avatars', async () => {
  const { data, error } = await dbClient.from('avatars').select();

  if (error) {
    console.log('error getting avatars: ', error);
    return new Response(JSON.stringify({ error: error }), {
      status: 500,
      statusText: 'Error getting avatars',
      headers: {
        'content-type': 'application/json',
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Max-Age': '86400',
      },
    });
  }

  return new Response(JSON.stringify({ users: data }), {
    headers: {
      'content-type': 'application/json',
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Max-Age': '86400',
    },
  });
});

addEventListener('fetch', (event) => {
  event.respondWith(router.handle(event.request));
});
