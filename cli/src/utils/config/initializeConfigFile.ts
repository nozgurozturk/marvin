import { homedir } from 'os';
import { join } from 'path';
import { existsSync, mkdir, writeFile } from 'fs';
import { getEnvVariable } from './env';

export const setDefaultConfigDir = async () => {
  const localFileDir: string = join(homedir(), getEnvVariable('LOCAL_DIR'));
  const exist = existsSync(localFileDir);

  if (exist) {
    return;
  }

  mkdir(localFileDir, (err) => {
    if (err) {
      throw err;
    }
  });
};

/**
 * @description gets configuration file in local
 * @param {string} fileName
 */
export const setDefaultConfigFile = async (fileName: string) => {
  const localFileDir: string = join(homedir(), getEnvVariable('LOCAL_DIR'));
  const path = join(localFileDir, fileName);
  const exist = existsSync(path);

  if (exist) {
    return;
  }
  writeFile(path, JSON.stringify({}), { encoding: 'utf-8' }, (err) => {
    if (err) {
      throw err;
    }
  });
};
