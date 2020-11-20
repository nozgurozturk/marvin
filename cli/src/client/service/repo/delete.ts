import client from '../../index';
import { APISuccess } from '../../../models/response';

const deleteRepository = async (repoID: string): Promise<APISuccess<null>> => {
  try {
    const body: APISuccess<null> = await client
      .delete('api/repository', {
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

export { deleteRepository };
