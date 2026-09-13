import type { IRequest } from "itty-router";
import type { Env } from "../bindings";

class ElementHandler {
	private key;
	private host;
	private prefix;

	constructor(host, key, prefix) {
		this.key = key;
		this.host = host;
		this.prefix = prefix;
	}

	element(element) {
		// An incoming element, such as `div`
		element.append(
			`<meta
    property="og:title"
    content="gcsim - simulation impact"
/>`,
			{ html: true },
		);
		element.append(
			`<meta
      property="og:site_name"
      content="gcsim"
  />`,
			{ html: true },
		);
		element.append(
			`<meta
        property="og:description"
        content=""
    />`,
			{ html: true },
		);
		element.append(
			`<meta property="og:image" content="${this.host}/api/preview/${this.prefix !== "" ? this.prefix + "/" : ""}${this.key}.png" />`,
			{ html: true },
		);
		element.append(`<meta property="og:image:width" content="540" />`, {
			html: true,
		});
		element.append(`<meta property="og:image:height" content="250" />`, {
			html: true,
		});
		element.append(`<meta property="og:image:type" content="image/png" />`, {
			html: true,
		});
		element.append(
			`<meta name="twitter:card" content="summary_large_image" />`,
			{
				html: true,
			},
		);
	}

	comments(comment) {
		// An incoming comment
	}

	text(text) {
		// An incoming piece of text
	}
}

export async function handleInjectHead(
	request: IRequest,
	env: Env,
): Promise<Response> {
	// Fetch the base HTML through the ASSETS binding (not a global fetch): with
	// `run_worker_first: true` a global fetch(request) re-enters the Worker and
	// yields a 404, so the head would be injected into the not-found page.
	const res = await env.ASSETS.fetch(request);
	const url = new URL(request.url);
	const segments = url.pathname.split("/");
	const key = segments.pop() || segments.pop();
	const host = url.protocol + "//" + url.host;
	console.log("received share request: " + key);

	return new HTMLRewriter()
		.on("head", new ElementHandler(host, key, ""))
		.transform(res);
}

export async function handleInjectHeadDB(
	request: IRequest,
	env: Env,
): Promise<Response> {
	// See handleInjectHead: fetch the base HTML via the ASSETS binding so the
	// SPA index.html is rewritten, not the Worker's own 404 page.
	const res = await env.ASSETS.fetch(request);
	const url = new URL(request.url);
	const segments = url.pathname.split("/");
	const key = segments.pop() || segments.pop();
	const host = url.protocol + "//" + url.host;
	console.log("received share request: " + key);

	return new HTMLRewriter()
		.on("head", new ElementHandler(host, key, "db"))
		.transform(res);
}
