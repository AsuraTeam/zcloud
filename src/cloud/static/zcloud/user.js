
// 添加直接
function addUser(userId) {
    if(!userId){
        userId = 0
    }
    var url = "/system/users/user/add";
    var result = get({UserId:userId}, url);
    $("#add_groups_html").html(result);
    $("#add_post_html").modal("toggle")
}

/**
 * 2018-08-28 09:08
 * 获取用户token
 * @param userId
 */
function showToken(userId) {
    if(!userId){
        userId = 0
    }
    var url = "/system/users/user/token/"+userId;
    var result = get({UserId:userId}, url);
    $("#add_groups_html").html(result);
    $("#add_post_html").modal("toggle")
}


/**
 * 删除用户弹出框
 * 2018-01-20 18:09
 */
function deleteUserSwal(id) {
    Swal("删除用户", "warning", "确认操作", "不操作", "成功", "失败", " deleteUser("+id+")", "loadUserData()");
}


/**
 * 加载数据
 * @param key
 */
function loadUserData(key) {
    if (!key) {
        key = $("#search_user_id").val();
    } else {
        if (key.length < 4) {
            return
        }
    }
    if(!key){
        key = "";
    }

    $("#user-data-table").dataTable({
        "filter": false,//去掉搜索框
        "ordering": false, // 是否允许排序
        "paginationType": "full_numbers", // 页码类型
        "destroy": true,
        "processing": true,
        "bPaginate": true, //是否显示（应用）分页器
        "serverSide": true,
        "bLengthChange": false,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/users?t=" + new Date().getTime() + "&search=" + key ,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "RealName"},
            {"data": "UserName"},
            {"data": "IsDel","mRender":function (data) {
                if(data == 0){
                    return  '<div class="Running"><div><i class="fa fa-circle"></i><span>&nbsp;有效</span>';
                }else{
                    return  '<div class="Fail"><div><i class="fa fa-circle"></i><span>&nbsp;已禁用</span>';
                }
            }},
            {"data": "Description"},
            {"data": "CreateTime"},
            {"data": "LastModifyTime"},
            {
                "sWidth": "150px", "data": "UserId", "mRender": function (data) {
                    return '<button type="button" title="更新" onclick="addUser(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>' +
                        '<button type="button"  title="删除" onClick="deleteUserSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper m-l-5"><i class="fa fa-trash-o"></i></button>'+
                        '<button type="button"  title="显示用户token" onClick="showToken(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper m-l-5"><i class="fa fa-user-secret"></i></button>';
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}



/**
 * 删除用户方法
 * @param id
 * @return {*}
 */
function deleteUser(id) {
    var url = "/api/users/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}



/**
 * 保存用户
 */
function saveUser(userId) {
    if(!userId){
        userId = 0
    }
    var data = get_form_data();
    data["UserId"] = parseInt(userId);
    if(!checkValue(data,"UserName")){
        return
    }
    var url = "/api/users";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadUserData();
    }else{
        faild(result);
    }
}
