import dayjs from "npm:dayjs";
import pc from "npm:picocolors";
import { BaseDto, UserHours } from "./types.ts";
import { req } from "./utils.ts";

export const STATE = {
  clients: [] as BaseDto[],
  projects: {} as Record<string, BaseDto[]>,
  releases: {} as Record<string, BaseDto[]>,
  user: {} as { id: string; token: string },
};

export const getWorkedHrs = async (inputDate: dayjs.Dayjs, userId: string) => {
  console.log(inputDate);
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
    (acc, { release, hours, hoursTag }, i) =>
      // deno-fmt-ignore
      acc + `${pc.bold(pc.yellow(release.project.name))}: ${ pc.italic(release.name) }: ${pc.green(`󱑂  ${hours} (${hoursTag.name.slice(0, 3)})`)}${i === r.data.length -1 ? "" : "\n" }`,
    "",
  );
};
