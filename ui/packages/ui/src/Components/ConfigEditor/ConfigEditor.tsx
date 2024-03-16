import AceEditor from "react-ace";
//THESE IMPORTS NEEDS TO BE AFTER IMPORTING AceEditor
import "ace-builds/src-noconflict/ext-language_tools";
import "../../util/mode-gcsim.js";
//manually import supported themes cause we can't get for loop to work here
import { FormGroup, HTMLSelect } from "@blueprintjs/core";
import "ace-builds/src-noconflict/theme-github";
import "ace-builds/src-noconflict/theme-kuroir";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/theme-solarized_dark";
import "ace-builds/src-noconflict/theme-solarized_light";
import "ace-builds/src-noconflict/theme-terminal";
import "ace-builds/src-noconflict/theme-textmate";
import "ace-builds/src-noconflict/theme-tomorrow";
import "ace-builds/src-noconflict/theme-tomorrow_night";
import "ace-builds/src-noconflict/theme-twilight";
import "ace-builds/src-noconflict/theme-xcode";
import React from "react";

export const themes = [
  "monokai",
  "github",
  "tomorrow",
  "tomorrow_night",
  "kuroir",
  "twilight",
  "xcode",
  "textmate",
  "solarized_dark",
  "solarized_light",
  "terminal",
];
export type Theme = (typeof themes)[number];

type Props = {
  cfg: string;
  onChange: (v: string) => void;
  hideThemeSelector?: boolean;
};

const LOCALSTORAGE_KEY = "gcsim-config-editor-theme";

export function ConfigEditor(props: Props) {
  const [theme, setTheme] = React.useState<Theme>(() => {
    return localStorage.getItem(LOCALSTORAGE_KEY) ?? "tomorrow_night";
  });
  React.useEffect(() => {
    localStorage.setItem(LOCALSTORAGE_KEY, theme);
  }, [theme]);
  const hideThemeSelector = props.hideThemeSelector ?? true;
  return (
    <div className="p-1 md:p-2">
      {hideThemeSelector ? null : (
        <div className="mb-1 w-full flex flex-row-reverse">
          <FormGroup label="Editor Theme" inline>
            <HTMLSelect onChange={(e) => setTheme(e.currentTarget.value)}>
              {themes.map((t) => (
                <option key={t} selected={t == theme}>
                  {t}
                </option>
              ))}
            </HTMLSelect>
          </FormGroup>
        </div>
      )}
      <AceEditor
        mode="gcsim"
        theme={theme}
        width="100%"
        onChange={props.onChange}
        value={props.cfg}
        name="config_editor"
        editorProps={{
          $blockScrolling: true,
        }}
        setOptions={{
          maxLines: Infinity,
          fontSize: 14,
          tabSize: 2,
          highlightActiveLine: false,
        }}
      />
    </div>
  );
}
