import { prompt } from 'inquirer';
import { program } from 'commander';
import ora from 'ora';
import { findAllRepositories } from '../../client/service/repo/findAll';
import { findAllSubscriber } from '../../client/service/subscriber/findAll';
import { sendMailToSubscriber } from '../../client/service/subscriber/sendMail';

const spinner = ora();

const sendMailCommand = () =>
  program
    .command('send')
    .description('Sends mail to subscriber')
    .action(async () => {
      try {
        spinner.start('Searching repositories...');
        const reposBody = await findAllRepositories();
        const { data, message } = reposBody;

        if (!reposBody.data) {
          spinner.warn('No data');
          return;
        }
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
            spinner.start('Saerching subscribers...');

            const repo = data.find((repo) => repo.name === repoName);

            if (!repo) {
              throw new Error('Undefined repository');
            }

            const subBody = await findAllSubscriber(repo.id!);

            spinner.succeed(subBody.message);

            const unconfirmed = subBody.data.filter((sub) => !sub.isConfirmed);
            if (unconfirmed.length === 0) {
              spinner.warn('There is no unconfirmed subscriber in this repo');
              return;
            }
            prompt([
              {
                type: 'list',
                message: 'Please select subscriber',
                name: 'subMail',
                choices: unconfirmed.map((sub) => sub.email),
              },
            ]).then(async ({ subMail }) => {
              try {
                spinner.start('Sending...');
                const mailBody = await sendMailToSubscriber(repo.id!, subMail);
                spinner.succeed(mailBody.message);
              } catch (error) {
                spinner.fail(error.message);
              }
            });
          } catch (error) {
            spinner.fail(error.message);
          }
        });
      } catch (error) {
        spinner.fail(error.message);
      }
    });

export default sendMailCommand;
