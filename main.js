const { spawn, Thread, Worker } = require("threads");
const { sleep } = require("sleep");

async function main() {
  console.log("subscribing ...");

  const join = await spawn(new Worker("./proc.js"));
  join();

  sleep(20);

  await Thread.terminate(join);

  main();
}

main().catch(console.error);
