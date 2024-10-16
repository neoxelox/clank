import Amazon from "$lib/components/icons/Amazon.svelte";
import AppStore from "$lib/components/icons/AppStore.svelte";
import IAgora from "$lib/components/icons/IAgora.svelte";
import PlayStore from "$lib/components/icons/PlayStore.svelte";
import Trustpilot from "$lib/components/icons/Trustpilot.svelte";
import type { SvelteComponent } from "svelte";

export enum FeedbackSource {
  TRUSTPILOT = "TRUSTPILOT",
  PLAY_STORE = "PLAY_STORE",
  APP_STORE = "APP_STORE",
  AMAZON = "AMAZON",
  IAGORA = "IAGORA",
}

export type FeedbackSourceDetail = {
  title: string;
  icon: typeof SvelteComponent;
};

export const FeedbackSourceDetails: Record<FeedbackSource | string, FeedbackSourceDetail> = {
  [FeedbackSource.TRUSTPILOT]: {
    title: "Trustpilot",
    icon: Trustpilot,
  },
  [FeedbackSource.PLAY_STORE]: {
    title: "Play Store",
    icon: PlayStore,
  },
  [FeedbackSource.APP_STORE]: {
    title: "App Store",
    icon: AppStore,
  },
  [FeedbackSource.AMAZON]: {
    title: "Amazon",
    icon: Amazon,
  },
  [FeedbackSource.IAGORA]: {
    title: "iAgora",
    icon: IAgora,
  },
};

export enum FeedbackLanguage {
  AFRIKAANS = "AFRIKAANS",
  ALBANIAN = "ALBANIAN",
  ARABIC = "ARABIC",
  ARMENIAN = "ARMENIAN",
  AZERBAIJANI = "AZERBAIJANI",
  BASQUE = "BASQUE",
  BELARUSIAN = "BELARUSIAN",
  BENGALI = "BENGALI",
  BOKMAL = "BOKMAL",
  BOSNIAN = "BOSNIAN",
  BULGARIAN = "BULGARIAN",
  CATALAN = "CATALAN",
  CHINESE = "CHINESE",
  CROATIAN = "CROATIAN",
  CZECH = "CZECH",
  DANISH = "DANISH",
  DUTCH = "DUTCH",
  ENGLISH = "ENGLISH",
  ESPERANTO = "ESPERANTO",
  ESTONIAN = "ESTONIAN",
  FINNISH = "FINNISH",
  FRENCH = "FRENCH",
  GANDA = "GANDA",
  GEORGIAN = "GEORGIAN",
  GERMAN = "GERMAN",
  GREEK = "GREEK",
  GUJARATI = "GUJARATI",
  HEBREW = "HEBREW",
  HINDI = "HINDI",
  HUNGARIAN = "HUNGARIAN",
  ICELANDIC = "ICELANDIC",
  INDONESIAN = "INDONESIAN",
  IRISH = "IRISH",
  ITALIAN = "ITALIAN",
  JAPANESE = "JAPANESE",
  KAZAKH = "KAZAKH",
  KOREAN = "KOREAN",
  LATIN = "LATIN",
  LATVIAN = "LATVIAN",
  LITHUANIAN = "LITHUANIAN",
  MACEDONIAN = "MACEDONIAN",
  MALAY = "MALAY",
  MAORI = "MAORI",
  MARATHI = "MARATHI",
  MONGOLIAN = "MONGOLIAN",
  NYNORSK = "NYNORSK",
  PERSIAN = "PERSIAN",
  POLISH = "POLISH",
  PORTUGUESE = "PORTUGUESE",
  PUNJABI = "PUNJABI",
  ROMANIAN = "ROMANIAN",
  RUSSIAN = "RUSSIAN",
  SERBIAN = "SERBIAN",
  SHONA = "SHONA",
  SLOVAK = "SLOVAK",
  SLOVENE = "SLOVENE",
  SOMALI = "SOMALI",
  SOTHO = "SOTHO",
  SPANISH = "SPANISH",
  SWAHILI = "SWAHILI",
  SWEDISH = "SWEDISH",
  TAGALOG = "TAGALOG",
  TAMIL = "TAMIL",
  TELUGU = "TELUGU",
  THAI = "THAI",
  TSONGA = "TSONGA",
  TSWANA = "TSWANA",
  TURKISH = "TURKISH",
  UKRAINIAN = "UKRAINIAN",
  URDU = "URDU",
  VIETNAMESE = "VIETNAMESE",
  WELSH = "WELSH",
  XHOSA = "XHOSA",
  YORUBA = "YORUBA",
  ZULU = "ZULU",
}

export const NO_LANGUAGE: string = "UNKNOWN";

export const NO_RELEASE: string = "UNKNOWN";

export type FeedbackCustomer = {
  email?: string;
  name: string;
  picture: string;
  location?: string;
  verified?: boolean;
  reviews?: number;
  link?: string;
};

export type FeedbackMetadata = {
  rating?: number;
  media?: string[];
  verified?: boolean;
  votes?: number;
  link?: string;
};

export type Feedback = {
  id: string;
  productID: string;
  source: string;
  customer: FeedbackCustomer;
  content: string;
  language: string;
  translation: string;
  release: string;
  metadata: FeedbackMetadata;
  postedAt: Date;
};
