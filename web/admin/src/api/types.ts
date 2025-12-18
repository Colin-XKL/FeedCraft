export interface APIResponse<T = any> {
  status?: number;
  msg: string;
  code: number;
  data: T;
}
