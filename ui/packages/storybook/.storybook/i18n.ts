import { initI18n } from "@gcsim/localization";
import Backend from "i18next-http-backend";

export default initI18n({ use: [Backend], init: { returnNull: false } });
