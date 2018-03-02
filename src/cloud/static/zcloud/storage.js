// 添加直接
function addStorage(storageId) {
    if (!storageId) {
        storageId = 0
    }
    var url = "/base/storage/add"
    var result = get({StorageId: storageId, ClusterName: getClusterName(1)}, url)
    $("#add_groups_html").html(result)
    $("#add_post_html").modal("toggle")
}


/**
 * 删除存储组弹出框
 * 2018-01-26 18:09
 */
function deleteStorageSwal(id, force) {
    var msg = "删除该存储";
    if (force) {
        msg = "强制删除该存储";
    }
    Swal(msg, "warning", "确认操作", "不操作", "成功", "失败", " deleteStorage(" + id + "," + force + ")", "loadStorageData()");
}


/**
 * 加载数据
 * @param key
 * @param grouptype
 */
function loadStorageData(key, grouptype) {
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
        "processing": true,
        "bPaginate": true, //是否显示（应用）分页器
        "serverSide": true,
        "bLengthChange":false,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/storage?t=" + new Date().getTime() + "&search=" + key + "&groupType=" + grouptype,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "Name", "sWidth": "10%", "mRender": function (data, type, full) {
                return "<a target='_self' href='/base/storage/detail/" + full["StorageId"] + "'>" + data + "</a>";
            }
            },

            {
                "data": "Status", "sWidth": "8%", "mRender": function (data) {
                if (!data) {
                    return '<div title="可以删除" class="Running"><div><i class="fa fa-circle"></i><span>&nbsp;未使用</span>';
                }
                return '<div title="可以删除" class="Fail"><div><i class="fa fa-circle"></i><span>&nbsp;已使用</span>';

            }
            },
            {
                "data": "ClusterName", "sWidth": "10%", "mRender": function (data,type,full) {
                if (data) {
                    data = data.replace(/"/g, "")
                }
                return full["Entname"] + "<br><a class='th-top-8' cl href='/base/cluster/detail/" + data + "'>" + data + "</a>";
            }
            },
            {
                "data": "SharedType", "sWidth": "10%", "mRender": function (data,type,full) {
                if (data == "0") {
                    return "<span class='left10'>共享型:" + "&nbsp;nfs</span>";
                }
                return "<span class='left10'>独享型:&nbsp;" + full["StorageType"] + "</span>";
            }
            },
            {
                "data": "StorageSize", "sWidth": "7%", "mRender": function (data, type, full) {
                if (full["SharedType"] == "0") {
                    return '<div title="共享型无法控制使用空间" class="Fail"><div><i class="fa fa-circle"></i><span>&nbsp;' + data + 'MB</span>';
                }
                return data + "MB";
            }
            },
            {"data": "Description", "sWidth": "13%"},
            {"data": "CreateTime", "sWidth": "9%"},
            {
                "data": "StorageId", "sWidth": "6%", "mRender": function (data) {
                return '<button type="button" title="更新" onclick="addStorage(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                    '<button type="button"  title="删除" onClick="deleteStorageSwal(' + data + ',\'\')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>&nbsp;' +
                    '<button type="button"  title="强制删除" onClick="deleteStorageSwal(' + data + ',1)" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-bolt"></i></button>';
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
function deleteStorage(id, force) {
    if(!force){
        force = "";
    }
    var url = "/api/storage/" + id + "?force=" + force;
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}


/**
 * 2018-01-31 09:19
 * 保存存储卷
 */
function saveStorage(storageId) {
    if (!storageId) {
        storageId = 0
    }
    var data = get_form_data();
    data["StorageId"] = parseInt(storageId);
    if (!checkValue(data, "Name,StorageSize,StorageType,Description,ClusterName")) {
        return
    }
    data["SharedType"] = $("input[name='SharedType']:checked").val();
    data["StorageType"] = $("input[name='StorageType']:checked").val();
    var url = "/api/storage";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        loadStorageData();
    } else {
        faild(result);
    }
}
