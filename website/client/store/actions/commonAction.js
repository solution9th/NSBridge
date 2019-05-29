import { LOADING } from '../actionTypes'

export function loading(payload) {
    return {
        type: LOADING,
        loading: payload
    }
}