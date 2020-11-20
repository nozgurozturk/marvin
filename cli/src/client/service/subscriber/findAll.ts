import client from '../../index';
import { APISuccess } from '../../../models/response';
import { Subscriber } from '../../../models/subscriber';

const findAllSubscriber = async (repoID: string): Promise<APISuccess<[Subscriber]>> => {
  try {
    const body: APISuccess<[Subscriber]> = await client
      .post('api/subscriber/all', {
        json: {
          id: repoID,
        },
      })
      .json();
    return body;
  } catch (error) {
    return error;
  }
};

export { findAllSubscriber };
