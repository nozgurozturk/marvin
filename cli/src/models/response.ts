export interface APISuccess<T> {
  status: number;
  message: string;
  data: T;
}

export interface APIError {
  status: number;
  error: string;
  message: string;
}
