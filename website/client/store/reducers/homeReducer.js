// reducer 也是作为一个例子使用

import { SAVE_NAME, SAVE_PASS, ASYNC_REQUEST, ASYNC_SUCCESS, ASYNC_FAILED, GET_USER_INFO, SAVE_SEARCH_VALUE } from '../actionTypes'

let initState = {
    name: 'huangyulong',
    password: '123456',
    data: {},
    userInfo: {user_name: '', user_id: ''},
    searchValue: ''
}

function homeReducer(state = initState, action) {
    switch(action.type) {
        case GET_USER_INFO: 
            return {...state, userInfo: action.payload}
        case SAVE_SEARCH_VALUE: 
            return {...state, searchValue: action.payload}
        case SAVE_NAME: 
            return Object.assign({}, state, {name: action.name})
        case SAVE_PASS:
            return {...state, password: action.password}
        case ASYNC_REQUEST: 
            return { ...state } 
        case ASYNC_SUCCESS: 
            return {...state, data: action.payload}
        case ASYNC_FAILED: 
            return {...state, data: {}}
        default: 
            return state
    }
}

export default homeReducer