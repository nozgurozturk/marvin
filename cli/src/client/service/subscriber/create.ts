import client from '../../index';
import { APISuccess } from '../../../models/response';
import { Subscriber } from '../../../models/subscriber';

const createSubscriber = async (repoID: string, email: string): Promise<APISuccess<Subscriber>> => {
  try {
    const body: APISuccess<Subscriber> = await client
      .post('api/subscriber', {
        json: {
          repoID: repoID,
          email: email
        },
      })
      .json();
    return body;
  } catch (error) {
    return error;
  }
};

export { createSubscriber };
