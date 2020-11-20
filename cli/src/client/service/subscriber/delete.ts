import client from '../../index';
import { APISuccess } from '../../../models/response';

const deleteSubscriber = async (id: string): Promise<APISuccess<null>> => {
  try {
    const body: APISuccess<null> = await client
      .delete('api/subscriber', {
        json: {
          id: id,
        },
      })
      .json();
    return body;
  } catch (error) {
    return error;
  }
};

export { deleteSubscriber };
