import React, { Component } from 'react'
import styles from '../../styles/sass/common.scss'
import PropTypes from 'prop-types'
import { withRouter } from 'react-router-dom'
import NavigateNext from '@material-ui/icons/NavigateNext'

 class DnsTitle extends Component {
    constructor(props) {
        super(props)
    }
    goPage(url) {
        if(url !== '') {
            this.props.history.push(url)
        }
    }
    render() {
        const { title } = this.props
        const titleLength = title.length - 1
        return (
            <React.Fragment>
                {
                    title.map((item,index)=>{
                        if(index === titleLength) {
                            return  <p key={index} onClick={this.goPage.bind(this, item.url)}><span style={{display:'inline-block', padding: '0 10px'}}>{item.label}</span></p>
                        }else {
                            return <p key={index} style={{cursor: 'pointer', display: 'flex', alignItems: 'center'}} onClick={this.goPage.bind(this, item.url)}><span className={styles.navHover}>{item.label}</span> <NavigateNext /> </p>
                        }
                    })
                }
               
            </React.Fragment>
        )
    }
}

export default withRouter(DnsTitle)