import { prompt } from 'inquirer';
import { program } from 'commander';
import ora from 'ora';
import { findAllRepositories } from '../../client/service/repo/findAll';
import { findAllSubscriber } from '../../client/service/subscriber/findAll';
import chalk from 'chalk';

const spinner = ora();

const findAllSubscribersCommand = () =>
  program
    .command('sub-list')
    .description('List all subscribers belongs to repository, reds are unverfied')
    .action(async () => {
      try {
        spinner.start('Searching repositories');

        const reposBody = await findAllRepositories();
        const { data, message } = reposBody;

        spinner.succeed(message);

        prompt([
          {
            type: 'list',
            message: 'Please select repository that you want to list subscribers',
            name: 'repoName',
            choices: data.map((repo) => repo.name),
          },
        ]).then(async ({ repoName }) => {
          try {
            spinner.start('Saerching subscribers');

            const repo = data.find((repo) => repo.name === repoName);
            if (!repo) {
              throw new Error('Undefined repository id');
            }

            const subBody = await findAllSubscriber(repo.id!);

            spinner.succeed(subBody.message);

            subBody.data.forEach((sub) => {
              let email = chalk.red(sub.email);

              if (sub.isConfirmed) {
                email = chalk.green(sub.email);
              }
              console.log(email);
            });
          } catch (error) {
            spinner.fail(error.message);
          }
        });
      } catch (error) {
        spinner.fail(error.message);
      }
    });

export default findAllSubscribersCommand;
