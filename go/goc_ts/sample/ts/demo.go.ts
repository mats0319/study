// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/goc_ts
// Version: goc_ts v0.4.0

import { ResBase } from "./common.go"

export enum UserIdentify {
    UserIdentify_Administrator = 0,
    UserIdentify_VIP = 1,
    UserIdentify_Visitor = 2,
}

export class Pagination {
    page_num: number = 0;
    page_size: number = 0;
}

export class ListUserReq {
    operator: string = "";
    list_identify: UserIdentify = new UserIdentify();
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
