import { confirm, input, select } from "npm:@inquirer/prompts";
import { STATE } from "../src/common.ts";
import { BaseDto, UploadBody } from "../src/types.ts";
import { req } from "../src/utils.ts";
import dayjs, { Dayjs } from "npm:dayjs";

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

const getInputHrs = async () =>
  await input({
    message: "Hours:",
    validate: (v) => {
      const float = parseFloat(v);
      if (isNaN(float)) return "Must be a number";
      if (0 > float) return "Must be positive";
      if (float > 8) return `Can't be more than ${8} hours`;
      return Number.isFinite(float);
    },
  });

export async function uploadWorkedHrs(
  body: Omit<UploadBody, "date">,
  date: dayjs.Dayjs,
) {
  const { data } = await req<{ notes: string; hours: string }>(
    "/userHours",
    { ...body, date: date.toISOString() },
  );
  return data.hours;
}

export async function fillWorkedHrs(user: { id: string; token: string }) {
  const client = await pickClient();
  const project = await pickProject(client._id);
  const release = await pickRelease(project._id);
  const tag = await pickTag();
  const toUploadBody = {
    notes: await input({ message: "Notes:" }),
    hours: await getInputHrs(),
    releaseId: release._id,
    hoursTagId: tag._id,
    userId: user.id,
  };
  return toUploadBody;
  // total += parseFloat(toUploadBody.hours);
  // console.log(`Saved ${data.notes} X ${data.hours}`);
  // if (total < 8) {
  //   const shouldContinue = await confirm({ message: "Continue?" });
  //   if (!shouldContinue) {
  //     break;
  //   }
  // }
}
