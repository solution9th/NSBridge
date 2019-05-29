import { LOADING } from '../actionTypes'

const initState = {
    loading: 'false'
}

function commonReducer(state=initState, action) {
    switch(action.type) {
        case LOADING: 
            return {...state, loading: action.loading}
        default:
            return state
    }
}

export default commonReducer