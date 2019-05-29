import React, { Component } from 'react'
import { Route , Switch} from 'react-router-dom'
import { connect } from 'react-redux' 
import { GET_USER_INFO } from '../../store/actionTypes'
import { getUserInfo } from '../../store/actions/homeAction'
import styles from '../../styles/sass/home.scss'
import { getData } from '../../utils/request'
import { userInfo } from '../../api/homeApi'
import TopPage from '../components/top'
import SidePage from '../components/side'
import DomainList from '../domainList/DomainList'
import AnalysisRecord from '../analysisRecord/AnalysisRecord'
import AuthorizationList from '../authorizationList/AuthorizationList'

 class Home extends Component {
    constructor(props) {
        super(props)
    }
    componentDidMount(){
        const { getUserInfo } = this.props
        getData(userInfo)
        .then(res => {
            if(res.data.err_code === 0) {
                    getUserInfo(res.data.data)
            }else {
                getUserInfo(GET_USER_INFO, {user_id: '', user_name: ''})
            }
        })
    }
    render() {
        return (
            <div className={styles.dnsContainer}> 
                <div className={styles.dnsTop}><TopPage/></div>
                <div className={styles.dnsBottom}>
                    <div className={styles.dnsSide}><SidePage /></div>
                    <div className={styles.dnsContent}>
                        <Route path="/" exact component={DomainList}/>
                        <Route path="/analysis/:name/:id"  component={AnalysisRecord}/>
                        <Route path="/authorization" exact component={AuthorizationList} />
                    </div>
                </div>
            </div>
        )
    }
}

const mapStateToProps = state => {
    return {}
}

const mapDispathToProps = dispatch => {
    return {
        getUserInfo: (data) => {dispatch(getUserInfo(data))}
    }
}

export default connect(mapStateToProps, mapDispathToProps)(Home)