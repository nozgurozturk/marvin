import got, { Got } from 'got';
import { getAuthConfig } from '../utils/config/auth';

let token:string;

const setBearerToken = ():string => {
  if (token) {
    return token
  }
  const auth = getAuthConfig()
  if (auth) {
    token = auth.accessToken
    return token
  }
  return ""
}

const client: Got = got.extend({
  prefixUrl: process.env.API_HOST || 'http://localhost:8081/',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `bearer ${setBearerToken()}`
  },
});

export default client;
