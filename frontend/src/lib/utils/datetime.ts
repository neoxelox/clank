import dayjs from "dayjs";
import advancedFormat from "dayjs/plugin/advancedFormat.js";
import localizedFormat from "dayjs/plugin/localizedFormat.js";
import relativeTime from "dayjs/plugin/relativeTime.js";

dayjs.extend(relativeTime);
dayjs.extend(advancedFormat);
dayjs.extend(localizedFormat);

export default dayjs;
