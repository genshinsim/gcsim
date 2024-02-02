import { MenuItem } from "@blueprintjs/core";
import { MultiSelect2 } from "@blueprintjs/select";
import { useContext } from "react";
import { useTranslation } from "react-i18next";
import { FilterDispatchContext, FilterContext, ItemFilterState, charNames } from "./FilterComponents/Filter.utils";

export function CharacterQuickSelect() {
    //dispatch
    const dispatch = useContext(FilterDispatchContext);
    const filter = useContext(FilterContext);
    const { t } = useTranslation();
  
    const includedChars = Object.entries(filter.charFilter)
      .map(([charName, charState]) => {
        if (charState.state === ItemFilterState.include) {
          return charName;
        }
      })
      .filter((charName) => charName) as string[];
  
    const translateCharName = (charName: string) =>
      t("game:character_names." + charName);
    return (
      <div className="grow max-w-xl">
        <MultiSelect2
          placeholder={t<string>("db.type_to_search")}
          items={charNames}
          itemRenderer={(charName, itemProps) => {
            return (
              <MenuItem
                key={charName}
                text={translateCharName(charName)}
                icon={
                  <img
                    src={`/api/assets/avatar/${charName}.png`}
                    className="w-6 h-6"
                  />
                }
                onClick={() => {
                  dispatch({
                    type: "includeChar",
                    char: charName,
                  });
                }}
                active={itemProps.modifiers.active}
              />
            );
          }}
          tagRenderer={(charName) => (
            <div className="flex flex-row gap-1" key={charName}>
              <img
                className="w-4 h-4"
                src={`/api/assets/avatar/${charName}.png`}
              />
              {translateCharName(charName)}
            </div>
          )}
          onItemSelect={(charName) => {
            if (!charName) {
              return;
            }
            dispatch({
              type: "handleChar",
              char: charName,
            });
          }}
          itemListPredicate={(query, items) => {
            return items.filter((item) => {
              const normalizedItem = item.toLowerCase();
              const normalizedLocalizedItem = translateCharName(item).toLowerCase();
              const normalizedQuery = query.toLocaleLowerCase();
              return normalizedItem.includes(normalizedQuery) || normalizedLocalizedItem.includes(normalizedQuery);
            });
          }}
          selectedItems={includedChars}
          onClear={() => {
            dispatch({
              type: "clearFilter",
            });
          }}
          onRemove={(charName) => {
            dispatch({
              type: "includeChar",
              char: charName,
            });
          }}
          resetOnSelect
          resetOnQuery
          openOnKeyDown
          tagInputProps={{
            tagProps: {
              minimal: true,
            },
            onRemove: (value) => {
              if (!value) return;
              if (!value["key"]) return;
              dispatch({
                type: "removeChar",
                char: value["key"],
              });
            },
          }}
        ></MultiSelect2>
      </div>
    );
  }