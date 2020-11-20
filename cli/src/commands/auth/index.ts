import login from './login';
import setUser from './setUser';
import signup from './signup';

export const auth = () => {
  login();
  signup();
  setUser();
};
