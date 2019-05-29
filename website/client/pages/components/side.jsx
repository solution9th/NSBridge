import React, { Component } from 'react'
import { withStyles } from '@material-ui/core/styles';
import { connect } from 'react-redux'
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import PropTypes from 'prop-types'
import Language from '@material-ui/icons/Language';
import Assignment from '@material-ui/icons/Assignment'
import { withRouter } from 'react-router-dom'
import { saveSearchValue } from '../../store/actions/homeAction'

const muiStyles = () => ({
    buttonIcon: {
        fontSize: '16px',
        marginRight: '4px'
    },
    sideList: {
        fontSize: "16px"
    },
    activeList: {
        backgroundColor: '#e1e1e1',
        '&:hover': {
            backgroundColor: '#ececec'
        }
    },
    hovers: {
        '&:hover': {
            backgroundColor: '#ececec'
        }
    },
    sizeStyle: {
        minWidth: '200px',
        height: '49px',
        paddingLeft: '24px'
    }
})

class Side extends Component {
    constructor(props) {
        super(props)
        this.showDomain = this.showDomain.bind(this)
        this.showAuthorization = this.showAuthorization.bind(this)
        this.state={
            activeDomain: true
        }
    }

    componentDidMount(){
        if(/authorization/g.test(this.props.history.location.pathname)) {
            this.setState({
                activeDomain: false
            })
        }else {
            this.setState({
                activeDomain: true
            })
        }
    }

    showDomain () {
        this.setState({
            activeDomain: true
        })
        this.props.history.push('/')
        this.props.saveSearchValue('')
    }

    showAuthorization () {
        this.setState({
            activeDomain: false
        })
        this.props.history.push('/authorization')
    }

    render() {
        const { classes } = this.props
        return (
            <div>
                <List component="nav">
                    <ListItem button onClick={this.showDomain} className={`${this.state.activeDomain ? classes.activeList : classes.hovers} ${classes.sizeStyle}`}>
                        <ListItemIcon className={classes.buttonIcon}>
                            <Language />
                        </ListItemIcon>
                        <span>域名列表</span>       
                    </ListItem>
                    <ListItem button onClick={this.showAuthorization}  className={`${!this.state.activeDomain ? classes.activeList : classes.hovers} ${classes.sizeStyle}`}>
                        <ListItemIcon  className={classes.buttonIcon}>
                            <Assignment />
                        </ListItemIcon>
                        <span>授权信息</span>     
                    </ListItem>
                </List>
            </div>
        )
    }
}

Side.propTypes = {
    classes: PropTypes.object.isRequired,
    saveSearchValue: PropTypes.func
}

const mapStateToProps = state => {
    return {}
}

const mapDispatchToProps = dispatch => {
    return {
        saveSearchValue: (data)=>{dispatch(saveSearchValue(data))}
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(withStyles(muiStyles)(Side)))