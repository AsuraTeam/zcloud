/**
 * 2018-02-06 20:33
 * 添加镜像同步
 * @param syncId
 */
function addSync(syncId,copy) {
    if (!syncId) {
        syncId = 0;
    }
    copy = getValue(copy);
    var url = "/image/sync/add";
    var result = get({SyncId: syncId,copy:copy}, url);
    $("#add_sync_html").html(result);
    $("#add_post_html").modal("toggle")
}

/**
 * 同意申请弹出框
 * 2018-02-06 21:24
 */
function ApprovedSyncSwal(id) {
    Swal("将同意该申请", "warning", "确认操作", "不操作", "成功", "失败", " ApprovedSync("+id+")", "loadRegistryData()");
}

/**
 * 同意镜像同步申请
 * @param id
 * @return {*}
 * @constructor
 */
function ApprovedSync(id) {
    var url = "/api/image/sync/approved/"+id
    var result = post({}, url)
    result = JSON.stringify(result)
    return result
}

/**
 * 删除镜像同步弹出框
 * 2018-01-20 18:09
 */
function deleteSyncSwal(id) {
    Swal("删除该同步申请", "warning", "确认操作", "不操作", "成功", "失败", " deleteSync("+id+")", "loadRegistryData()");
}

/**
 * 2018-02-06 20:40
 * 删除镜像同步
 * @param id
 * @return {*}
 */
function deleteSync(id) {
    var url = "/api/image/sync/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}

/**
 * 执行同步
 * 2018-02-06 21:48
 * @param id
 * @return {*}
 */
function imageSync(id) {
    var url = "/api/image/sync/"+id
    var result = get({}, url)
    result = JSON.stringify(result)
    success(result)
    return result
}
/**
 * 2018-02-06 20:36
 * 加载数据
 * @param key
 */
function loadSyncData(key) {
    if (!key) {
        key = $("#search_user_id").val();
    } else {
        if (key.length < 4) {
            return
        }
    }
    if (!key) {
        key = "";
    }

    $("#history-data-table").dataTable({
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
            "url": "/api/image/sync?t=" + new Date().getTime() + "&search=" + key,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "Entname", "sWidth": "9%", "mRender": function (data,type,full) {
                return "<span style='top: 5px;    position: relative;'>源站:&nbsp;</span><span class='RunningTop5'>" + data + "</span><br><br><span style='    position: relative;top: 5px;'>目标:</span>&nbsp;<span class='FailTop5'>" + full["TargetEntname"]+"</span>";
            }
            },
            {
                "data": "Registry", "sWidth": "10%", "mRender": function (data, type, full) {
                return "<span>" + data + "</span><br><br>"+"/" + full["ImageName"];
            }
            },
            {"data": "Version", "sWidth": "9%"},
            {"data": "TargetRegistry", "sWidth": "10%"},
            {
                "data": "CreateUser", "sWidth": "10%", "mRender": function (data, type, full) {
                return "<span>" + data + "</span><br><br>" + full["CreateTime"];
            }
            },
            {
                "data": "ApprovedBy", "sWidth": "10%", "mRender": function (data, type, full) {
                    if(!data){
                        var r = '<button type="button" title="点击同意该申请" onclick="ApprovedSyncSwal(' + full["SyncId"] + ')" class="btn btn-xs rb-btn-oper"><i class="fa  fa-check-square-o"></i>&nbsp;点击同意</button>&nbsp;<br><span class="FailTop5">请审批</span>'
                        return r;
                    }else {
                        return "<span>" + data + "</span><br>" + full["ApprovedTime"]+"<br><span class='RunningTop5'>已审批</span>";
                    }
            }
            },
            {"data": "Description", "sWidth": "8%","mRender":function (data) {
                return "<div style='word-wrap:break-word'><a>"+data+"</a></div>";
            }},
            {
                "data": "SyncId", "sWidth": "6%", "mRender": function (data, type, full) {
                if (full["ApprovedBy"]) {
                    var r = '<button type="button" title="更新" onclick="addSync(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;';

                } else {
                    var r = '<button type="button" title="更新" onclick="addSync(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;';
                }
                r += '<button type="button"  title="复制一个" onClick="addSync(' + data + ',1)" class="delete-groups btn btn-xs rb-btn-oper"><i class="mdi mdi-copyright"></i></button>&nbsp;';
                if (full["ApprovedBy"]) {
                    r += '<button type="button"  title="执行镜像同步" onClick="imageSync(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="mdi mdi-sync"></i></button><br><span class="RunningTop5">可同步</span>';
                }else{
                    r += '<button type="button"  title="删除" onClick="deleteSyncSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button><br><span class="FailTop5">等待审批</span>';
                }
                return r;
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });

}/**
 * Created by zhaoyun on 2018/2/6.
 */
