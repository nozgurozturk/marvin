import client from '../../index';
import { IAuth } from '../../../models/auth';
import { APISuccess } from '../../../models/response';


const login = async (email: string, password: string): Promise<APISuccess<IAuth>> => {
  try {
    const body: APISuccess<IAuth> = await client
      .post('auth/login', {
        json: {
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

export { login };
