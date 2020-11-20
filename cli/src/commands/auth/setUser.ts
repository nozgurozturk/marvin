import { program } from 'commander';
import { prompt } from 'inquirer';
import { readAuthConfigFile, setDefaultUserToConfig } from '../../utils/config/auth';

const setUser = () => {
  const config = readAuthConfigFile();
  const emails = Object.keys(config);

  const defaultUser = emails.find((email) => config[email].isDefault);

  prompt({
    type: 'list',
    message: `Please select default user's email. Current: ${defaultUser}`,
    name: 'email',
    choices: emails,
  }).then(({ email }) => {
    setDefaultUserToConfig(email);
  });
};

const setDefaultUserCommand = () =>
  program
    .command('use')
    .description('Sets default user')
    .action(() => setUser());

export default setDefaultUserCommand;
