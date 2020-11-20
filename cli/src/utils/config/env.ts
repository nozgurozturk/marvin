import { config, DotenvParseOutput } from 'dotenv';
// Read .env file
const ENV = config();

/**
 * @description gets enviromental variable in .env file
 * @param {string} key
 */
export const getEnvVariable = (key: string): string => {
  if (ENV.error) {
    throw ENV.error;
  }

  const variables: DotenvParseOutput | undefined = ENV.parsed;

  if (!variables || !Object.keys(variables).length) {
    throw new Error('Enviroment file does not contain any variable');
  }

  if (!String(variables[key])) {
    throw new Error(`${key} is not found in enviroment file`);
  }
  return variables[key];
};
