
// 添加直接
function addGroup(groupsId) {
    if(!groupsId){
        groupsId = 0
    }
    var url = "/system/users/groups/add"
    var result = get({GroupsId:groupsId}, url)
    $("#add_groups_html").html(result)
    $("#add_post_html").modal("toggle")
}


/**
 * 删除部门团队弹出框
 */
function deleteGroupSwal(id) {
    Swal("删除部门", "warning", "确认操作", "不操作", "成功", "失败", " deleteGroups("+id+")", "loadGroupData()");
}


/**
 * 加载数据
 * @param key
 */
function loadGroupData(key) {
    if (!key) {
        key = $("#search_groups_data").val();
    } else {
        if (key.length < 2) {
            return
        }
    }
    if(!key){
        key = "";
    }

    $("#groups-data-table").dataTable({
        "filter": false,//去掉搜索框
        "ordering": false, // 是否允许排序
        "paginationType": "full_numbers", // 页码类型
        "destroy": true,
        "processing": true,
        "bPaginate": false, //是否显示（应用）分页器
        "serverSide": true,
        "bLengthChange": false,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/groups?t=" + new Date().getTime() + "&search=" + key ,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "GroupsName"},
            {"data": "Users"},
            {"data": "Description"},
            {"data": "CreateTime"},
            {"data": "LastModifyTime"},
            {
                "sWidth": "150px", "data": "GroupsId", "mRender": function (data, type, full) {
                    return '<button type="button" title="更新" onclick="addGroup(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                        '<button type="button"  title="删除" onClick="deleteGroupSwal('+data+')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}



/**
 * 删除部门团队方法
 * @param id
 * @return {*}
 */
function deleteGroups(id) {
    var url = "/api/groups/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}



/**
 * 保存部门团队
 */
function saveGroups(groupsId) {
    if(!groupsId){
        groupsId = 0
    }
    var data = get_form_data();
    data["GroupsId"] = parseInt(groupsId);
    var users = [];
    $("#undo_contact_group_redo_to option").each(function () {
        users.push($(this).val());
    });
    data["Users"] = users.join(",");
    if(!checkValue(data,"GroupsName")){
        return
    }
    var url = "/api/groups";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadGroupData();
    }else{
        faild(result);
    }
}
