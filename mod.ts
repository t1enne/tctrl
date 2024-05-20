import dayjs from "npm:dayjs";
import utc from "npm:dayjs/plugin/utc.js";
import pc from "npm:picocolors";
import { input, password } from "npm:@inquirer/prompts";
import { existsSync } from "https://deno.land/std@0.224.0/fs/mod.ts";
import { parseArgs } from "https://deno.land/std@0.224.0/cli/parse_args.ts";
import { API_BASE, AUTH_HEADERS } from "./src/utils.ts";
import { CliOptions, ILoginPrompt } from "./src/types.ts";
import { getWorkedHrs, STATE } from "./src/common.ts";
import list from "./commands/list.ts";
import { uploadDate } from "./commands/uploadDate.ts";

// load utc plugin
dayjs.extend(utc);

const CACHE_FILE = "tcontrol.json";

const authUser = async ({ username, password }: ILoginPrompt) =>
  await fetch(`${API_BASE}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email: `${username}@raintonic.com`, password }),
  }).then((j) => j.json());

const getInput = async () => ({
  username: await input({ message: "Username: " }),
  password: await password({ message: "Password: ", mask: "*" }),
});

const getUser = (
  { data }: { data: { user: { _id: string }; token: string } },
) => ({ id: data.user._id, token: data.token });

const getInputDate = (arg: string | number | undefined) => {
  console.log(arg);
  const today = dayjs().utc();
  if (typeof arg === "undefined") return today;
  if (typeof arg === "number") {
    return arg > 0 ? today.subtract(arg, "day") : today.add(arg, "day");
  }
  // arg == string
  return dayjs(arg, "DD/MM/YY");
};

async function main() {
  Deno.addSignalListener("SIGINT", () => {
    console.log(pc.red("\nBye!"));
    Deno.exit(0);
  });
  const args = parseArgs<Partial<CliOptions>>(Deno.args);

  // SHOW HELP
  if (args.h || args.help || Object.keys(args).length === 1) {
    console.dir(
      `RAINTONIC CLI TPCA
Usage: tpca [options]
-d, --date 			Specify date to fill worked hours in the format DD/MM/YY
-r, --relative 	Same as above but allows specifying days as a relative integer (1,-1)
-m, --month 		List worked hours for the month specified as integer. Default is current month
-h, --help 			Show help`,
    );
    Deno.exit(0);
  }

  // INVALID ARGS
  if (Object.keys(args).length > 2) {
    console.log("You can only specify one argument. Use --help for usage info");
    Deno.exit(0);
  }

  const CACHE = existsSync(CACHE_FILE)
    ? Deno.readTextFileSync(CACHE_FILE)
    : undefined;

  STATE.user = CACHE
    ? JSON.parse(CACHE) as ReturnType<typeof getUser>
    : getUser(
      await authUser(await getInput()),
    );
  const { user } = STATE;
  if (!CACHE && user) {
    console.log("writing cache file...");
    Deno.writeTextFileSync(CACHE_FILE, JSON.stringify(user, null, 2));
  }
  AUTH_HEADERS.Authorization = `Bearer ${user.token}`;

  // list worked hours
  if (args.m || args.month) {
    await list();
    Deno.exit(0);
  }

  // RELATIVE DATE -1,1
  if (args.r || args.relative) {
    const inputDate = getInputDate(args.r || args.relative);
    console.log(getWorkedHrs(inputDate, user.id));
    await uploadDate(user, inputDate);
  }

  if (args.d || args.date) {
    const inputDate = getInputDate(args.d || args.date);
    console.log(getWorkedHrs(inputDate, user.id));
    await uploadDate(user, inputDate);
  }
  Deno.exit(0);
}

await main();
