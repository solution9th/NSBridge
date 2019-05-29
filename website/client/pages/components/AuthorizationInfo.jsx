import React, { Component } from 'react'
import { withStyles } from '@material-ui/core'
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle'
import PropTypes from 'prop-types'
import Button from '@material-ui/core/Button'
import { postData, getData } from '../../utils/request';
import { authorizationInfo } from '../../api/authorizationApi';
import Snackbar from '@material-ui/core/Snackbar'
import CopeToClipboard from 'react-copy-to-clipboard'
import styles from '../../styles/sass/common.scss'


const muiStyles = () => ({
    progress: {
        margin: '10vh 0 0 49%',
        color: '#4A90E2'
    },
    iconFont: {
        fontSize: '20px'
    },
   closeBtn: {
        padding: 0,
        '&>h2': {
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            fontSize: '16px',
            fontWeight: 600,
            color: "#111111",
        }

    },
    rightButton: {
        color: '#4A90E2',
        fontSize: '14px'
       },
   search: {
       fontSize: '20px'
   },
   inputStyle: {
     flex: 1
   }
})

class AuthorizationInfo extends Component {
    constructor(props) {
        super(props) 
        this.state = {
            snackbarStatus: false,
            snackbarMessage: '新增成功',
            copyRight: false,
            row: {}, // 当前authorization的信息
            data: { domain_key: 'KKKKKK', domain_secret: 'SSSSSSSSSSSS' }
        }
    }

    setData(row) {
        this.setState({
            row: row
        })
    }

    handleClose () {
        this.setState({
            value: '',
            valueLength: 0,
            checkValue: false
        })
        this.props.close('cancle')
    }

    closeSnackbar() {
        this.setState({
            snackbarMessage: '',
            snackbarStatus: false
        })
    }

    getDomainInfo() {
        getData(authorizationInfo,{domain_key: this.state.row.domain_key})
        .then(res => {
            if(res.data.err_code === 0) {
                this.setState({
                   data: res.data.data[0]
                })
            }else {
                this.setState({
                    data: {}
                 })
            }
        })
    }

    copyData(flag) {
        // flag 为 k / s  如果是k则是复制的K值 否则是S值
        if(flag === 'k') {
            this.setState({
                snackbarMessage: '已将K值复制到剪切板',
                snackbarStatus: true
            })
        }else {
            this.setState({
                snackbarMessage: '已将S值复制到剪切板',
                snackbarStatus: true
            })
        }
    }

    render() {
        const { classes } = this.props
        const { data, row } = this.state
        return (
        <React.Fragment>
            <Dialog open={this.props.open} onClose={this.handleClose.bind(this)} aria-labelledby="simple-dialog-title" className={styles.dialogContainer}>
                <DialogTitle id="simple-dialog-title" className={classes.closeBtn}><span>{this.props.title}</span></DialogTitle>
                <div className={`${styles.dialogContent} ${styles.authorDialogContent}`}>
                   <div>
                       <p>K值： <span>{row.record_key}</span></p>
                       <CopeToClipboard text={row.record_key} onCopy={this.copyData.bind(this, 'k')}>
                            <Button className={classes.rightButton}>复制</Button>
                       </CopeToClipboard>
                   </div>
                   <div>
                       <p>S值： <span>{row.record_secret}</span></p>
                       <CopeToClipboard text={row.record_secret} onCopy={this.copyData.bind(this, 's')}>
                            <Button className={classes.rightButton}>复制</Button>
                       </CopeToClipboard>
                   </div>
                </div>
                <div className={styles.dialogButtons}>
                    <Button color="secondary"  onClick={this.handleClose.bind(this)}>关闭</Button>
                  </div>
            </Dialog>
            <Snackbar
                anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
                open={this.state.snackbarStatus}
                onClose={this.closeSnackbar.bind(this)}
                autoHideDuration={1000}
                className={styles.snackbarStyle}
                ContentProps={{
                    'aria-describedby': 'message-id',
                }}
                message={<span id="message-id">{this.state.snackbarMessage}</span>}
                />
        </React.Fragment>
        )
    }
}

AuthorizationInfo.propTypes = {
    classes: PropTypes.object,
    open: PropTypes.bool,
    close: PropTypes.func,
    userData: PropTypes.array,
    title: PropTypes.string
}

export default withStyles(muiStyles)(AuthorizationInfo)