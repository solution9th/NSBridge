import React, { Component } from 'react'
import { connect } from 'react-redux'
import Button from '@material-ui/core/Button';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import Divider from '@material-ui/core/Divider';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import Tooltip from '@material-ui/core/Tooltip'
import AccountBox from '@material-ui/icons/AccountBox'
import Apps from '@material-ui/icons/Apps'
import Notifications from '@material-ui/icons/Notifications'
import SwapHorizIcon from '@material-ui/icons/SwapHoriz'
import ExitToAppIcon from '@material-ui/icons/ExitToApp'
import { getData } from '../../utils/request';
import { currentUser, logout, baseUrl } from '../../api/homeApi';
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core'
import styles from '../../styles/sass/home.scss'
import { withRouter } from 'react-router-dom'

const muiStyles = () => ({
    iconFonts: {
        fontSize: '24px',
        minWidth: 0,
        width: '48px',
        height: '48px'
    }
})

class Top extends Component {
    constructor(props) {
        super(props)
        this.state = {
            user: {},
            mineAnchorEl: null, // 我的图标menu控件属性
            showUserToolTip: false,
        }
        this.mineOpen = this.mineOpen.bind(this)
        this.mineClose = this.mineClose.bind(this)
        this.goFrontPage = this.goFrontPage.bind(this)
    }

    mineOpen (event) {
        event.preventDefault()
        this.setState({
            mineAnchorEl: event.currentTarget,
            showUserToolTip: false
        })
    }

    mineClose () {
        getData(logout)
            .then(res => {
                if(res.data.err_code === 0) {
                    window.location.href = baseUrl + '/saml/login?RelayState=/'
                    this.setState({
                        mineAnchorEl: null
                    })
                }
            })
    }

    goFrontPage() {
        this.props.history.push('/')
    }
    render() {
        const { mineAnchorEl, showUserToolTip } = this.state
        const { classes } = this.props
        return (
            <div className={styles.dnsTopContainer}>
                <div className={styles.dnsTopLogo}>
                    <span>DNS</span>
                </div> 
                <div className={styles.dnsTopRight}>
                    <div>
                        <Tooltip 
                            title="个人"
                            open={showUserToolTip}
                            disableFocusListener
                            onOpen={()=>{this.setState({showUserToolTip: true})}}
                            onClose={()=>{this.setState({showUserToolTip: false})}}
                        >
                            <Button
                                aria-owns={mineAnchorEl ? 'mine-menu' : null}
                                aria-haspopup="true"
                                onClick={this.mineOpen}
                                className={classes.iconFonts}
                                >
                                <AccountBox />
                            </Button>
                        </Tooltip>
                        <Menu
                            id="mine-menu"
                            className={styles.dnsTRMine}
                            anchorEl={mineAnchorEl}
                            open={Boolean(mineAnchorEl)}
                            onClose={()=>{this.setState({mineAnchorEl: null})}}
                            >
                            <MenuItem>
                                <div>
                                    欢迎,{this.props.user.user_name}
                                </div>
                            </MenuItem>
                            <MenuItem onClick={this.mineClose}><ExitToAppIcon />退出</MenuItem>
                        </Menu>
                    </div> 
                </div>
            </div>
        )
    }
}

Top.propTypes = {
    classes: PropTypes.object,
    user: PropTypes.object
}

const mapStateToProps = state => {
    return {
        user: state.homeReducer.userInfo
    }
}

export default withRouter(connect(mapStateToProps)(withStyles(muiStyles)(Top)))
