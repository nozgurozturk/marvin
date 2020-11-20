import { prompt } from 'inquirer';
import { program } from 'commander';
import ora from 'ora';
import { createRepository } from '../../client/service/repo/create';
import { createSubscriber } from '../../client/service/subscriber/create';
import { getDefaultUserEmail } from '../../utils/config/auth';

const spinner = ora('Creating');

const createRepositoryCommand = () =>
  program
    .command('create <url>')
    .description('Creates new repository')
    .action((url) =>
      prompt([
        {
          type: 'confirm',
          message: 'Do you want to add your email as a subscriber',
          name: 'addYourself',
        },
      ]).then(async ({ addYourself }) => {
        spinner.start('Creating');

        try {
          const body = await createRepository(url);
          const { data, message } = body;

          spinner.succeed(message);

          if (addYourself) {
            spinner.start('Adding...');

            const email = getDefaultUserEmail();
            const addBody = await createSubscriber(data.id!, email!);

            spinner.succeed(addBody.message);
          }
        } catch (error) {
          spinner.fail(error.message);
        }
      }),
    );

export default createRepositoryCommand;
