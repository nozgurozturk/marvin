import { writeFile, readFileSync } from 'fs';
import { getConfigPath } from './file-path';

type AuthSet = {
  isDefault: boolean;
  accessToken: string;
  refreshToken: string;
  email: string;
};

type AuthGet = {
  isDefault: boolean;
  accessToken: string;
  refreshToken: string;
};

/**
 * @description read auth configuration file
 */

export const readAuthConfigFile = () => {
  const AUTH_PATH = getConfigPath('auth.json');
  const authFile = readFileSync(AUTH_PATH, { encoding: 'utf-8' });
  return JSON.parse(authFile);
};

/**
 * @description set auth credential to config file
 * @param {AuthSet} auth
 */
export const writeAuthConfigFile = async (auth: AuthSet) => {
  const AUTH_PATH = getConfigPath('auth.json');
  let config = readAuthConfigFile();

  if (Object.keys(config).length > 0) {
    Object.keys(config).map((key) => {
      config[key].isDefault = false;
    });
  }

  config[auth.email] = {
    isDefault: true,
    accessToken: auth.accessToken,
    refreshToken: auth.refreshToken,
  };

  writeFile(AUTH_PATH, JSON.stringify(config), { encoding: 'utf-8' }, (err) => {
    if (err) {
      console.error(err);
      throw err;
    }
  });
};

/**
 * @description get auth credentails of user from file
 * @param {string} email
 */

export const getAuthConfig = (): AuthGet | undefined => {
  const config = readAuthConfigFile();
  const emails = Object.keys(config);

  const defaultUser = emails.find((email) => config[email].isDefault);

  if (!defaultUser) {
    console.error('Please set default user use command: marvin use <email>');
    return;
  }

  return config[defaultUser];
};

export const getDefaultUserEmail = () => {
  const config = readAuthConfigFile();
  const emails = Object.keys(config);

  const defaultUser = emails.find((email) => config[email].isDefault);

  return defaultUser;
};

/**
 * @description set default user and save the file
 * @param {string} email
 */
export const setDefaultUserToConfig = (email: string) => {
  const AUTH_PATH = getConfigPath('auth.json');
  let config = readAuthConfigFile();

  if (Object.keys(config).length > 0) {
    Object.keys(config).map((key) => {
      config[key].isDefault = false;
    });
  }

  config[email] = {
    isDefault: true,
    ...config[email],
  };

  writeFile(AUTH_PATH, JSON.stringify(config), { encoding: 'utf-8' }, (err) => {
    if (err) {
      console.error(err);
      throw err;
    }
  });
};
