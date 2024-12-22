import dayjs from "npm:dayjs";
import utc from "npm:dayjs/plugin/utc.js";
import customParseFormat from "npm:dayjs/plugin/customParseFormat.js";
import pc from "npm:picocolors";
import { parseArgs } from "node:util";
import { Commands } from "./src/commands.ts";

// load utc plugin
dayjs.extend(utc);
dayjs.extend(customParseFormat);

await main();

async function main() {
  Deno.addSignalListener("SIGINT", () => {
    console.log(pc.red("\nBye!"));
    Deno.exit(0);
  });

  const args = parseArgs({
    args: Deno.args,
    options: {
      help: { type: "boolean", short: "h" },
      relative: { type: "string", short: "r" },
      month: {
        type: "string",
        short: "m",
        default: `${dayjs.utc().month() + 1}`,
      },
      date: { type: "string", short: "d" },
    },
    allowPositionals: true,
  });

  // SHOW HELP
  switch (true) {
    case args.values.help:
      Commands.help();
      break;
    case !!args.values.date:
      await Commands.date([args.values.date].concat(args.positionals));
      break;
    case !!args.values.month:
      await Commands.month(+args.values.month);
      break;
    case !!args.values.relative:
      await Commands.relative(+args.values.relative);
      break;
    default:
      Deno.exit(1);
  }

  Deno.exit(0);
}
