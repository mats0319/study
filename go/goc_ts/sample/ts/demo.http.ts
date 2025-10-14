// Generate File, Should not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/goc_ts
// Version: goc_ts v0.3.1

import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
import { ListUserRes, ListUserReq, CreateUserRes, Pagination } from "./demo.go"
import { objectToFormData } from "./utils"

class DemoAxios {
    public listUser(operator: string, page: Pagination): Promise<AxiosResponse<ListUserRes>> {
        let req: ListUserReq = {
            operator: operator,
            page: page,
        }

        return axiosWrapper.post("/user/list", objectToFormData(req))
    }

    public createUser(): Promise<AxiosResponse<CreateUserRes>> {
        return axiosWrapper.post("/user/create")
    }
}

export const demoAxios: DemoAxios = new DemoAxios()
