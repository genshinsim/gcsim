import { model } from "@gcsim/types";
import axios from "axios";

export const fetchDataFromDB = async (
  urlParams: any,
  setData: React.Dispatch<
    React.SetStateAction<model.IDBEntry[] | null | undefined>
  >
) => {
  try {
    const response = await axios.get(
      "https://simimpact.app/api/db" + urlParams
    );

    setData(response.data);
  } catch (error) {
    console.log(error);
  }
};
