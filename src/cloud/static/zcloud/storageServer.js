/**
 * 添加存储服务
 * @param storageId
 */
function addStorageServer(storageId) {
    if (!storageId) {
        storageId = 0
    }
    var url = "/base/storage/server/add"
    var result = get({ServerId: storageId, ClusterName: getClusterName(1)}, url)
    $("#add_groups_html").html(result)
    $("#add_post_html").modal("toggle")
}


/**
 * 删除存储组弹出框
 * 2018-02-07 20:45
 * */
function deleteStorageServerSwal(id, force) {
    var msg = "删除该存储服务";
    if (force) {
        msg = "强制删除该存储服务";
    }
    Swal(msg, "warning", "确认操作", "不操作", "成功", "失败", " deleteStorageServer(" + id + "," + force + ")", "loadStorageServerData()");
}


/**
 * 加载数据
 * @param key
 * @param grouptype
 */
function loadStorageServerData(key, grouptype) {
    if (!key) {
        key = $("#search_user_id").val();
    } else {
        if (key.length < 4) {
            return
        }
    }
    key = getValue(key);
    grouptype = getValue(grouptype);

    $("#storage-data-table").dataTable({
        "filter": false,//去掉搜索框
        "ordering": false, // 是否允许排序
        "paginationType": "full_numbers", // 页码类型
        "destroy": true,
        "bLengthChange": false,
        "processing": true,
        "bPaginate": true, //是否显示（应用）分页器
        "serverSide": true,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/storage/server?t=" + new Date().getTime() + "&search=" + key + "&groupType=" + grouptype,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "Entname", "sWidth": "12%", "mRender": function (data, type, full) {
                return data;
            }
            },
            {
                "data": "ClusterName", "sWidth": "12%", "mRender": function (data) {
                if (data) {
                    data = data.replace(/"/g, "")
                }
                return "<a href='/base/cluster/detail/" + data + "'>" + data + "</a>";
            }
            },
            {"data": "StorageType", "sWidth": "12%"},
            {"data": "CreateTime", "sWidth": "12%"},
            {"data": "Description", "sWidth": "22%"},
            {
                "data": "ServerId", "sWidth": "9%", "mRender": function (data) {
                return '<button type="button" title="更新" onclick="addStorageServer(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                    '<button type="button"  title="删除" onClick="deleteStorageServerSwal(' + data + ',\'\')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>&nbsp;' ;
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}


/**
 * 删除存储卷
 * @param id
 * @return {*}
 */
function deleteStorageServer(id, force) {
    if(!force){
        force = "";
    }
    var url = "/api/storage/server/" + id + "?force=" + force;
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}


/**
 * 2018-02-08 09:19
 * 保存存储卷
 */
function saveStorageServer(storageId) {
    if (!storageId) {
        storageId = 0
    }
    var data = get_form_data();
    data["StorageServerId"] = parseInt(storageId);
    if (!checkValue(data, "StorageType,Entname,Description,ClusterName")) {
        return
    }
    data["StorageType"] = $("input[name='StorageType']:checked").val();
    var url = "/api/storage/server";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        loadStorageServerData();
    } else {
        faild(result);
    }
}
