import React from 'react'
import { BrowserRouter, Switch, Route } from 'react-router-dom'
import { connect } from 'react-redux'
import styles from '../styles/css/common.css'
import loadable from '@loadable/component' 

const Loading = () => (
    <div className={styles.loadingStyle}>
        页面正在加载中...
    </div>
)

// 动态路由（组件拆分）动态加载组件
// 参考https://reacttraining.com/react-router/web/guides/code-splitting
// 其中fallback 是在加载路由过程中展示的组件，加载完成以后就会停止显示
// const HomePage = loadable(() => import('../pages/home/home'), {fallback: Loading})
const HomePage = loadable(() => import('../pages/home/home'))

class RouterComponent extends React.Component {
    render(){
        return (
            <BrowserRouter>
                {
                    this.props.loading === 'true'?
                    <Loading /> :
                    null
                }
                <Switch>
                    <Route path="/" component={HomePage} />
                </Switch>
            </BrowserRouter>
        )
    }
}

const mapStateToProps = (state) => {
    return {
        loading: state.commonReducer.loading
    }
}

const Routers = connect(mapStateToProps)(RouterComponent) 

export default Routers