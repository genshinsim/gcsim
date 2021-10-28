import React from "react";
import { Link } from "wouter";
import SearchBar from "./SearchBar";

export default function ActionDB() {
  return (
    <div className="flex-grow flex flex-col items-center justify-center">
      <SearchBar />
      <div className="p-2">
        <Link href="/browse">
          <a href="/browse">Browse All</a>
        </Link>
      </div>
    </div>
  );
}
