// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/goc_ts
// Version: goc_ts v1.0.0

import axios, { AxiosInstance } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
    baseURL: "http://127.0.0.1:9693",
    timeout: 3000,
});
