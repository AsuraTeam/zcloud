// 添加权限资源
function addRole(roleId) {
    if(!roleId){
        roleId = 0
    }
    var url = "/system/perm/role/add";
    var result = get({RoleId:roleId}, url);
    $("#add_groups_html").html(result);
    $("#add_post_html").modal("toggle")
}

// 2018-09-10 08:23
// 添加角色权限资源
function addRolePerm(roleId) {
    if(!roleId){
        roleId = 0
    }
    var url = "/system/perm/role/perm/add";
    var result = get({RoleId:roleId}, url);
    $("#add_post_html").modal("toggle")
    $("#add_groups_html").html(result);
}

// 2018-09-11 10:11
// 添加角色用户
function addRolePermUser(roleId) {
    if(!roleId){
        roleId = 0
    }
    var url = "/system/perm/role/user/add";
    var result = get({RoleId:roleId}, url);
    $("#add_groups_html").html(result);
    $("#add_post_html").modal("toggle")
}

/**
 * 删除角色管理弹出框
 * 2018-01-20 18:09
 */
function deleteRoleSwal(id) {
    Swal("删除角色管理", "warning", "确认操作", "不操作", "成功", "失败", " deleteRole("+id+")", "loadRoleData()");
}


/**
 * 加载数据
 * @param key
 */
function loadRoleData(key) {
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
        "ordering": false, // 是否允许7排序
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
            "url": "/api/perm/role?t=" + new Date().getTime() + "&search=" + key ,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "RoleName"},
            {"data": "RoleName"},
            {"data": "Description"},
            {"data": "CreateTime"},
            {
                "sWidth": "150px", "data": "RoleId", "mRender": function (data) {
                    return '<button type="button" title="更新" onclick="addRole(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>' +
                        '<button type="button"  title="分配权限" onClick="addRolePerm(' + data + ')" class="delete-groups m-l-5 btn btn-xs rb-btn-oper"><i class="fa fa-send-o"></i></button>'+
                        '<button type="button"  title="分配角色用户" onClick="addRolePermUser(' + data + ')" class="delete-groups m-l-5 btn btn-xs rb-btn-oper"><i class="fa fa-user-o"></i></button>'+
                        '<button type="button"  title="删除" onClick="deleteRoleSwal(' + data + ')" class="delete-groups m-l-5 btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}



/**
 * 删除角色管理方法
 * @param id
 * @return {*}
 */
function deleteRole(id) {
    var url = "/api/perm/role/"+id;
    var result = del({}, url);
    result = JSON.stringify(result);
    return result
}



/**
 * 保存角色管理
 */
function saveRole(roleId) {
    if(!roleId){
        roleId = 0
    }
    var data = get_form_data();
    data["RoleId"] = parseInt(roleId);
    if(!checkValue(data,"RoleName")){
        return
    }
    var url = "/api/perm/role";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadRoleData();
    }else{
        faild(result);
    }
}
