// The formatter only rewrites whitespace between tokens, and sameTokens()
// checks that, so a formatter bug can make a config ugly but never change
// what it means.

type TokenKind =
	| "ident"
	| "number"
	| "string"
	| "field"
	| "comment"
	| "op"
	| "punct"
	| "error";

interface Token {
	kind: TokenKind;
	text: string;
	// newlines in the whitespace before this token
	newlines: number;
	// whether any whitespace precedes this token
	spaced: boolean;
}

// The tokenizer mirrors pkg/gcs/ast/lex.go rather than the Lezer grammar so
// token boundaries match the real lexer: identifiers may contain "-" and "%"
// (`i-1` is one identifier), ".foo" is a single field token, "<>" means !=.
const isSpace = (c: string | undefined) =>
	c === " " || c === "\t" || c === "\n" || c === "\r";
const isDigit = (c: string | undefined) =>
	c !== undefined && c >= "0" && c <= "9";
const isIdentChar = (c: string | undefined) =>
	c !== undefined && /[\p{L}\p{N}_%-]/u.test(c);

const TWO_CHAR_OPS = new Set(["==", "!=", "<=", ">=", "<>", "&&", "||"]);
const ONE_CHAR_OPS = new Set(["=", "<", ">", "!", "+", "-", "*", "/"]);
const PUNCT = new Set([";", ":", ",", "(", ")", "[", "]", "{", "}", "."]);
// characters the Go lexer accepts directly after an identifier or field
const TERMINATORS = new Set(".,|:)(+=><&!;[]{}/*");

function tokenize(src: string): Token[] {
	const tokens: Token[] = [];
	let i = 0;
	for (;;) {
		const gap = i;
		let newlines = 0;
		while (isSpace(src[i])) if (src[i++] === "\n") newlines++;
		if (i >= src.length) return tokens;

		const from = i;
		const c = src[i];
		const next = src[i + 1];
		let kind: TokenKind;
		if (c === "#" || (c === "/" && next === "/")) {
			while (i < src.length && src[i] !== "\n") i++;
			while (isSpace(src[i - 1])) i--;
			kind = "comment";
		} else if (c === '"') {
			i++;
			while (i < src.length && src[i] !== '"' && src[i] !== "\n")
				i += src[i] === "\\" && src[i + 1] !== "\n" ? 2 : 1;
			if (src[i] === '"') i++;
			kind = "string";
		} else if (isDigit(c) || (c === "." && isDigit(next))) {
			while (isDigit(src[i])) i++;
			if (src[i] === ".") {
				i++;
				while (isDigit(src[i])) i++;
			}
			kind = "number";
		} else if (c === "." && isIdentChar(next)) {
			i++;
			while (isIdentChar(src[i])) i++;
			kind = "field";
		} else if (isIdentChar(c) && c !== "-") {
			while (isIdentChar(src[i])) i++;
			kind = "ident";
		} else if (TWO_CHAR_OPS.has(c + next)) {
			i += 2;
			kind = "op";
		} else if (ONE_CHAR_OPS.has(c)) {
			i++;
			kind = "op";
		} else if (PUNCT.has(c)) {
			i++;
			kind = "punct";
		} else {
			i++;
			kind = "error";
		}
		tokens.push({
			kind,
			text: src.slice(from, i),
			newlines,
			spaced: from > gap,
		});
	}
}

export function sameTokens(a: string, b: string): boolean {
	const x = tokenize(a);
	const y = tokenize(b);
	return (
		x.length === y.length &&
		x.every((t, i) => t.kind === y[i].kind && t.text === y[i].text)
	);
}

// Statements made of `key=value` pairs that the Go parser reads into the
// ActionList rather than the AST. They print tight: `lvl=90/90 talent=9,9,9`.
const CONFIG_STARTERS = new Set([
	"options",
	"target",
	"energy",
	"hurt",
	"active",
]);
const CONFIG_VERBS = new Set(["char", "add"]);
const KEYWORDS = new Set([
	"if",
	"else",
	"while",
	"for",
	"switch",
	"case",
	"default",
	"return",
	"let",
	"fn",
	"break",
	"continue",
	"fallthrough",
]);
// tokens that stay on the `}` line instead of starting a new one
const JOINS_CLOSE_BRACE = new Set(["else", ";", ")", ",", "]"]);
// tokens that never start a continuation line
const NO_BREAK_BEFORE = new Set(["{", ";", ",", ")", "]", ":"]);

interface Frame {
	isSwitch: boolean;
	inCase: boolean;
	brackets: string[];
}

export interface FormatOptions {
	indent?: string;
}

export function formatGcsim(
	src: string,
	{ indent = "\t" }: FormatOptions = {},
): string {
	const tokens = tokenize(src);
	const lines: string[] = [];
	let line: string | null = null;
	let depth = 0;
	const frames: Frame[] = [{ isSwitch: false, inCase: false, brackets: [] }];
	let frame = frames[0];
	let mode: "config" | "script" = "script";
	let stmtStart = true;
	let pendingBreak = false;
	// no blank line directly after `{` or `case x:`
	let justOpened = false;
	let forHeader = false;
	let switchHeader = false;
	let caseHeader = false;
	let prev: Token | undefined;
	let prevUnary = false;

	const breakLine = () => {
		if (line !== null) lines.push(line);
		line = null;
	};
	const blankLine = () => {
		if (lines.length > 0 && lines[lines.length - 1] !== "") lines.push("");
	};
	const pad = (level: number) => indent.repeat(Math.max(0, level));

	const isUnary = (p: Token | undefined, t: Token) =>
		t.text === "!" ||
		((t.text === "-" || t.text === "+") &&
			(!p ||
				p.kind === "op" ||
				(p.kind === "punct" && !")]}".includes(p.text)) ||
				(p.kind === "ident" && KEYWORDS.has(p.text))));

	const configSpace = (p: Token, t: Token) => {
		// `+params=[...]` reads as a prefixed key
		if (t.text === "+") return p.text !== "=" && p.text !== ",";
		if (p.kind === "op" || t.kind === "op") return false;
		if (p.text === "," || t.text === "[" || t.text === "(") return false;
		return true;
	};

	const scriptSpace = (p: Token, t: Token) => {
		if (prevUnary) return false;
		const callee =
			(p.kind === "ident" && (!KEYWORDS.has(p.text) || p.text === "fn")) ||
			p.kind === "field" ||
			p.text === ")";
		if (t.text === "(") return !callee;
		// action params `skill[hold=1]`, map literals `= [a=1]`
		if (t.text === "[") return !callee && p.text !== "]";
		if (t.kind === "field" && p.kind === "field") return false;
		const inSquare = frame.brackets[frame.brackets.length - 1] === "[";
		if (inSquare && (t.text === "=" || p.text === "=")) return false;
		// repeat counts `attack:3`
		if (p.text === ":") return t.kind !== "number";
		return true;
	};

	const needsSpace = (p: Token, t: Token) => {
		const joined = tokenize(p.text + t.text);
		if (
			joined.length !== 2 ||
			joined[0].text !== p.text ||
			joined[1].text !== t.text
		)
			return true;
		if (
			(p.kind === "ident" || p.kind === "field") &&
			!TERMINATORS.has(t.text[0])
		)
			return true;
		if (p.kind === "error" || t.kind === "error") return t.spaced;
		if (p.text === "." || t.text === ".") return t.spaced;
		if (",;)]:".includes(t.text) || p.text === "(" || p.text === "[")
			return false;
		return mode === "config" ? configSpace(p, t) : scriptSpace(p, t);
	};

	for (const [k, t] of tokens.entries()) {
		const trailingComment =
			t.kind === "comment" && t.newlines === 0 && line !== null;
		if (
			pendingBreak &&
			!trailingComment &&
			!(prev?.text === "}" && JOINS_CLOSE_BRACE.has(t.text))
		)
			breakLine();
		pendingBreak = false;

		if (t.kind === "comment") {
			if (line !== null && trailingComment) {
				line = `${line} ${t.text}`;
			} else {
				breakLine();
				if (t.newlines > 1 && !justOpened) blankLine();
				line = pad(depth + (stmtStart ? 0 : 1)) + t.text;
			}
			// a comment runs to end of line, so the next token must start a new one
			breakLine();
			justOpened = false;
			continue;
		}

		// leaving a block or a case body dedents before the token is placed
		if (t.text === "}") {
			breakLine();
			if (frames.length > 1) {
				if (frame.inCase) depth--;
				depth--;
				frames.pop();
				frame = frames[frames.length - 1];
			}
		} else if (
			stmtStart &&
			(t.text === "case" || t.text === "default") &&
			frame.isSwitch &&
			frame.inCase
		) {
			depth--;
			frame.inCase = false;
		}

		const unary = isUnary(prev, t);
		let head: string;
		if (line === null) {
			if (stmtStart && t.newlines > 1 && !justOpened && t.text !== "}")
				blankLine();
			head = pad(depth + (stmtStart || t.text === "}" ? 0 : 1));
		} else if (t.newlines > 0 && !stmtStart && !NO_BREAK_BEFORE.has(t.text)) {
			// keep the author's line breaks inside a long statement
			breakLine();
			head = pad(depth + 1);
		} else {
			head = prev && needsSpace(prev, t) ? `${line} ` : line;
		}
		line = head + t.text;
		justOpened = false;

		const atTop = frame.brackets.length === 0;
		if (t.text === "{") {
			frames.push({ isSwitch: switchHeader, inCase: false, brackets: [] });
			frame = frames[frames.length - 1];
			depth++;
			forHeader = switchHeader = false;
			stmtStart = pendingBreak = justOpened = true;
		} else if (t.text === "}") {
			stmtStart = pendingBreak = true;
		} else if (t.text === ";" && atTop && !forHeader) {
			stmtStart = pendingBreak = true;
			caseHeader = false;
		} else if (t.text === ":" && atTop && caseHeader) {
			caseHeader = false;
			frame.inCase = true;
			depth++;
			stmtStart = pendingBreak = justOpened = true;
		} else {
			if (t.text === "(" || t.text === "[") frame.brackets.push(t.text);
			if (t.text === ")" || t.text === "]") frame.brackets.pop();
			if (stmtStart) {
				const verb = tokens[k + 1]?.text ?? "";
				mode =
					CONFIG_STARTERS.has(t.text) ||
					(t.kind === "ident" && CONFIG_VERBS.has(verb))
						? "config"
						: "script";
				caseHeader =
					frame.isSwitch && (t.text === "case" || t.text === "default");
			}
			if (t.text === "for") forHeader = true;
			if (t.text === "switch") switchHeader = true;
			stmtStart = false;
		}
		prev = t;
		prevUnary = unary;
	}
	breakLine();
	return lines.length > 0 ? `${lines.join("\n")}\n` : "";
}
