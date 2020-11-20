import { prompt } from 'inquirer';
import { program } from 'commander';
import ora from 'ora';
import { findAllRepositories } from '../../client/service/repo/findAll';
import { createSubscriber } from '../../client/service/subscriber/create';
import { requireValidEmail } from '../../utils/validation/email';

const spinner = ora();

const createSubscriberCommand = () =>
  program
    .command('sub-add')
    .description('Add subscriber to repository')
    .action(async () => {
      try {
        spinner.start('Searching repositories');
        const reposBody = await findAllRepositories();
        const { data, message } = reposBody;
        spinner.succeed(message);
        prompt([
          {
            type: 'list',
            message: 'Please select repository that you want to add subscriber',
            name: 'repoName',
            choices: data.map((repo) => repo.name),
          },
          {
            type: 'input',
            message: 'Enter a email',
            name: 'email',
            validate: (value: string) => requireValidEmail(value),
          },
        ]).then(async ({ repoName, email }) => {
          try {
            spinner.start('Adding subscriber to repository');

            const repo = data.find((repo) => repo.name === repoName);
            if (!repo) {
              throw new Error('Undefined repository');
            }

            const subBody = await createSubscriber(repo.id!, email);

            spinner.succeed(subBody.message);
          } catch (error) {
            spinner.fail(error.message);
          }
        });
      } catch (error) {
        spinner.fail(error.message);
      }
    });

export default createSubscriberCommand;
