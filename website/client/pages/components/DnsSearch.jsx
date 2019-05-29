import React, { Component } from 'react'
import PropTypes from 'prop-types';
import { connect } from 'react-redux'
import styles from '../../styles/sass/common.scss'
import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Search from '@material-ui/icons/Search'
import Select from '@material-ui/core/Select'
import MenuItem from '@material-ui/core/MenuItem'
import { getData } from '../../utils/request';
import { analysisTypes } from '../../api/domainApi'
import { saveSearchValue } from '../../store/actions/homeAction'

const muiStyles = () => ({
    button: {
        color: 'white',
        height: '34px',
        width: "110px",
        fontSize: '14px',
        // padding: '10px 20px',
        boxSizing: 'border-box',
        backgroundColor: '#4A90E2',
        marginRight: '16px',
        '&:hover': {
            backgroundColor: '#4A70e0'
        }
    },
    search: {
        fontSize: '24px'
    },
    input: {
        width: '100%',
        color: '#9b9b9b',
        fontSize: '14px',
        '&:before': {
            borderBottom: 'none',
            borderColor: 'transparent!important'
        },
        '&:before:hover': {
            borderBottom: 'none'
        },
        '&:after': {
            borderBottom: 'none',
            borderColor: 'transparent!important'
        },
        '&:after:hover': {
            borderBottom: 'none',
        },
        '&:hover': {
            border: 'none',
        },
        '&>div>div:focus' :{
            backgroundColor: 'transparent'
        }
    },
    menuItem: {
        color: '#6a6a6a',
        fontSize: '14px',
        '&:hover': {
            backgroundColor: '#ececec'
        },
        '&:focus': {
            backgroundColor: '#e1e1e1'
        }
    }
})

class DnsSearch extends Component {
    constructor(props) {
        super(props)
        this.state = {
            status: '全部',
            recordType: -1,
            options: [],
            inputValue: ''
        }
        this.openFn = this.openFn.bind(this)
    }

    openFn(){
        this.getOptions()
    }

    getOptions() {
        let url = analysisTypes + '/' +this.props.domainId + '/types'
        getData(url)
            .then(res => {
                if(res.data.err_code === 0 && res.data.data !== null) {
                    this.setState({
                        options: res.data.data
                    })
                }
            })
    }    

    openDialog() {
        this.props.openDialog()
    }
    
    search(flag) {
        if(flag === 'keydown') {
            if(event.keyCode === 13) {
                this.props.saveSearchValue(this.state.inputValue)
                this.props.search(this.state.inputValue)
            }
        }else {
            this.props.search(this.state.inputValue)
            this.props.saveSearchValue(this.state.inputValue)
        }
      
    }
    
    selectStatus(event) {
        this.setState({
            status: event.target.value
        },()=>{
            this.props.selectStatus(this.state.status)
        })

    }

    selectRecordType(event) {
        this.setState({
            recordType: event.target.value
        },()=>{
            this.props.selectRecordType(this.state.recordType)
        })

    }

    changeInputValue(event) {
        this.setState({
            inputValue: event.target.value
        })
    }

    componentDidMount() {
        if(this.props.btnText === 'domain') {
            this.setState({
                inputValue: this.props.inputValue
            },()=>{
                this.props.search(this.state.inputValue)
            })
        }
       
    }

    render() {
        const { classes, btnText } = this.props
        const { options, inputValue } = this.state
        return (
            <div className={styles.searchContainer}> 
                {
                    btnText === 'authorization' ?
                        (
                            <Button variant="contained" onClick={this.openDialog.bind(this)} className={classes.button}>
                                <span>新增KS</span>
                            </Button>
                        ) : null

                }
                {
                    btnText === 'analysis' ?
                        (
                            <div className={styles.searchSelect}>
                                <Select
                                    value={this.state.status}
                                    onFocus={this.openFn}

                                    onChange={this.selectStatus.bind(this)}
                                    className={classes.input}
                                    renderValue={value => `记录类型: ${value}`}
                                >
                                    <MenuItem className={classes.menuItem} value="全部">全部</MenuItem>
                                    {
                                        options.map((item, index) => (
                                            <MenuItem key={index} className={classes.menuItem} value={item.record_type}>{item.record_type}</MenuItem>
                                        ))
                                    }
                                </Select>
                            </div>
                        ) : null 
                }
                {
                    btnText === 'authorization' ?
                        (
                            <div className={styles.searchSelect} >
                                <Select
                                    value={this.state.recordType}
                                    onChange={this.selectRecordType.bind(this)}
                                    className={classes.input}
                                    renderValue={value => `状态: ${value===-1?'全部': value===0 ? '启用' : '停用'}`}
                                >
                                    <MenuItem className={classes.menuItem} value={-1}>全部</MenuItem>
                                    <MenuItem className={classes.menuItem} value={0}>启用</MenuItem>
                                    <MenuItem className={classes.menuItem} value={1}>停用</MenuItem>
                                </Select>
                            </div>
                        ) : null 
                }
                 <div className={styles.searchInput}>
                    <input type="text" placeholder={btnText === 'domain' ? "请输入域名" : (btnText === 'analysis' ? "请输入主机记录或记录值" : "请输入K值")} 
                        spellCheck={false}
                        onChange={this.changeInputValue.bind(this)}
                        value={inputValue}
                        onKeyDown={this.search.bind(this, 'keydown')}/> 
                    <div className={styles.searchIcon} onClick={this.search.bind(this, 'click')}><Search className={classes.search}/></div>  
                </div>
            </div>
        )
    }
}  

DnsSearch.propTypes = {
    classes: PropTypes.object.isRequired,
    btnText: PropTypes.string,
    managerSearch: PropTypes.bool,
    openDialog: PropTypes.func,
    search: PropTypes.func, // 父组件传来的输入框搜索的触发的方法
    selectStatus: PropTypes.func,
    selectRecordType: PropTypes.func, 
    addKS: PropTypes.func,
    domainId: PropTypes.string,
    inputValue: PropTypes.string,
    saveSearchValue: PropTypes.func
}

const mapStateToProps = state => {
    return {
        inputValue: state.homeReducer.searchValue
    }
}

const mapDispatchToProps = dispatch => {
    return {
        saveSearchValue: (data)=>{dispatch(saveSearchValue(data))}
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(withStyles(muiStyles)(DnsSearch))