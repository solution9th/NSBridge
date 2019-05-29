import { combineReducers } from 'redux'
import commonReducer from './commonReducer'
import homeReducer from './homeReducer'

const reducers = combineReducers({
    commonReducer: commonReducer,
    homeReducer: homeReducer
})

export default reducers