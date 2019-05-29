// 被多个reducer用到的action
export const LOADING = 'LOADING'

/**
 * loginPage used actions
 */
export const SAVE_NAME = 'SAVE_NAME';
export const SAVE_PASS = 'SAVE_PASS';

// 异步action types
export const ASYNC_REQUEST = 'ASYNC_REQUEST';
export const ASYNC_SUCCESS = 'ASYNC_SUCCESS';
export const ASYNC_FAILED = 'ASYNC_FAILED';

/**
 * 以上的都是例子
 *  */

// 获取用户信息
export const GET_USER_INFO = 'GET_USER_INFO'

// 保存域名列表中的搜索框的值， 为了从解析记录跳过来以后还能回显
export const SAVE_SEARCH_VALUE = 'SAVE_SEARCH_VALUE'