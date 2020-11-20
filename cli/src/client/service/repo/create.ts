import client from '../../index';
import { APISuccess } from '../../../models/response';
import { Repo } from '../../../models/repo';

const createRepository = async (url: string): Promise<APISuccess<Repo>> => {
  try {
    const body: APISuccess<Repo> = await client
      .post('api/repository', {
        json: {
          url: url,
        },
      })
      .json();
    return body;
  } catch (error) {
    return error;
  }
};

export { createRepository };
