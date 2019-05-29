// action 中的内容暂时没用  作为一个例子参考
import { SAVE_NAME, SAVE_PASS, ASYNC_REQUEST, ASYNC_SUCCESS, ASYNC_FAILED, LOADING, GET_USER_INFO, SAVE_SEARCH_VALUE } from '../actionTypes'
import { getData, postData } from '../../utils/request'
import { loading } from './commonAction'

export function getUserInfo(payload) {
    return {
        type: GET_USER_INFO,
        payload: payload
    }
}

export function saveSearchValue (payload) {
    return {
        type: SAVE_SEARCH_VALUE,
        payload: payload
    }
}


export function saveName(payload) {
    return {
        type: SAVE_NAME,
        name: payload
    }
}

export function savePass(payload) {
    return {
        type: SAVE_PASS,
        password: payload
    }
}

export function asyncRequest() {
    return {
        type: ASYNC_REQUEST
    }
}

export function asyncSuccess(payload) {
    return {
        type: ASYNC_SUCCESS,
        payload
    }
}

export function asyncFailed(payload) {
    return {
        type: ASYNC_FAILED,
        payload
    }
}

export function asyncGetData(url, data) {
    return (dispatch) => {
        dispatch(loading('true'));
        return getData(url,data).then((res) => {
            dispatch(asyncSuccess(res))
            dispatch(loading('false'));
        }).catch(err => {
            dispatch(asyncFailed(err))
            dispatch(loading('false'));
        })
    }
}

export function asyncPostData(url, data) {
    return (dispatch) => {
        dispatch(loading('true'));
        return postData(url, data).then(res => {
            dispatch(asyncSuccess(res))
            dispatch(loading('false'));
        }).catch(err => {
            dispatch(asyncFailed(err))
            dispatch(loading('false'));
        })
    }
}

