
// 添加直接
function addBatch(fileId) {
    if(!fileId){
        fileId = 0
    }
    var url = "/ci/batch/add";
    var result = get({BatchId:fileId}, url);
    $("#add_file_html").html(result);
    $("#add_post_html").modal("toggle")
}


/**
 * 删除batch弹出框
 * 2018-08-25 10:24
 */
function deleteBatchSwal(id,detail) {
    if(detail){
        Swal("删除批量任务", "warning", "确认操作", "不操作", "成功", "失败", " deleteBatch("+id+")", "loadDockerfileList()");
    }else{
        Swal("删除批量任务", "warning", "确认操作", "不操作", "成功", "失败", " deleteBatch("+id+")", "loadBatchData()");
    }
}

/**
 * 将页面刷新到列表页面
 */
function loadDockerfileList() {
    window.location.href = "/ci/batch/list";
}

/**
 * 加载数据
 * @param key
 */
function loadBatchData(key) {
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

    $("#file-data-table").dataTable({
        "filter": false,//去掉搜索框
        "ordering": false, // 是否允许排序
        "paginationType": "full_numbers", // 页码类型
        "destroy": true,
        "processing": true,
        "bLengthChange": false,
        "bPaginate": true, //是否显示（应用）分页器
        "serverSide": true,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/ci/batch?t=" + new Date().getTime() + "&search=" + key ,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "BatchName","sWidth":"8%","mRender":function (data) {
                return data;
            }},
            {"data": "JobIds","sWidth":"17%"},
            {"data": "CreateTime","sWidth":"10%"},
            {"data": "LastModifyUser","sWidth":"15%"},
            {"data": "LastModifyTime","sWidth":"10%"},
            {"data": "Description","sWidth":"15%"},
            {"data": "BatchId","sWidth":"7%", "mRender": function (data) {
                    return '<button type="button" title="更新" onclick="addBatch(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                        '<button type="button"  title="删除" onClick="deleteBatchSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}



/**
 * 删除batch方法
 * @param id
 * @return {*}
 */
function deleteBatch(id) {
    var url = "/api/ci/batch/"+id;
    var result = del({}, url);
    result = JSON.stringify(result);
    return result
}



/**
 * 2018-08-27 11:23
 * 保存batch
 */
function saveBatch(fileId) {
    if(!fileId){
        fileId = 0
    }
    var data = get_form_data();
    data["BatchId"] = parseInt(fileId);
    if(!checkValue(data,"Ent,Cluster,Version")){
        return
    }
    data["BuildType"] = "1";
    if($("#selectStorageType1").is(":checked")){
        data["BuildType"] = "0";
    }
    var url = "/api/ci/batch";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadBatchData();
    }else{
        faild(result);
    }
}
