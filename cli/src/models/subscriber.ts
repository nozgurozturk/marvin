enum Frequency {
  hour = 'hour',
  day = 'way',
  week = 'week',
}

type Notify = {
  hour: number;
  minute: number;
  weekday: number;
  frequency: Frequency;
};

export type Subscriber = {
  id: string;
  repoID: string;
  email: string;
  isConfirmed: boolean;
  notify: Notify;
};
