import { homedir } from 'os';
import { join } from 'path';
import { existsSync } from 'fs';
import { getEnvVariable } from './env';

/**
 * @description gets configuration file in local
 * @param {string} fileName
 */
const getConfigPath = (fileName: string): string => {
  const localFileDir: string = join(homedir(), getEnvVariable('LOCAL_DIR') || './marvin');
  const marvinAuthPath = join(localFileDir, fileName);
  const marvinAuthExist = existsSync(marvinAuthPath);

  if (!marvinAuthExist) {
    throw new Error('Config file is not found');
  }

  return marvinAuthPath;
};

export { getConfigPath };
