import client from '../../index';
import { IAuth } from '../../../models/auth';
import { APISuccess } from '../../../models/response';

const signup = async (
  name: string,
  email: string,
  password: string,
): Promise<APISuccess<IAuth>> => {
  try {
    const body: APISuccess<IAuth> = await client
      .post('auth/signup', {
        json: {
          name: name,
          email: email,
          password: password,
        },
      })
      .json();

    return body;
  } catch (error) {
    return error;
  }
};

export { signup };
