import { prompt } from 'inquirer';
import { program } from 'commander';
import ora from 'ora';
import { findAllRepositories } from '../../client/service/repo/findAll';
import { deleteRepository } from '../../client/service/repo/delete';

const spinner = ora();

const deleteRepoCommand = () =>
  program
    .command('delete')
    .description('Deletes repository')
    .action(async () => {
      try {
        spinner.start('Searching');

        const reposBody = await findAllRepositories();
        const { data } = reposBody;

        spinner.succeed(reposBody.message);

        prompt([
          {
            type: 'list',
            message: 'Please select repository that you want to delete',
            name: 'repoName',
            choices: data.map((repo) => repo.name),
          },
        ]).then(async ({ repoName }) => {
          try {
            const repo = data.find((repo) => repo.name === repoName);

            if (!repo) {
              throw new Error('Undefined repository');
            }

            spinner.start('Deleting');

            const deleteBody = await deleteRepository(repo.id!);

            spinner.succeed(deleteBody.message);
          } catch (error) {
            spinner.fail(error.message);
          }
        });
      } catch (error) {
        spinner.fail(error.message);
      }
    });

export default deleteRepoCommand;
