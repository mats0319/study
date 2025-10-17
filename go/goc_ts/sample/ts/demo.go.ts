// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/goc_ts
// Version: goc_ts v0.4.0

import { ResBase } from "./common.go"

export class Pagination {
    page_num: number = 0;
    page_size: number = 0;
}

export class ListUserReq {
    operator: string = "";
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
