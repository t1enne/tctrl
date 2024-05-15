import { confirm, input, password, select } from "npm:@inquirer/prompts";
import { existsSync } from "https://deno.land/std@0.224.0/fs/mod.ts";
import { parseArgs } from "https://deno.land/std@0.224.0/cli/parse_args.ts";
import { API_BASE, AUTH_HEADERS, req } from "./utils.ts";
import dayjs from "npm:dayjs";
import utc from "npm:dayjs/plugin/utc.js";
import pc from "npm:picocolors";
import { BaseDto, ILoginPrompt, UserHours } from "./types.ts";

// load utc plugin
dayjs.extend(utc);

const CACHE_FILE = "tcontrol.json";
const STATE = {
  clients: [] as BaseDto[],
  projects: {} as Record<string, BaseDto[]>,
  releases: {} as Record<string, BaseDto[]>,
};

const getWorkedHrs = async (inputDate: dayjs.Dayjs, userId: string) => {
  const from = dayjs(inputDate).utc().startOf("day").toISOString();
  const to = dayjs(inputDate).utc().endOf("day").toISOString();
  const body = {
    relations: [
      "release",
      "release.project",
      "release.project.customer",
      "hoursTag",
    ],
    where: {
      userId,
      date: {
        _fn: 17,
        args: [from, to],
      },
    },
  };
  const r = await req<UserHours[]>("/userHours/fb", body);
  return r.data.reduce(
    (acc, { release, hours, hoursTag }) =>
      // deno-fmt-ignore
      acc + `${pc.bold(pc.yellow(release.project.name))}: ${ pc.italic(release.name) }: ${pc.green(`󱑂  ${hours} (${hoursTag.name.slice(0, 3)})`)}\n`,
    "",
  );
};

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

const pickTag = async () => {
  const { data } = await req<Array<BaseDto>>("/hoursTags/fb", {});
  return await select({
    message: "Tag:",
    choices: data.map((c) => ({ value: c, name: c.name })),
  });
};
const pickRelease = async (projectId: string) => {
  const body = { order: { name: "ASC" }, where: { projectId } };
  const releases = STATE.releases[projectId]
    ? STATE.releases[projectId]
    : await req<Array<BaseDto>>("/releases/fb", body).then((r) => r.data);
  STATE.releases[projectId] = releases;
  return await select({
    message: "Release:",
    choices: releases.map((c) => ({ value: c, name: c.name })),
  });
};

const pickProject = async (customerId: string) => {
  const body = { order: { name: "ASC" }, where: { customerId } };
  const projects = STATE.projects[customerId]
    ? STATE.projects[customerId]
    : await req<Array<BaseDto>>("/projects/fb", body).then((r) => r.data);
  STATE.projects[customerId] = projects;
  return await select({
    message: "Project:",
    choices: projects.map((c) => ({ value: c, name: c.name })),
  });
};

const pickClient = async () => {
  const body = { order: { name: "ASC" } };
  const clients = STATE.clients.length
    ? STATE.clients
    : await req<Array<BaseDto>>("/customers/fb", body).then((r) => r.data);
  STATE.clients = clients;
  return await select({
    message: "Client:",
    choices: clients.map((c) => ({ value: c, name: c.name })),
  });
};

const getInputHrs = async (hrsLeft: number) =>
  await input({
    message: "Hours:",
    validate: (v) => {
      const float = parseFloat(v);
      if (isNaN(float)) return "Must be a number";
      if (0 > float) return "Must be positive";
      if (float > hrsLeft) return `Can't be more than ${hrsLeft} hours`;
      return Number.isFinite(float);
    },
  });

const getInputDate = (arg: string | number | undefined) => {
  const today = dayjs().utc();
  if (typeof arg === "undefined") return today;
  if (typeof arg === "number") {
    return arg > 0 ? today.subtract(arg, "day") : today.add(arg, "day");
  }
  // arg == string
  return dayjs.utc(arg, "DD-MM-YY");
};

async function main() {
  Deno.addSignalListener("SIGINT", () => {
    console.log(pc.red("Bye!"));
    Deno.exit();
  });
  console.log(pc.cyan("RAINTONIC"));
  const args = parseArgs(Deno.args) as { d?: number; l?: string; a?: boolean };
  const inputDate = getInputDate(args.d || args.l);

  const CACHE = existsSync(CACHE_FILE)
    ? Deno.readTextFileSync(CACHE_FILE)
    : undefined;

  const user = CACHE
    ? JSON.parse(CACHE) as ReturnType<typeof getUser>
    : getUser(
      await authUser(await getInput()),
    );

  if (!CACHE && user) {
    console.log("writing cache file...");
    Deno.writeTextFileSync(CACHE_FILE, JSON.stringify(user, null, 2));
  }
  AUTH_HEADERS.Authorization = `Bearer ${user.token}`;
  // list worked hours
  if (args.a) {
    for (let i = 1; i <= dayjs.utc().date(); i++) {
      const current = dayjs().set("date", i);
      const workedHrs = await getWorkedHrs(current, user.id);
      if (workedHrs.length) {
        console.log(pc.bgWhite(pc.black(current.format("DD-MM-YY"))));
        console.log(workedHrs);
      }
    }
    return;
  }

  // write worked hours
  console.log(await getWorkedHrs(inputDate, user.id));

  let total = 0;
  while (total < 8) {
    const client = await pickClient();
    const project = await pickProject(client._id);
    const release = await pickRelease(project._id);
    const tag = await pickTag();
    const toUploadBody = {
      notes: await input({ message: "Notes:" }),
      hours: await getInputHrs(8 - total),
      date: dayjs(Date.now()).utc().toISOString(),
      releaseId: release._id,
      hoursTagId: tag._id,
      userId: user.id,
    };

    const { data } = await req<{ notes: string; hours: string }>(
      "/userHours",
      toUploadBody,
    );
    total += parseFloat(toUploadBody.hours);
    console.log(`Saved ${data.notes} X ${data.hours}`);
    if (total < 8) {
      const shouldContinue = await confirm({ message: "Continue?" });
      if (!shouldContinue) {
        break;
      }
    }
  }
  Deno.exit(0);
}

await main();
