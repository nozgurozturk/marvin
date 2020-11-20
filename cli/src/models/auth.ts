type Token = {
  accessToken: string;
  refreshToken: string;
};

type User = {
  name: string;
  email: string;
};

export interface IAuth {
  tokens: Token;
  user: User;
}
