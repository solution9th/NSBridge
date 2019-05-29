import React, { Component } from 'react'
import DnsTitle from '../components/DnsTitle'
import DnsSearch from '../components/DnsSearch'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import AddKSDialog from '../components/AddKSDialog'
import ToolTip from '@material-ui/core/Tooltip'
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import Button from '@material-ui/core/Button'
import MoreVert from '@material-ui/icons/MoreVert'
import Create from '@material-ui/icons/Create'
import Delete from '@material-ui/icons/Delete'
import Switch from '@material-ui/core/Switch';
import { getData , deleteData, postData, putData } from '../../utils/request'
import { authorizationInfo, deleteAuthorization, startEndAuthorization } from '../../api/authorizationApi'
import CircularProgress from '@material-ui/core/CircularProgress'
import Snackbar from '@material-ui/core/Snackbar'
import styles from '../../styles/sass/common.scss'
import { Divider } from '@material-ui/core';
import CopyToClipboard from 'react-copy-to-clipboard'

const muiStyles = () => ({
    tableHeader: {
        color: '#d2d2d2',
        tr: {
            th: {fontSize: '1rem', boxSizing: 'border-box'}
        },
        '&>tr>th:nth-child(1),&>tr>th:nth-child(2),&>tr>th:nth-child(4),&>tr>th:nth-child(5)': {
            width: '17.5%'
        },
        '&>tr>th:nth-child(3)': {
            width: '14%'
        },
        '&>tr>th:nth-child(6),&>tr>th:nth-child(7)': {
            width: '8%'
        }
    },
    tableBody: {
        color: '#4a4a4a',
        tr: {
            td: {fontSize: '1rem', boxSizing: 'borer-box'}
        },
        '&>tr>td:nth-child(1),&>tr>td:nth-child(2),&>tr>td:nth-child(4),&>tr>td:nth-child(5)': {
            width: '17.5%'
        },
        '&>tr>td:nth-child(3)': {
            width: '14%'
        },
        '&>tr>td:nth-child(6),&>tr>td:nth-child(7)': {
            width: '8%'
        }
    },
    status: {
        display: 'flex',
    },
    switch: {
        color: '#3986ff'
    },
    progress: {
        margin: '10vh 0 0 49%',
        color: '#4A90E2'
    },
    toolStyle: {
        fontSize: '12px'
    },
    menuItem: {
        color: '#6a6a6a',
        fontSize: '12px',
        '&>svg': {
            fontSize: '20px',
            width: '20px',
            height: '20px',
            marginRight: '12px'
        },
        '&:hover': {
            backgroundColor: '#ececec'
        }
    },
    hrStyle: {
        backgroundColor: '#ececec'
    }
})

class AuthorizationList extends Component {
    constructor(props) {
        super(props)
        this.state = {
            title: [{label: '授权信息', url: ''}],
            btnText: 'authorization',
            tableHeader: [
                {label: 'K值', key: 'templateName'},
                {label: 'S值', key: 'httpType'},
                {label: '备注', key: 'requestType'},
                {label: '创建时间', key: 'createTime'},
                {label: '最后修改时间',key: 'modifyTime'},
                {label: '状态',key: 'status'},
                {label: '操作', key: 'cancle'},
            ],
            snackbarStatus: false,
            snackbarMessage: '',
            snackbarDuration: 1000,
            sortTable: 'up',
            recordType: -1,
            tableData: [],
            statusAnchorEl: null,
            editAnchorEl: null,
            options: [
                {label: '不限分类', value: -1},
                {label: '已开启', value: 1},
                {label: '未开启', value: 0},
            ],
            currentRow: '-1', // 当前行
            pageCount: 30, 
            openDialog: false, // 是否展示模板dialog
            isAddBack: false, // 打开弹窗时是否是添加备注
            searchValue: '', // 搜索框的值
            templateStatus: '0', // 模板状态
            dialogTitle: '新增KS', // 弹框的title
            currentAuthorizationData: {}, // 当前授权信息的数据
            dialogTemplateData: {}, // 点击编辑的时候传给dialog数据
            showCircle: false,
        }
        this.id = 0,
        this.canScroll = true,
        this.currentPage = 0,
        this.localData = [] // 用来存储当前所有信息
        this.currentPageListLength = 0 // 用来判断当前页取回的数据是否达到每一页请求数据条数
    }

    componentDidMount() {
        this.getAuthorization()
    }

     // 下拉加载更多
     onScroll() {
        let that = this
        if(!this.canScroll) {
            return
        }
        this.canScroll = false
        setTimeout(() => {
            that.scrollData()
            that.canScroll = true
        }, 300)
    }

    scrollData() {
        let scrollHeight = this.container.scrollHeight
        let clientHeight = this.container.clientHeight
        let scrollTop = this.container.scrollTop
        if(scrollHeight === (clientHeight + scrollTop)){
            if(this.currentPageListLength === this.state.pageCount) {
                this.currentPage = this.currentPage +  this.state.pageCount ;
                this.getAuthorization('scroll')
            }else {
                this.setState({
                    snackbarMessage: '暂无更多数据',
                    snackbarStatus: true,
                    snackbarDuration: 5000,
                })
            }
        }
    }


    // 获取模板列表数据
    getAuthorization(flag) {
        this.setState({
            showCircle: true
        })
        const { pageCount, searchValue, tableData, recordType } = this.state
        if(flag && flag === 'scroll') {
            
        }else {
            this.currentPage = 0
        }

        let params = {
            count: pageCount, // 分页条数 
            offset: this.currentPage, // 第几页
            domain_key: searchValue !== '' ? searchValue : '', // 搜索条件
            disable: recordType === -1 ? '' : recordType, // 是否禁用(1禁用 ， 0启用)
        }

        getData(authorizationInfo, params)
            .then(res => {
                if(res.data.err_code === 0 && res.data.data !== null) {
                    this.currentPageListLength = res.data.data.length
                    if(flag && flag ==='scroll') {
                        let data = tableData.slice()
                        let newTableData = []
                        res.data.data.forEach(item => {
                            data.push(item)
                            this.localData.push(item)
                        })
                        // 开启状态过滤--开始
                        if(recordType !== -1) {
                            data.forEach((item) => {
                                if(item.disable === recordType) {
                                    newTableData.push(item)
                                }
                            })
                            this.setState({
                                tableData: newTableData,
                                showCircle: false
                            })
                        }else {
                            this.setState({
                                tableData: data,
                                showCircle: false
                            })
                        }
                       // 开启状态过虑--结束 
                    }else {
                        this.localData = res.data.data
                        let newTableData = []
                         // 开启状态过滤--开始
                         if(recordType !== -1) {
                            res.data.data.forEach((item) => {
                                if(item.disable === recordType) {
                                    newTableData.push(item)
                                }
                            })
                            this.setState({
                                tableData: newTableData,
                                showCircle: false
                            })
                        }else {
                            this.setState({
                                tableData: res.data.data,
                                showCircle: false
                            })
                        }
                       // 开启状态过虑--结束 
                    }
                }else {
                    this.setState({
                        tableData: [],
                        showCircle: false
                    })
                }
            })
    }

    sortTable() {
        let arrow = this.state.sortTable === 'up' ? 'down' : 'up'
        let temp = this.state.tableData.slice()
        let newTableData = temp.reverse()
        this.setState({
            sortTable: arrow,
            tableData: newTableData
        })
    }

    deleteRow(){
        deleteData(deleteAuthorization + '/' + this.state.currentAuthorizationData.id)
            .then(res => {
                if(res.data.err_code === 0) {
                    this.getAuthorization()
                    this.setState({
                        editAnchorEl: null,
                        snackbarMessage: '删除成功',
                        snackbarStatus: true,
                        snackbarDuration: 5000,
                    })
                }else {
                    this.setState({
                        snackbarMessage: res.data.err_msg,
                        snackbarStatus: true,
                        snackbarDuration: 5000
                    })
                }
            })
    }

    openMenu(event) {
        this.setState({
            statusAnchorEl: event.currentTarget
        })
    }
    
    closeMenu() {
        this.setState({
            statusAnchorEl: null
        })
    }

    addBack() {
        event.preventDefault()
        this.setState({
            editAnchorEl: null,
            openDialog: true,
            dialogTitle: '备注',
            isAddBack: false
        },() => {
            this.dialogRef.setData(this.state.currentAuthorizationData)
        })
    }

    onMouseover(id) {
        this.setState({
            currentRow: id
        })
    }

    onMouseout() {
        this.setState({
            currentRow: '-1'
        })
    }

    openAddKSDialog () {
        this.setState({
            openDialog: true,
            dialogTitle: '新增KS',
            isAddBack: true
        })
    }

    closeDialog (){
        this.setState({
            openDialog: false,
            currentAuthorizationData: {}
        },()=>{
            this.getAuthorization()
        })
       
        
    }

    searchTems(input) {
        this.setState({
            searchValue: input
        },()=> {
            this.getAuthorization()
        })
    }

    showTemplateMenu(rowData,event) {
        this.setState({
            editAnchorEl: event.currentTarget,
            currentAuthorizationData: rowData
        })
    }

    changeStatus(event) {
        let status = event.target.checked ? 1 : 0
        putData(startEndAuthorization+'/'+this.state.currentAuthorizationData.id)
            .then((res) => {
                if(res.data.err_code === 0) {
                    this.setState({
                        snackbarMessage: '修改状态成功',
                        snackbarStatus: true,
                        snackbarDuration: 5000,
                        editAnchorEl: null
                    })
                    this.getAuthorization()                    
                }else {
                    this.setState({
                        snackbarMessage: res.data.err_msg,
                        snackbarStatus: true,
                        snackbarDuration: 5000,
                        editAnchorEl: null
                    })
                }
            })
    }

    closeSnackbar() {
        this.setState({
            snackbarMessage: '',
            snackbarStatus: false,
            snackbarDuration: 5000
        })
    }

    selectRecordType(value) {
        this.setState({
            recordType: value
        },()=>{
            this.getAuthorization()
        })
    }

    copyToClipboard(flag) {
        if(flag === 'k') { // 判断是复制的域名
            this.setState({
                snackbarMessage: '已将K值复制到剪切板',
                snackbarStatus: true,
                snackbarDuration: 1000
            })
        }else {
            this.setState({
                snackbarMessage: '已将S值复制到剪切板',
                snackbarStatus: true,
                snackbarDuration: 1000
            })
        }
    }

    render() {
        const { title, btnText, tableHeader,snackbarDuration, tableData, currentRow, statusAnchorEl, options,isAddBack, editAnchorEl, openDialog, dialogTitle, currentAuthorizationData, dialogTemplateData, sortTable, showCircle  } = this.state
        const { classes } = this.props
        return (
            <React.Fragment>
            <AddKSDialog open={openDialog} close={this.closeDialog.bind(this)} title={dialogTitle} isAddBack={isAddBack} innerRef={(element) => {this.dialogRef = element}}/>
            <div className={styles.dnsRightContainer} ref={(container) => {this.container = container}} onScroll={this.onScroll.bind(this)}>
                <div className={styles.dnsTitle}><DnsTitle title={title} /></div>
                <div className={styles.dnsSearch}>
                    <DnsSearch 
                        btnText={btnText} 
                        managerSearch={false} 
                        openDialog={this.openAddKSDialog.bind(this)} 
                        selectRecordType={this.selectRecordType.bind(this)}
                        search={this.searchTems.bind(this)}/>
                </div>
                <div className={styles.dnsTable}> 
                <Table>
                    <TableHead className={classes.tableHeader}>
                        <TableRow>
                        {
                            tableHeader.map((item, index) => {
                                if(item.key === 'cancle'){
                                    return (
                                        <TableCell key={index}  className={classes.cancleColumn}>
                                            {item.label}
                                        </TableCell>
                                    )
                                }else {
                                    return (
                                        <TableCell key={index}>
                                            {item.label}
                                        </TableCell>
                                    )
                                }
                            })
                        }
                        </TableRow>
                    </TableHead>
                    <TableBody className={classes.tableBody}>
                        {
                            tableData.length>0 ?(tableData.map((item, index) => {
                                return (
                                    <TableRow 
                                        key={item.id.toString() + index} 
                                        className={classes.tableRow}
                                        hover
                                        onMouseEnter={this.onMouseover.bind(this,item.id)}
                                        onMouseLeave= {this.onMouseout.bind(this)}
                                        >
                                        <TableCell>
                                            
                                            {
                                                item.domain_key.length > 10 ?
                                                    (
                                                        <ToolTip title={item.domain_key} className={classes.toolStyle}>
                                                            <span>{item.domain_key.substr(0,10)+'...'}</span>
                                                        </ToolTip>
                                                    ) :
                                                    (<span>{item.domain_key}</span>)
                                            }
                                            {
                                                currentRow === item.id ? 
                                                    <CopyToClipboard text={item.domain_key} onCopy={this.copyToClipboard.bind(this, 'k')}>
                                                        <button className={styles.tableLineButton}>复制</button>
                                                    </CopyToClipboard> : null
                                            }
                                            
                                        </TableCell>
                                        <TableCell>
                                            {
                                                item.domain_secret.length > 10 ?
                                                    (
                                                        <ToolTip title={item.domain_secret} className={classes.toolStyle}>
                                                            <span>{item.domain_secret.substr(0,10)+'...'}</span>
                                                        </ToolTip>
                                                    ) :
                                                    (<span>{item.domain_secret}</span>)
                                            }
                                            {
                                                currentRow === item.id ? 
                                                    <CopyToClipboard text={item.domain_secret} onCopy={this.copyToClipboard.bind(this, 's')}>
                                                        <button className={styles.tableLineButton}>复制</button>
                                                    </CopyToClipboard> : null
                                            }
                                            
                                        </TableCell>
                                        <TableCell>
                                            {
                                                item.remark.length > 10 ?
                                                    (
                                                        <ToolTip title={item.remark} className={classes.toolStyle}>
                                                            <span>{item.remark.substr(0,10)+'...'}</span>
                                                        </ToolTip>
                                                    ) :
                                                    (<span>{item.remark}</span>)
                                            }
                                        </TableCell>
                                        <TableCell>{item.create_at ? item.create_at.replace(/(T|Z)/ig, '  ') : ''}</TableCell>
                                        <TableCell>{item.update_at ? item.update_at.replace(/(T|Z)/ig, '  ') : ''}</TableCell>
                                        <TableCell>{item.disable === 1 ? <span style={{color: '#d0021b'}}>停用</span> : <span style={{color: '#4a90e2'}}>启用</span>}</TableCell>
                                        <TableCell  
                                            className={classes.cancleColumn} 
                                        >
                                            <ToolTip title="操作" disableFocusListener>
                                                <Button
                                                    aria-owns={editAnchorEl ? 'template-menu' : null}
                                                    aria-haspopup="true"
                                                    onClick={this.showTemplateMenu.bind(this,item)}
                                                    className={styles.dnsTableButton}
                                                    >
                                                    <MoreVert />
                                                </Button>
                                            </ToolTip>
                                        </TableCell>
                                    </TableRow>
                                )
                            })): null
                        }
                    </TableBody>
                </Table>
                {
                    showCircle ? <CircularProgress className={classes.progress} size={60}/> : null
                }
                <Menu
                    id={"template-menu"}
                    anchorEl={editAnchorEl}
                    open={Boolean(editAnchorEl)}
                    onClose={()=>{this.setState({editAnchorEl: null})}}
                    className={styles.dnsSelectMenu}
                >
                    {
                        currentAuthorizationData.disable === 1 ? 
                            <MenuItem onClick={this.deleteRow.bind(this)} className={classes.menuItem}><Delete />删除</MenuItem> : null
                    }
                    
                    <MenuItem onClick={this.addBack.bind(this)} className={classes.menuItem}><Create />备注</MenuItem>
                    <Divider className={classes.hrStyle}/>
                    <MenuItem className={classes.menuItem}>
                        状态
                        <Switch
                            checked={currentAuthorizationData.disable === 0}
                            onChange={this.changeStatus.bind(this)}
                            className={currentAuthorizationData.disable === 0 ? styles.switchButtonActive : styles.switchButton}
                            />
                    </MenuItem>
                </Menu>
                </div>
            </div>
            <Snackbar
                anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
                open={this.state.snackbarStatus}
                onClose={this.closeSnackbar.bind(this)}
                autoHideDuration={snackbarDuration}
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

AuthorizationList.propTypes = {
    classes: PropTypes.object.isRequired
}

export default withStyles(muiStyles)(AuthorizationList)