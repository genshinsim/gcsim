import axios from "axios";
import { model } from "../../protos_gen/protos";

export const fetchDataFromDB = async (
  urlParams: any,
  setData: React.Dispatch<
    React.SetStateAction<model.IDBEntry[] | null | undefined>
  >
) => {
  const response = await axios.get("https://simimpact.app/api/db" + urlParams);
  setData(response.data);
};
