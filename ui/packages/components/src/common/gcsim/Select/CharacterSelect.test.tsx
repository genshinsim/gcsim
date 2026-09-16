import { initI18n } from "@gcsim/localization";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeAll, describe, expect, it, vi } from "vitest";
import { CharacterSelect } from "./CharacterSelect";

beforeAll(() => {
	initI18n();
});

describe("CharacterSelect", () => {
	it("finds a character by its key or its translated name", async () => {
		const onSelect = vi.fn();
		render(<CharacterSelect isOpen onClose={() => {}} onSelect={onSelect} />);

		await userEvent.type(screen.getByRole("combobox"), "xiangling");
		expect(await screen.findByText("Xiangling")).toBeInTheDocument();

		await userEvent.click(screen.getByText("Xiangling"));
		expect(onSelect).toHaveBeenCalledWith("xiangling");
	});

	it("does not match unrelated characters", async () => {
		render(<CharacterSelect isOpen onClose={() => {}} onSelect={() => {}} />);
		await userEvent.type(screen.getByRole("combobox"), "xiangling");
		expect(screen.queryByText("Amber")).not.toBeInTheDocument();
	});
});
