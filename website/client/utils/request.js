import axios from 'axios'

export function getData(url, data) {
    return axios.get(url, { params: data })
}

export function putData(url, data) {
    return axios.put(url, data)
}

export function postData(url, data) {
    const options = {
        method: 'POST',
        url: url,
        headers: {'Content-Type': 'application/json,charset=utf-8'},
        data: data
    }
    return axios(options)
}

export function deleteData(url, params) {
   return axios.delete(url, {params})
}