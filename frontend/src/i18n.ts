import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import English from "./messages/en.json"
import Turkish from "./messages/tr.json"
import German from "./messages/de.json"

const resources = {
  en: { translations: English },
  de: { translations: German },
  tr: { translations: Turkish }
};

i18n
  .use(initReactI18next)
  .init({
    resources,
    lng: "en",
    ns: ['translations'],
    defaultNS: 'translations',
    debug: true,
    interpolation: {
      escapeValue: false
    }
  });

export default i18n;
