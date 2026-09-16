import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { OmniSelect } from "./OmniSelect";

type Item = { key: string; label: string };

const items: Item[] = [
	{ key: "amber", label: "Amber" },
	{ key: "xiangling", label: "Xiangling" },
	{ key: "xingqiu", label: "Xingqiu" },
];

function renderSelect(
	overrides: Partial<Parameters<typeof OmniSelect<Item>>[0]> = {},
) {
	const onSelect = vi.fn();
	const onClose = vi.fn();
	render(
		<OmniSelect<Item>
			isOpen
			onClose={onClose}
			items={items}
			itemKey={(item) => item.key}
			itemPredicate={(item, query) =>
				item.label.toLowerCase().includes(query.toLowerCase())
			}
			itemRenderer={(item, { selected, onSelect: select }) => (
				<button type="button" data-selected={selected} onClick={select}>
					{item.label}
				</button>
			)}
			onSelect={onSelect}
			{...overrides}
		/>,
	);
	return { onSelect, onClose };
}

describe("OmniSelect", () => {
	it("renders all items when the query is empty", () => {
		renderSelect();
		expect(screen.getByText("Amber")).toBeInTheDocument();
		expect(screen.getByText("Xiangling")).toBeInTheDocument();
		expect(screen.getByText("Xingqiu")).toBeInTheDocument();
	});

	it("filters items using the supplied predicate", async () => {
		renderSelect();
		await userEvent.type(screen.getByRole("combobox"), "xing");
		expect(screen.queryByText("Amber")).not.toBeInTheDocument();
		expect(screen.getByText("Xingqiu")).toBeInTheDocument();
		expect(screen.queryByText("Xiangling")).not.toBeInTheDocument();
	});

	it("calls onSelect with the chosen item without closing on its own", async () => {
		const { onSelect, onClose } = renderSelect();
		await userEvent.click(screen.getByText("Amber"));
		expect(onSelect).toHaveBeenCalledWith(items[0]);
		expect(onClose).not.toHaveBeenCalled();
	});

	it("marks the current value as selected", () => {
		renderSelect({ value: items[1] });
		expect(screen.getByText("Xiangling")).toHaveAttribute(
			"data-selected",
			"true",
		);
		expect(screen.getByText("Amber")).toHaveAttribute("data-selected", "false");
	});

	it("resets the query each time it is reopened", () => {
		const { rerender } = render(
			<OmniSelect<Item>
				isOpen
				onClose={() => {}}
				items={items}
				itemKey={(item) => item.key}
				itemPredicate={(item, query) =>
					item.label.toLowerCase().includes(query.toLowerCase())
				}
				itemRenderer={(item) => <span>{item.label}</span>}
				onSelect={() => {}}
			/>,
		);
		rerender(
			<OmniSelect<Item>
				isOpen={false}
				onClose={() => {}}
				items={items}
				itemKey={(item) => item.key}
				itemPredicate={(item, query) =>
					item.label.toLowerCase().includes(query.toLowerCase())
				}
				itemRenderer={(item) => <span>{item.label}</span>}
				onSelect={() => {}}
			/>,
		);
		rerender(
			<OmniSelect<Item>
				isOpen
				onClose={() => {}}
				items={items}
				itemKey={(item) => item.key}
				itemPredicate={(item, query) =>
					item.label.toLowerCase().includes(query.toLowerCase())
				}
				itemRenderer={(item) => <span>{item.label}</span>}
				onSelect={() => {}}
			/>,
		);
		expect(screen.getByText("Amber")).toBeInTheDocument();
	});
});
