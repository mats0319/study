// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/goc_ts
// Version: goc_ts v0.4.0

// objectToFormData 泛型用于解决'obj[key]'报错问题
export function objectToFormData<T extends object>(obj: T): FormData {
    let data: FormData = new FormData()
    for (let key in obj) {
        if (typeof obj[key] == "object") { // if field type is another object
            objectToFormData(obj[key] as object).forEach((value: FormDataEntryValue, key: string) => {
                data.append(key, value)
            })
        } else { // normal
            data.append(key, obj[key] as string)
        }
    }

    return data
}
