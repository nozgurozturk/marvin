import { prompt } from 'inquirer';
import { program } from 'commander';
import chalk from 'chalk';
import ora from 'ora';
import { findAllRepositories } from '../../client/service/repo/findAll';
import { updateRepositoryPackages } from '../../client/service/repo/update';

const spinner = ora();

const updateRepoCommand = () =>
  program
    .command('update')
    .description('Updates repository packages')
    .action(async () => {
      try {
        const reposBody = await findAllRepositories();
        const { data } = reposBody;

        prompt([
          {
            type: 'list',
            message: 'Please select repository that you want to update',
            name: 'repoName',
            choices: data.map((repo) => repo.name),
          },
        ]).then(async ({ repoName }) => {
          try {
            const repo = data.find((repo) => repo.name === repoName);
            if (!repo) {
              throw new Error('Undefined repository');
            }

            const updatedBody = await updateRepositoryPackages(repo.id!);

            const { packageList } = updatedBody.data;
            if (!packageList) {
              spinner.warn('There is no package in this repository');
              return;
            }

            packageList.forEach((pkg) => {
              let name = chalk.whiteBright(pkg.name);
              let current = chalk.green(pkg.version.current);
              const last = chalk.blue(pkg.version.last);

              if (pkg.isOutdated) {
                name = chalk.bgRed(pkg.name);
                current = chalk.red(pkg.version.current);
              }

              console.log(current, '•', last, '-', name);
            });
          } catch (error) {
            spinner.fail(error.message);
          }
        });
      } catch (error) {
        spinner.fail(error.message);
      }
    });

export default updateRepoCommand;
