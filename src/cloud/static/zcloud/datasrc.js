
// 添加直接
function addDataSource(userId) {
    if(!userId){
        userId = 0
    }
    var url = "/log/datasrc/add";
    var result = get({DataSourceId:userId}, url);
    $("#add_groups_html").html(result);
    $("#add_post_html").modal("toggle")
}


/**
 * 删除数据源弹出框
 */
function deleteDataSourceSwal(id) {
    Swal("删除数据源", "warning", "确认操作", "不操作", "成功", "失败", " deleteDataSource("+id+")", "loadDataSourceData()");
}


/**
 * 加载数据
 * @param key
 */
function loadDataSourceData(key) {
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
            "url": "/api/log/datasrc?t=" + new Date().getTime() + "&search=" + key +"&type=datasrc",
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "Ent"},
            {"data": "ClusterName"},
            {"data": "Address"},
            {"data": "Description"},
            {"data": "LastModifyTime"},
            {
                "sWidth": "150px", "data": "DataSourceId", "mRender": function (data) {
                    return '<button type="button" title="更新" onclick="addDataSource(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                        '<button type="button"  title="删除" onClick="deleteDataSourceSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
                }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}



/**
 * 删除数据源方法
 * @param id
 * @return {*}
 */
function deleteDataSource(id) {
    var url = "/api/log/datasrc/"+id;
    var result = del({}, url);
    result = JSON.stringify(result);
    return result
}



/**
 * 2018-03-25 10:41
 * 保存数据源
 */
function saveDataSource(userId) {
    if(!userId){
        userId = 0
    }
    var data = get_form_data();
    data["DataSourceId"] = parseInt(userId);
    if(!checkValue(data,"Ent,ClusterName,Address,Description")){
        return
    }
    var url = "/api/log/datasrc";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadDataSourceData();
    }else{
        faild(result);
    }
}
