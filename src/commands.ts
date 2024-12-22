import dayjs from "dayjs";
import pc from "picocolors";
import { getDateInfo, getWorkedHrs, STATE } from "../src/common.ts";
import { input, password } from "npm:@inquirer/prompts";
import { existsSync } from "https://deno.land/std@0.224.0/fs/mod.ts";
import { API_BASE, AUTH_HEADERS } from "../src/utils.ts";
import { ILoginPrompt } from "../src/types.ts";
import { fillWorkedHrs, uploadWorkedHrs } from "../commands/uploadDate.ts";

type UserRes = { data: { user: { _id: string }; token: string } };
const CACHE_FILE = "tcontrol.json";

const loadUser = async () => {
  const CACHE = existsSync(CACHE_FILE)
    ? Deno.readTextFileSync(CACHE_FILE)
    : undefined;

  STATE.user = CACHE
    ? JSON.parse(CACHE) as ReturnType<typeof getUser>
    : await authUser(await getInput());

  const { user } = STATE;
  if (!CACHE && user) {
    console.log("writing cache file...");
    Deno.writeTextFileSync(CACHE_FILE, JSON.stringify(user, null, 2));
  }
  AUTH_HEADERS.Authorization = `Bearer ${user.token}`;
  return user;
};

const authUser = async ({ username, password }: ILoginPrompt) =>
  await fetch(`${API_BASE}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email: `${username}@raintonic.com`, password }),
  }).then((j) => j.json()).then(getUser);

const getInput = async () => ({
  username: await input({ message: "Username: " }),
  password: await password({ message: "Password: ", mask: "*" }),
});

const getInputDate = (arg?: string) => {
  return typeof arg === "string" ? dayjs.utc(arg, "YY-MM-DD") : dayjs.utc();
};
const getUser = ({ data }: UserRes) => ({
  id: data.user._id,
  token: data.token,
});

export const Commands = {
  help() {
    console.dir(
      `RAINTONIC CLI TPCA
Usage: tpca [options]
-d, --date 			Specify date to fill worked hours in the format YY-MM-DD
-r, --relative 	Same as above but accepts an int day to go back to
-m, --month 		List worked hours for the month specified as integer. Default is current month
-h, --help 			Show help`,
    );
    Deno.exit(0);
  },
  async month(m: number) {
    const user = await loadUser();
    const now = dayjs.utc();
    const selectedMonth = now.utc().set("month", (m as number) - 1);
    const endOfMonth = selectedMonth.endOf("month").date();
    const days = await getWorkedHrs(
      user.id,
      selectedMonth.set("date", 1),
      selectedMonth.endOf("month").endOf("day"),
    );
    for (let i = 1; i <= endOfMonth; i++) {
      const date = selectedMonth.utc().set("date", i);
      if (date.isAfter(now)) {
        continue;
      }
      const userHours = days.filter((d) =>
        d.date.slice(0, 10) === date.startOf("day").toISOString().slice(0, 10)
      );
      console.log(getDateInfo(date, userHours));
    }
  },
  async date(dates: string[]) {
    const user = await loadUser();
    const workedHrs = await fillWorkedHrs(user);
    for (const date of dates) {
      const inputDate = getInputDate(date);
      const worked = await getWorkedHrs(user.id, inputDate) || [];
      console.log(pc.dim(getDateInfo(inputDate, worked)));
      await uploadWorkedHrs(workedHrs, inputDate);
    }
    Deno.exit(0);
  },
  async relative(relative: number) {
    const today = dayjs.utc();
    const date = today.set("date", today.date() - relative).format("YY-MM-DD");
    await Commands.date([date]);
  },
};
