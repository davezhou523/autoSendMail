package xerr

// 成功返回
const OK int64 = 200

/**(前3位代表业务,后三位代表具体功能)**/

// 全局错误码
const Fail int64 = 1000
const SERVER_COMMON_ERROR int64 = 1001
const REUQEST_PARAM_ERROR int64 = 1002
const TOKEN_EXPIRE_ERROR int64 = 1003
const TOKEN_GENERATE_ERROR int64 = 1004
const DB_ERROR int64 = 1005
const DB_UPDATE_AFFECTED_ZERO_ERROR int64 = 1006

//用户模块
