import client from '../../index';
import { APISuccess } from '../../../models/response';
import { Repo } from '../../../models/repo';

const updateRepositoryPackages = async (repoID: string): Promise<APISuccess<Repo>> => {
  try {
    const body: APISuccess<Repo> = await client
      .put('api/repository', {
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

export { updateRepositoryPackages };
