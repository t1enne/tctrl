import dayjs from "npm:dayjs";
import { WEEKDAYS } from "../src/utils.ts";
import { getWorkedHrs, STATE } from "../src/common.ts";
import pc from "npm:picocolors";

export default async function list(_?: number) {
  const apiCalls = [];
  // const isPastMonth = month === undefined;
  for (let i = 1; i <= dayjs.utc().date(); i++) {
    const current = dayjs().set("date", i);
    const weekDay = current.day();
    if (!(0 < weekDay && weekDay < 6)) {
      // saturday and sunday
      continue;
    }
    apiCalls.push(getWorkedHrs(current, STATE.user.id));
  }
  const responses = await Promise.all(apiCalls);
  for (let i = 0; i < responses.length; i++) {
    const current = dayjs().set("date", i + 1);
    const weekDay = current.day();
    const workedHrs = responses[i];
    // deno-fmt-ignore
    console.log( pc.bgWhite( pc.black(`󰃭 ${current.format("DD-MM-YY")} (${WEEKDAYS[weekDay]})`)) + `\n${workedHrs}`);
  }
}
