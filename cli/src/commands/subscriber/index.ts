import createSubscriber from './create';
import findAllSubscribers from './findAll';
import sendMail from './sendEmail';

export const subscriber = () => {
  createSubscriber();
  findAllSubscribers();
  sendMail();
};
