import React from "react";
import SearchButton from "../search_white_24dp.svg";
import { useLocation } from "wouter";
import { AppContext } from "../Store";

const SearchContext = React.createContext<string>("");

export const SearchProvider = SearchContext.Provider;

export default function SearchBar() {
  const { state, dispatch } = React.useContext(AppContext);
  const [_, setLocation] = useLocation();

  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    dispatch({ type: "set", str: e.currentTarget.value });
  };

  return (
    <div
      className="md:w-1/2 lg:w-2/3 sm:w-full bg-gray-500 rounded-lg p-3 m-1 h-auto"
      style={{ position: "relative" }}
    >
      <img
        src={SearchButton}
        alt="search"
        className="rounded-md hover:bg-gray-400"
        style={{
          objectFit: "contain",
          position: "absolute",
          right: 0,
          top: 0,
          height: "100%",
          maxHeight: "60px",
          padding: "0.5rem",
        }}
        onClick={() => {
          setLocation("/db/results");
        }}
      />
      <form
        onSubmit={(e) => {
          e.preventDefault();
          setLocation("/db/results");
        }}
      >
        <input
          type="text"
          className="p-1 w-full text-2xl md:text-xl bg-gray-500 outline-none"
          placeholder="Search for an action list..."
          value={state}
          onChange={handleChange}
        />
      </form>
    </div>
  );
}
