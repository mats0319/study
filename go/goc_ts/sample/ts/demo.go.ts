// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/goc_ts
// Version: goc_ts v1.0.0

import { ResBase } from "./common.go"

export class Pagination {
    page_num: number = 0;
    page_size: number = 0;
}

export class ListUserReq {
    operator: string = "";
    list_identify: UserIdentify = 0;
    page: Pagination = new Pagination();
}

export class ListUserRes {
    res: ResBase = new ResBase();
    summary: number = 0;
    users: Array<string> = new Array<string>();
}

export class CreateUserReq {
}

export class CreateUserRes {
}

export enum UserIdentify {
    Value0 = 0,
    Value1 = 1,
    Value2 = 2,
}
