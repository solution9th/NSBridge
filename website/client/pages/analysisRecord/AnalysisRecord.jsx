import React, { Component } from 'react'
import DnsTitle from '../components/DnsTitle'
import DnsSearch from '../components/DnsSearch'
import PropTypes from 'prop-types'
import styles from '../../styles/sass/common.scss'
import { analysisList } from '../../api/domainApi'
import { getData } from '../../utils/request'
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import CircularProgress from '@material-ui/core/CircularProgress'
import Snackbar from '@material-ui/core/Snackbar'

const muiStyles = () => ({
    tableHeader: {
        color: '#d2d2d2',
        fontSize: '1rem', 
        '&>tr>th:nth-child(1),&>tr>th:nth-child(2),&>tr>th:nth-child(3),&>tr>th:nth-child(4),&>tr>th:nth-child(5)': {
            width: '20%'
        },
    },
    tableBody: {
        color: '#4a4a4a',
        fontSize: '1rem',
        '&>tr>td:nth-child(1),&>tr>td:nth-child(2),&>tr>td:nth-child(3),&>tr>td:nth-child(4),&>tr>td:nth-child(5)': {
            width: '20%'
        },
        
    },
    progress: {
        margin: '10vh 0 0 49%',
        color: '#4A90E2'
    },
})

class AnalysisRecord extends Component {
    constructor(props) {
        super(props)
        this.state = {
            title: [],
            btnText: 'analysis',
            tableHeader: [
                {label: '记录类型', key: 'id'},
                {label: '记录值', key: 'name'},
                {label: '主机记录', key: 'mail'},
                {label: '运营商', key: 'addTime'},
                {label: 'TTL', key: 'updateTime'},
                // {label: '操作', key: ''}
            ],
            tableData: [
                {
                    "id": 10,
                    "domain_id": 6,  
                    "fone_domain_id": 25,
                    "fone_record_id": 45,
                    "sub_domain": "mafeng",  // 主机记录
                    "record_type": "A",  // 记录类型 
                    "value": "mf.sddeznsm.com", // 记录值
                    "line_id": 2,  // 线路
                    "ttl": 3,  // ttl
                    "unit": "hour", // 时间单位
                    "priority": 3,  // 权重
                    "disable": 0,  // 是否禁用
                    "create_at": "2019-03-23T10:24:07Z",
                    "update_at": "2019-03-23T10:24:07Z"
               },
            ],
            // currentRow: '-1',
            pageCount: 30,
            searchValue: '',
            showCircle: false,
            snackbarStatus: false,
            snackbarMessage: '',
            selectValue: '',
        }
        this.id = 0;
        this.canScroll = true,
        this.currentPage = 0,
        this.currentPageListLength = 0
    }

    componentWillMount(){
        
        this.domainId = this.props.match.params.id
    }

    componentDidMount() {
        this.setState({
            title: [
                {label: '域名列表', url: '/'},
                // {label: this.props.match.params.name + '解析记录', url: ''},
                {label: ' 解析记录', url: ''}
            ]
        })

        this.getAnalysisList()
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
                this.getAnalysisList('scroll')
            }else {
                this.setState({
                    snackbarMessage: '暂无更多数据',
                    snackbarStatus: true,
                })
            }
        }
    }

    getAnalysisList(flag) {
        this.setState({
            showCircle: true
        })
        const { pageCount, searchValue, tableData, selectValue } = this.state
        if(flag && flag === 'scroll') {
            
        }else {
            this.currentPage = 0
        }
        let params = {
            sub_or_val: searchValue.trim(), // 主机记录模糊值 记录类型模糊值
            record_type: selectValue === '全部' ? '' : selectValue, // 记录类型
            count: pageCount,
            offset: this.currentPage,
        }
        getData(analysisList+'/'+this.domainId, params).then((res) => {
            if(res.data.err_code === 0 && res.data.data !== null) {
                this.currentPageListLength = res.data.data.length
                res.data.data.forEach(item => {
                    switch(item.unit){
                        case 'sec':
                            item['unitName'] = 's'
                            break;
                        case 'min': 
                            item['unitName'] = 'm'
                            break;
                        case 'hour': 
                            item['unitName'] = 'h'
                            break;
                        case 'day': 
                            item['unitName'] = 'd'
                            break;
                        default: 
                            item['unitName'] = 's'
                    }    
                })
                if(flag && flag === 'scroll') {
                    let data = tableData.slice()
                    res.data.data.forEach(item => {
                        data.push(item)
                    })
                    this.setState({
                        tableData: data,
                        showCircle: false
                    })
                   
                }else {
                    this.setState({
                        tableData: res.data.data,
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

    searchTems(input) {
        this.setState({
            searchValue: input
        },()=> {
            this.getAnalysisList()
        })
    }

    closeSnackbar() {
        this.setState({
            snackbarMessage: '',
            snackbarStatus: false
        })
    }

    selectStatus(selectValue) {
        this.setState({
            selectValue: selectValue
        },()=>{
            this.getAnalysisList()
        })
    }

    render() {
        const { title, btnText, tableHeader, tableData, showCircle } = this.state
        const { classes } = this.props
        return (
            <React.Fragment>
            <div className={styles.dnsRightContainer} ref={(container) => {this.container = container}} onScroll={this.onScroll.bind(this)}>
                <div className={styles.dnsTitle}><DnsTitle title={title} /></div>
                <div className={styles.dnsSearch}>
                    <DnsSearch 
                        btnText={btnText} 
                        managerSearch={true}
                        domainId={this.domainId}
                        selectStatus={this.selectStatus.bind(this)}
                        search={this.searchTems.bind(this)}/>
                </div>
                <div className={styles.dnsTable}> 
                    <Table>
                        <TableHead className={classes.tableHeader}>
                            <TableRow>
                            {
                                tableHeader.map((item, index) => {
                                    return (
                                        <TableCell key={index}>
                                            {item.label}
                                        </TableCell>
                                    ) 
                                })
                            }
                            </TableRow>
                        </TableHead>
                        <TableBody className={classes.tableBody}>
                            {
                                tableData.length > 0 ? 
                                    (tableData.map((item) => {
                                        return (
                                            <TableRow 
                                                key={item.id} 
                                                className={classes.tableRow} 
                                                hover
                                                >
                                                <TableCell>{item.record_type}</TableCell>
                                                <TableCell>{item.value}</TableCell>
                                                <TableCell>{item.sub_domain}</TableCell>
                                                <TableCell>{item.line_name}</TableCell>
                                                <TableCell>{item.ttl}{item.unitName}</TableCell>
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

AnalysisRecord.propTypes = {
    classes: PropTypes.object.isRequired
}

export default withStyles(muiStyles)(AnalysisRecord)