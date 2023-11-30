export class ListReq {
  public cursor?: string;
  public q?: string;
  public limit?: number;
  public ids?: Array<string>;
}
