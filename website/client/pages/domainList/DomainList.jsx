import React, { Component } from 'react'
import DnsTitle from '../components/DnsTitle'
import DnsSearch from '../components/DnsSearch'
import PropTypes from 'prop-types'
import styles from '../../styles/sass/common.scss'
import { domainList } from '../../api/domainApi'
import { getData, postData, deleteData } from '../../utils/request'
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableSortLabel from '@material-ui/core/TableSortLabel';
import Tooltip from '@material-ui/core/Tooltip'
import ArrowDownward from '@material-ui/icons/ArrowDownward'
import ArrowUpward from '@material-ui/icons/ArrowUpward'
import CircularProgress from '@material-ui/core/CircularProgress'
import Snackbar from '@material-ui/core/Snackbar'
import Button from '@material-ui/core/Button'
import Menu from '@material-ui/core/Menu'
import MenuItem from '@material-ui/core/MenuItem'
import MoreVert from '@material-ui/icons/MoreVert'
import AuthorizationInfo from '../components/AuthorizationInfo';
import Assignment from '@material-ui/icons/Assignment'
import Poll from '@material-ui/icons/Poll'

const muiStyles = () => ({
    tableHeader: {
        color: '#d2d2d2',
        fontSize: '1rem', 
        '&>tr>th:nth-child(4)': {
            width: '22%'
        },
        '&>tr>th:nth-child(1)': {
            width: '22%'
        },
        '&>tr>th:nth-child(2)': {
            width: '10%'
        },
        '&>tr>th:nth-child(5)': {
            width: '10%'
        },
        '&>tr>th:nth-child(3)': {
            width: '36%'
        }
    },
    tableBody: {
        color: '#4a4a4a',
        fontSize: '1rem',
        '&>tr>td:nth-child(4)': {
            width: '22%'
        },
        '&>tr>td:nth-child(1)': {
            width: '22%'
        },
        '&>tr>td:nth-child(2)': {
            width: '10%'
        },
        '&>tr>td:nth-child(5)': {
            width: '10%'
        },
        '&>tr>td:nth-child(3)': {
            width: '36%'
        }
    },
    cancleColumn: {
        width: '5rem',
        cursor: 'pointer'
    },
    progress: {
        margin: '10vh 0 0 49%',
        color: '#4A90E2'
    },
    toolTip: {
        fontSize: '12px'
    },
    arrow: {
        width: '20px',
        height: '20px',
        marginLeft: '10px'
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
    }
    

})

class DomainList extends Component {
    constructor(props) {
        super(props)
        this.state = {
            title: [{label: '域名列表', url: ''}],
            btnText: 'domain',
            tableHeader: [
                {label: '域名', key: 'id'},
                {label: '记录数', key: 'name'},
                {label: 'DNS服务器', key: 'mail'},
                {label: '添加时间', key: 'addTime'},
                {label: '操作', key: ''}
            ],
            tableData: [
                {   "id": 9, "fone_domain_id": 28, "domain_key": "", "domain": "dev.newio.cc",
                    "name_server": [  "ns1.newio.cc",  "ns2.newio.cc","ns2.newio.cc","ns2.newio.cc","ns2.newio.cc","ns2.newio.cc","ns2.newio.cc","ns2.newio.cc", ], "soa_email": "",
                    "remark": "测试",  "is_take_over": 1,  "is_open_key": 1,  "record_key": "cmVjb3Jk2cab9d2edfd19538",
                    "record_secret": "wT3gFaOxhj1nylxJQ1s8",  "create_at": "2019-03-22T17:11:55Z", "update_at": "2019-03-23T14:00:26Z"
                }
            ],
            // currentRow: '-1',
            openAuthorDialog: false, // 是否打开添加管理员dialog
            authorDialogTitle: '授权信息',
            pageCount: 30,
            sortTable: 'down',
            searchValue: '',
            showCircle: false,
            showAddManagerCircle: false, // 添加管理员弹框的progress
            snackbarStatus: false,
            snackbarMessage: '',
            editAnchorEl: null,
            currentColumn: {}, // 当前点击的行信息
        }
        this.id = 0;
        this.canScroll = true,
        this.currentPage = 0,
        this.currentPageListLength = 0
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
            if(this.currentPageListLength === this.state.pageCount){
                this.currentPage = this.currentPage +  this.state.pageCount ;
                this.getDomainList('scroll')
            }else {
                this.setState({
                    snackbarMessage: '暂无更多数据',
                    snackbarStatus: true,
                })
            }
        }
    }

    getDomainList(flag) {
        this.setState({
            showCircle: true
        })
        const { pageCount, searchValue, tableData } = this.state
        if(flag && flag === 'scroll') {
            
        }else {
            this.currentPage = 0
        }
        let params = {
            count: pageCount, // 每页的数量
            offset: this.currentPage, // 第几页
            domain: searchValue.trim() // 域名模糊值
        }
        getData(domainList, params).then((res) => {
            if(res.data.err_code === 0 && res.data.data !== null) {
                
                let ownData = res.data.data.slice()
                ownData.forEach(item => {
                    let dnsServer = ''
                    item.name_server.forEach((elem, index) => {
                        if(index === item.name_server.length-1) {
                            dnsServer += elem 
                        }else {
                            dnsServer += elem + ','
                        }
                    })
                    item.name_server = dnsServer 
                })
                this.currentPageListLength = ownData.length
                if(flag && flag === 'scroll') { // 判断是否是下拉加载更多的操作
                    let data = tableData.slice()
                    ownData.forEach(item => {
                        data.push(item)
                    })
                    this.setState({
                        tableData: data,
                        showCircle: false
                    })
                   
                }else {
                    this.setState({
                        tableData: ownData,
                        showCircle: false
                    })
                }
            }else {
                this.setState({
                    tableData: [],
                    showCircle: false
                })
            }
        }).catch(err => {
            this.setState({
                tableData: [],
                showCircle: false
            })
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

    getUsers(){
        this.setState({
            showAddManagerCircle: true
        })
    }

    closeAuthorDialog (flag){
        this.setState({
            openAuthorDialog: false,
        })
    }

    searchFn(input) {
        this.setState({
            searchValue: input
        },()=> {
            this.getDomainList()
        })
    }

    closeSnackbar() {
        this.setState({
            snackbarMessage: '',
            snackbarStatus: false
        })
    }

    showMenu(rowData,event) {
        this.setState({
            editAnchorEl: event.currentTarget,
            currentColumn: rowData
        })
    }

    analysisRecord() {
        this.props.history.push('/analysis/'+this.state.currentColumn.domain+'/'+this.state.currentColumn.id)
    }

    authorizInfo() {
        this.setState({
            openAuthorDialog: true,
            editAnchorEl: null,
        },()=>{
            this.dialogRef.setData(this.state.currentColumn)
        })           
    }

    render() {
        const { title, btnText, tableHeader, tableData,editAnchorEl, openAuthorDialog, sortTable, showCircle, showAddManagerCircle, authorDialogTitle } = this.state
        const { classes } = this.props
        return (
            <React.Fragment>
            <AuthorizationInfo open={openAuthorDialog} close={this.closeAuthorDialog.bind(this)} title={authorDialogTitle} showCircle={showAddManagerCircle} innerRef={(element)=>{this.dialogRef = element}}/>
            <div className={styles.dnsRightContainer} ref={(container) => {this.container = container}} onScroll={this.onScroll.bind(this)}>
                <div className={styles.dnsTitle}><DnsTitle title={title}/></div>
                <div className={styles.dnsSearch}>
                    <DnsSearch 
                        btnText={btnText} 
                        managerSearch={true}
                        search={this.searchFn.bind(this)}/>
                </div>
                <div className={styles.dnsTable}> 
                    <Table>
                        <TableHead className={classes.tableHeader}>
                            <TableRow>
                            {
                                tableHeader.map((item, index) => {
                                    if(item.key === 'addTime') {
                                        return (
                                            <TableCell key={index}>
                                                <Tooltip title="排序">
                                                    <TableSortLabel onClick={this.sortTable.bind(this)}>
                                                        {item.label}
                                                       {
                                                            sortTable === 'up' ?
                                                                <ArrowUpward  className={classes.arrow}/> :
                                                                <ArrowDownward className={classes.arrow}/>
                                                        }
                                                    </TableSortLabel>
                                                </Tooltip>
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
                                tableData.length > 0 ? 
                                    (tableData.map((item, index) => {
                                        return (
                                            <TableRow 
                                                key={item.id} 
                                                className={classes.tableRow} 
                                                hover
                                                // selected={currentRow === ''+index}
                                                // onMouseOver={this.onMouseover.bind(this,index)}
                                                // onMouseOut= {this.onMouseout.bind(this)}
                                                >
                                                <TableCell>
                                                    {
                                                        item.domain.length > 40 ?
                                                            (
                                                                <Tooltip title={item.domain}>
                                                                    <span>{item.domain.substr(0,40)+'...'}</span>
                                                                </Tooltip>
                                                            ): <span>{item.domain}</span>
                                                    }
                                                </TableCell>
                                                <TableCell>{item.record_count}</TableCell>
                                                <TableCell>
                                                    {
                                                        item.name_server.length > 40 ?
                                                            (
                                                                <Tooltip title={item.name_server}>
                                                                    <span>{item.name_server.substr(0,40)+'...'}</span>
                                                                </Tooltip>
                                                            ): <span>{item.name_server}</span>
                                                    }
                                                </TableCell>
                                                <TableCell>{item.create_at ? item.create_at.replace(/(T|Z)/ig, '  ') : ''}</TableCell>
                                                <TableCell className={classes.cancleColumn}>
                                                    <Tooltip title="操作" disableFocusListener>
                                                        <Button
                                                            aria-owns={editAnchorEl ? 'template-menu' : null}
                                                            aria-haspopup="true"
                                                            onClick={this.showMenu.bind(this,item)}
                                                            className={styles.dnsTableButton}
                                                            >
                                                            <MoreVert />
                                                        </Button>
                                                    </Tooltip>
                                                </TableCell>
                                            </TableRow>
                                        )
                                    })) : null
                            }
                        </TableBody>
                    </Table>
                    {
                        showCircle ? <CircularProgress className={classes.progress} size={60}/> : null
                    }
                </div>
            </div>
            
            <Menu
                id="template-menu" 
                anchorEl={editAnchorEl}
                open={Boolean(editAnchorEl)}
                onClose={()=>{this.setState({editAnchorEl: null})}}
                className={styles.dnsSelectMenu} >
                <MenuItem onClick={this.analysisRecord.bind(this)} className={classes.menuItem}><Poll />解析记录</MenuItem>
                <MenuItem onClick={this.authorizInfo.bind(this)} className={classes.menuItem}><Assignment />授权信息</MenuItem>
            </Menu>
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

DomainList.propTypes = {
    classes: PropTypes.object.isRequired
}

export default withStyles(muiStyles)(DomainList)