import { prompt } from 'inquirer';
import { program } from 'commander';
import ora from 'ora';
import { writeAuthConfigFile } from '../../utils/config/auth';
import { requireMinCharacter } from '../../utils/validation/minCharacter';
import { requireValidEmail } from '../../utils/validation/email';
import { APISuccess } from '../../models/response';
import { IAuth } from '../../models/auth';
import { signup } from '../../client/service/auth/signup';

const spinner = ora();

const signupCommand = () =>
  program
    .command('signup')
    .description('Creates user')
    .action(() =>
      prompt([
        {
          type: 'input',
          message: 'Enter a user name',
          name: 'name',
          validate: (value: string) => requireMinCharacter(value, 6, 'email'),
        },
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
        {
          type: 'password',
          message: 'Confirm password',
          name: 'confirmPassword',
          mask: '•',
          validate: (value: string, { password }) => {
            if (value === password) {
              return true;
            }
            return 'Passwords must match';
          },
        },
      ]).then(
        async ({ password, email, name }): Promise<void> => {
          spinner.start('Waiting authentication...');
          try {
            const body: APISuccess<IAuth> = await signup(name, email, password);

            const { data } = body;
            const { tokens, user } = data;

            await writeAuthConfigFile({
              isDefault: true,
              email: user.email,
              accessToken: tokens.accessToken,
              refreshToken: tokens.refreshToken,
            });

            spinner.succeed(body.message);
          } catch (error) {
            spinner.fail(error.message);
          }
        },
      ),
    );

export default signupCommand;
