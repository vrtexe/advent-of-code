import { dirname } from '@std/path';

const DenoTaskScript = 'deno run --watch --allow-read';

async function main() {
  const [task, yearArg] = Deno.args;
  const year = yearArg ?? getCurrentYear();

  const paddedTask = task.padStart(2, '0');

  const taskFileName = `src/${year}/task-${paddedTask}.ts`;
  const dataFileName = `${year}/task${task}.txt`;
  const testDataFileName = `${year}/task${task}-test.txt`;

  if (await pathExists(taskFileName)) {
    console.error(`Script with name: ${taskFileName}, Already Exists!`);
    return;
  }

  createDataFiles(dataFileName, testDataFileName);
  createScriptFile(testDataFileName, taskFileName);
  createScriptCommand(paddedTask, year, taskFileName);
}

function getCurrentYear() {
  return new Date().getFullYear().toString();
}

async function createScriptFile(dataFileName: string, fileName: string) {
  const scriptTemplate = await Deno.readTextFile('assets/template.ts');
  const startingScript = scriptTemplate.replace('{{dataFile}}', dataFileName);

  const directory = dirname(fileName);
  if (!(await pathExists(directory))) {
    await Deno.mkdir(directory, { recursive: true });
  }

  await Deno.writeTextFile(fileName, startingScript);
}

async function createScriptCommand(
  day: string,
  year: string,
  fileName: string,
) {
  const denoConfigData = JSON.parse(
    await Deno.readTextFile('deno.json'),
  ) as DenoConfig;

  const denoTaskName = `task:${year}:${day}`;

  if (denoConfigData.tasks[denoTaskName]) {
    console.error(`Deno task with name: ${denoTaskName}, Already exists!`);
    return;
  }

  denoConfigData.tasks[denoTaskName] = `${DenoTaskScript} ${fileName}`;

  await Deno.writeTextFile('deno.json', JSON.stringify(denoConfigData));

  const formatCommand = new Deno.Command(Deno.execPath(), {
    args: ['fmt', 'deno.json'],
  });

  await formatCommand.output();

  console.log(`deno run ${denoTaskName}`);
}

async function createDataFiles(...dataFiles: string[]) {
  for (const dataFile of dataFiles) {
    const directory = dirname(dataFile);
    if (!(await pathExists(directory))) {
      await Deno.mkdir(`assets/${directory}`, { recursive: true });
    }
    Deno.create(`assets/${dataFile}`);
  }
}

async function pathExists(fileName: string) {
  try {
    await Deno.lstat(fileName);
    return true;
  } catch (e) {
    if (!(e instanceof Deno.errors.NotFound)) {
      throw e;
    }

    return false;
  }
}

if (import.meta.main) {
  main();
}

type DenoConfig = {
  tasks: Record<string, string>;
};
