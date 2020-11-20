import client from '../../index';
import { APISuccess } from '../../../models/response';
import { Repo } from '../../../models/repo';;

const findAllRepositories = async (): Promise<APISuccess<[Repo]>> => {
  try {
    const body: APISuccess<[Repo]> = await client
      .get('api/repository')
      .json();
    return body;
  } catch (error) {
    return error;
  }
};

export { findAllRepositories };
