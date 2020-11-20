import { program } from 'commander';
import ora from 'ora';
import { findAllRepositories } from '../../client/service/repo/findAll';

const spinner = ora();

const listRepoCommand = () =>
  program
    .command('list-repo')
    .description('List all repositories of user')
    .action(async () => {
      try {
        spinner.start('Searching\n');

        const body = await findAllRepositories();
        const { data, message } = body;

        data.forEach((repo) => {
          console.log(repo.name);
        });

        spinner.succeed(message);
      } catch (error) {
        spinner.fail(error.message);
      }
    });

export default listRepoCommand;
