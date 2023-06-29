import { useContext } from "react";
import { FaChevronLeft, FaChevronRight } from "react-icons/fa";
import {
  FilterContext,
  FilterDispatchContext,
} from "./FilterComponents/Filter.utils";

//
export function PaginationButtons() {
  const dispatch = useContext(FilterDispatchContext);

  const filter = useContext(FilterContext);
  // write a component that uses the filter context to render the pagination buttons
  return (
    <div className="flex flex-row justify-center">
      <div className="flex flex-row gap-2">
        <button
          className="bp4-button bp4-large"
          onClick={() => {
            dispatch({ type: "decrementPage" });
            scrollToTop();
          }}
        >
          <FaChevronLeft />
        </button>
        {
          //set page number input
          <input
            className="bp4-input bp4-large w-12 text-center"
            type="number"
            min={1}
            value={filter.pageNumber}
            onChange={(e) => {
              dispatch({
                type: "setPage",
                pageNumber: parseInt(e.target.value),
              });
            }}
          />
        }
        <button
          className="bp4-button bp4-large"
          onClick={() => {
            dispatch({ type: "incrementPage" });
            scrollToTop();
          }}
        >
          <FaChevronRight />
        </button>
      </div>
    </div>
  );
}

function scrollToTop() {
  const isBrowser = () => typeof window !== "undefined";
  if (!isBrowser()) return;
  window.scrollTo({ top: 0, behavior: "smooth" });
}
