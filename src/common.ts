import dayjs from "dayjs";
import { BaseDto, UserHours } from "./types.ts";
import { req, WEEKDAYS } from "./utils.ts";
import pc from "npm:picocolors";

export const STATE = {
  clients: [] as BaseDto[],
  projects: {} as Record<string, BaseDto[]>,
  releases: {} as Record<string, BaseDto[]>,
  user: {} as { id: string; token: string },
};

const pad = (dir: "l" | "r") => (l: number, str: string, ch = " ") =>
  dir === "l"
    ? ch.repeat(Math.max(0, l - str.length)) + str
    : str + ch.repeat(Math.max(0, l - str.length));

export const leftPad = pad("l");
export const rightPad = pad("r");

export const isWeekday = (d: number) => 0 < d && d < 6;

const getDayStart = (d: dayjs.Dayjs) => dayjs(d).startOf("day").toISOString();
const getDayEnd = (d: dayjs.Dayjs) => dayjs(d).endOf("day").toISOString();

export const getDateInfo = (
  date: dayjs.Dayjs,
  uploadedDay: UserHours[] = [],
) => {
  const today = dayjs.utc().isSame(date) ? "* " : "";
  const dayOfWeek = `${today}${WEEKDAYS[date.day()]}`;

  const fmt = (a: string) =>
    uploadedDay.length || !isWeekday(date.day()) ? a : pc.bold(pc.yellow(a));

  const str = `${date.format("YY-MM-DD")}     ${dayOfWeek}`;
  return fmt(str) + getHoursInfo(uploadedDay);
};

const getHoursInfo = (data: UserHours[]) => {
  return data.reduce(
    (acc, { release, hours, hoursTag }) => {
      const prjName = `[${(release.project.name.slice(0, 15))}] `;
      const relName = `(${release.name})`;
      const hrsTag = `{${hoursTag.name.slice(0, 3).toLowerCase()}}`;
      const hrs = `${hours} `;
      return acc +
        pc.gray(
          pc.dim(`\n${" ".repeat(13)}${hrs} ${prjName} ${relName} ${hrsTag}`),
        );
    },
    "",
  );
};

export const getWorkedHrs = async (
  userId: string,
  from: dayjs.Dayjs,
  to?: dayjs.Dayjs,
) => {
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
        args: [getDayStart(from), getDayEnd(to || from)],
      },
    },
  };
  const r = await req<UserHours[]>("/userHours/fb", body);
  return r.data;
};
