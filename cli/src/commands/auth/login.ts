import { prompt } from 'inquirer';
import { program } from 'commander';
import ora from 'ora';
import { writeAuthConfigFile } from '../../utils/config/auth';
import { requireMinCharacter } from '../../utils/validation/minCharacter';
import { requireValidEmail } from '../../utils/validation/email';
import { APISuccess } from '../../models/response';
import { IAuth } from '../..//models/auth';
import { login } from '../../client/service/auth/login';

const spinner = ora();

const loginCommand = () =>
  program
    .command('login')
    .description('authenticates users')
    .action(() => {
      prompt([
        {
          type: 'input',
          message: 'Enter a email',
          name: 'email',
          validate: (value: string) => requireValidEmail(value),
        },
        {
          type: 'password',
          message: 'Enter a password',
          name: 'password',
          mask: '•',
          validate: (value: string) => requireMinCharacter(value, 8, 'password'),
        },
      ]).then(
        async ({ password, email }): Promise<void> => {
          spinner.start('Waiting authetication...');
          try {
            const body: APISuccess<IAuth> = await login(email, password);

            const { data, message } = body;
            const { tokens, user } = data;

            await writeAuthConfigFile({
              isDefault: true,
              email: user.email,
              accessToken: tokens.accessToken,
              refreshToken: tokens.refreshToken,
            });

            spinner.succeed(message);
          } catch (error) {
            spinner.fail(error.message);
          }
        },
      );
    });

export default loginCommand;
