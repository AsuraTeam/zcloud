// 添加权限资源
function addResource(userId) {
    if(!userId){
        userId = 0
    }
    var url = "/system/perm/resource/add"
    var result = get({ResourceId:userId}, url)
    $("#add_groups_html").html(result)
    $("#add_post_html").modal("toggle")
}


/**
 * 删除资源管理弹出框
 * 2018-01-20 18:09
 */
function deleteResourceSwal(id) {
    Swal("删除资源管理", "warning", "确认操作", "不操作", "成功", "失败", " deleteResource("+id+")", "loadResourceData()");
}


/**
 * 加载数据
 * @param key
 */
function loadResourceData(key) {
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
        "displayLength": 30, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/perm/resource?t=" + new Date().getTime() + "&search=" + key ,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "Name"},
            {"data": "ApiUrl"},
            {"data": "Method"},
            {"data": "ApiType"},
            {"data": "Parent"},
            {"data": "CreateTime"},
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}



/**
 * 删除资源管理方法
 * @param id
 * @return {*}
 */
function deleteResource(id) {
    var url = "/api/perm/resource/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}



/**
 * 保存资源管理
 */
function saveResource(userId) {
    if(!userId){
        userId = 0
    }
    var data = get_form_data();
    data["ResourceId"] = parseInt(userId);
    if(!checkValue(data,"Name")){
        return
    }
    var url = "/api/perm/resource";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadResourceData();
    }else{
        faild(result);
    }
}
