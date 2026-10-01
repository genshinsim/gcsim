import { LanguageSupport, LRLanguage } from "@codemirror/language";
import { styleTags, tags as t } from "@lezer/highlight";
import { gcsimCompletionSource } from "./autocomplete";
import { parser } from "./parser.gen";

export const gcsimLanguage = LRLanguage.define({
	name: "gcsim",
	parser: parser.configure({
		props: [
			styleTags({
				LineComment: t.lineComment,
				Number: t.number,
				StringLiteral: t.string,
				Keyword: t.keyword,
				Bool: t.bool,
				CharacterName: t.className,
				ActionName: t.function(t.variableName),
				StatName: t.attributeName,
				ElementName: t.atom,
				Identifier: t.variableName,
				Operator: t.operator,
				Punctuation: t.punctuation,
			}),
		],
	}),
	languageData: {
		commentTokens: { line: "#" },
		closeBrackets: { brackets: ["(", "[", "{", '"'] },
	},
});

export function gcsim(): LanguageSupport {
	return new LanguageSupport(gcsimLanguage, [
		gcsimLanguage.data.of({ autocomplete: gcsimCompletionSource }),
	]);
}
