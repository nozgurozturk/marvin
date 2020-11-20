#!/usr/bin/env node
import { program } from 'commander';
import { version } from '../package.json';
import { auth } from './commands/auth';
import { repository } from './commands/repo';
import { subscriber } from './commands/subscriber';
import { setDefaultConfigDir, setDefaultConfigFile } from './utils/config/initializeConfigFile';

const main = async () => {
  await setDefaultConfigDir();
  await setDefaultConfigFile('auth.json');
};

main().then(() => {
  program.version(version);

  //Commands
  auth();
  repository();
  subscriber();

  program.parse(process.argv);
});
