export type PackageVersion = {
  current: string;
  last: string;
};

export type Package = {
  name: string;
  version: PackageVersion;
  file: string;
  isOutdated: boolean;
};

export type Repo = {
  id?: string;
  userID: string;
  name: string;
  owner: string;
  path: string;
  provider: string;
  packageList: [Package];
};
