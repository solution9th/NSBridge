import React, { Component } from 'react'
import { withStyles } from '@material-ui/core'
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle'
import PropTypes from 'prop-types'
import Button from '@material-ui/core/Button'
import { postData, getData, putData } from '../../utils/request';
import { updateAuthorization } from '../../api/authorizationApi';
import Snackbar from '@material-ui/core/Snackbar'
import Input from '@material-ui/core/Input'
import InputAdornment from '@material-ui/core/InputAdornment'
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

class AddKSDialog extends Component {
    constructor(props) {
        super(props) 
        this.state = {
            snackbarStatus: false,
            snackbarMessage: '新增成功',
            currentAuthorization: {}, // 当前的授权信息的数据
            value: '', // input框输入的内容
            valueLength: 0 , // input框输入内容的长度
            checkValue: false, // 判断是否输入的内容符合条件
        }
        this.changeValue = this.changeValue.bind(this)
        this.setData = this.setData.bind(this)
    }
    changeValue(event) {
        let length = event.target.value.length
        if(length > 30) {
            this.setState({
                value: event.target.value,
                valueLength: length,
                checkValue: true
            })
        }else {
            this.setState({
                value: event.target.value,
                valueLength: length,
                checkValue: false
            })
        }
       
    }

    setData(data) {
        this.setState({
            currentAuthorization: data,
            value: data.remark
        })
    }

    handleClose () {
        this.setState({
            value: '',
            valueLength: 0,
            checkValue: false,
            currentAuthorization: {}
        })
        this.props.close()
    }

    saveBack() {
        const { value, checkValue, currentAuthorization } = this.state
        if(!checkValue) {
            if(this.props.isAddBack) {
                postData(updateAuthorization, {id: currentAuthorization.id,remark: value})
                .then(res => {
                    if(res.data.err_code === 0) {
                        this.setState({
                            snackbarStatus: true,
                            snackbarMessage: '成功新增KS',
                            value: '',
                            valueLength: 0,
                                checkValue: false,
                        })
                        this.props.close()
                    }else {
                        this.setState({
                            snackbarStatus: true,
                            snackbarMessage: res.data.err_msg
                        })
                    }
                })
                .catch(err => {
                    this.setState({
                        snackbarStatus: true,
                        snackbarMessage: '添加失败'
                    })
                })
            }else {
                putData(updateAuthorization, {id: currentAuthorization.id,remark: value})
                .then(res => {
                    if(res.data.err_code === 0) {
                        this.setState({
                            snackbarStatus: true,
                            snackbarMessage: '成功添加备注' ,
                            value: '',
                            valueLength: 0,
                                checkValue: false,
                        })
                        this.props.close()
                    }else {
                        this.setState({
                            snackbarStatus: true,
                            snackbarMessage: res.data.err_msg
                        })
                    }
                })
                .catch(err => {
                    this.setState({
                        snackbarStatus: true,
                        snackbarMessage: '添加失败'
                    })
                })
            }
           

        }
    }



    closeSnackbar() {
        this.setState({
            snackbarMessage: '',
            snackbarStatus: false
        })
    }

    render() {
        const { classes } = this.props
        const { value, valueLength, checkValue } = this.state
        return (
        <React.Fragment>
            <Dialog open={this.props.open} onClose={this.handleClose.bind(this)} aria-labelledby="simple-dialog-title" className={styles.dialogContainer}>
                <DialogTitle id="simple-dialog-title" className={classes.closeBtn}><span>{this.props.title}</span></DialogTitle>
                <div className={`${styles.dialogContent} ${styles.addKSDialogContent}`}>
                    <span>备注：</span>
                    <div className={styles.dialogInput}>
                        <Input
                            value={value}
                            error={checkValue}
                            className={classes.inputStyle}
                            onChange={this.changeValue.bind(this)}
                            endAdornment={<InputAdornment position="end" style={{color: checkValue ? '#f44336' : '#9b9b9b', fontSize: '12px'}}>{valueLength}/30</InputAdornment>}
                        />
                        {
                            checkValue ? <p className={styles.inputCheck}>输入字数过多</p> : null
                        }
                        
                    </div>
                    
                </div>
                <div className={styles.dialogButtons}>
                    <Button color="secondary" onClick={this.handleClose.bind(this)}>取消</Button>
                    <Button color="primary" className={classes.rightButton} onClick={this.saveBack.bind(this)}>保存</Button>
                </div>
            </Dialog>
            <Snackbar
                anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
                open={this.state.snackbarStatus}
                onClose={this.closeSnackbar.bind(this)}
                autoHideDuration={5000}
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

AddKSDialog.propTypes = {
    classes: PropTypes.object,
    open: PropTypes.bool,
    close: PropTypes.func,
    title: PropTypes.string,
    isAddBack: PropTypes.bool
}

export default withStyles(muiStyles)(AddKSDialog)