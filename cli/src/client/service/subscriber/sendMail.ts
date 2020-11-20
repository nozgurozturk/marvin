import client from '../../index';
import { APISuccess } from '../../../models/response';

const sendMailToSubscriber = async (repodID: string, email: string): Promise<APISuccess<null>> => {
  try {
    const body: APISuccess<null> = await client
      .post('api/subscriber/send', {
        json: {
          repoID: repodID,
          email: email,
        },
      })
      .json();
    return body;
  } catch (error) {
    return error;
  }
};

export { sendMailToSubscriber };
