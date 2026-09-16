import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { MultiSelect } from "./MultiSelect";

type Item = { key: string; label: string };

const items: Item[] = [
	{ key: "amber", label: "Amber" },
	{ key: "xiangling", label: "Xiangling" },
	{ key: "xingqiu", label: "Xingqiu" },
];

function renderSelect(
	overrides: Partial<Parameters<typeof MultiSelect<Item>>[0]> = {},
) {
	const onChange = vi.fn();
	render(
		<MultiSelect<Item>
			items={items}
			itemKey={(item) => item.key}
			itemPredicate={(item, query) =>
				item.label.toLowerCase().includes(query.toLowerCase())
			}
			itemRenderer={(item, { selected, onSelect }) => (
				<button type="button" data-selected={selected} onClick={onSelect}>
					{item.label}
				</button>
			)}
			tagRenderer={(item, { onRemove }) => (
				<span key={item.key}>
					{item.label}
					<button
						type="button"
						aria-label={`remove-${item.key}`}
						onClick={onRemove}
					>
						x
					</button>
				</span>
			)}
			value={[]}
			onChange={onChange}
			{...overrides}
		/>,
	);
	return { onChange };
}

describe("MultiSelect", () => {
	it("renders a tag for each selected value", () => {
		renderSelect({ value: [items[0], items[2]] });
		expect(screen.getByText("Amber")).toBeInTheDocument();
		expect(screen.getByText("Xingqiu")).toBeInTheDocument();
		expect(screen.queryByText("Xiangling")).not.toBeInTheDocument();
	});

	it("opens the item list and adds an item on select", async () => {
		const { onChange } = renderSelect();
		await userEvent.click(screen.getByRole("button", { name: /select/i }));
		await userEvent.click(screen.getByRole("button", { name: "Amber" }));
		expect(onChange).toHaveBeenCalledWith([items[0]]);
	});

	it("removes an already-selected item when chosen again from the list", async () => {
		const { onChange } = renderSelect({ value: [items[0]] });
		await userEvent.click(screen.getByRole("button", { name: /select/i }));
		await userEvent.click(screen.getByRole("button", { name: "Amber" }));
		expect(onChange).toHaveBeenCalledWith([]);
	});

	it("removes an item via its tag's remove control", async () => {
		const { onChange } = renderSelect({ value: [items[0], items[1]] });
		await userEvent.click(screen.getByLabelText("remove-amber"));
		expect(onChange).toHaveBeenCalledWith([items[1]]);
	});

	it("marks items already in value as selected in the list", async () => {
		renderSelect({ value: [items[1]] });
		await userEvent.click(screen.getByRole("button", { name: /select/i }));
		expect(screen.getByRole("button", { name: "Xiangling" })).toHaveAttribute(
			"data-selected",
			"true",
		);
		expect(screen.getByRole("button", { name: "Amber" })).toHaveAttribute(
			"data-selected",
			"false",
		);
	});
});
