
// 添加直接
function addRegistry(serverId) {
    if(!serverId){
        serverId = 0
    }
    var url = "/image/registry/add"
    var result = get({ServerId:serverId,ClusterName:getClusterName(1)}, url)
    $("#add_groups_html").html(result)
    $("#add_post_html").modal("toggle")
}


/**
 * 删除仓库弹出框
 * 2018-01-20 18:09
 */
function deleteRegistrySwal(id) {
    Swal("删除该仓库", "warning", "确认操作", "不操作", "成功", "失败", " deleteRegistry("+id+")", "loadRegistryData()");
}


/**
 * 重建仓库弹出框
 * 2018-03-02 10:36
 */
function recreateRegistrySwal(id) {
    Swal("重新部署仓库<br>该操作在仓库服务不存在时生效", "warning", "确认操作", "不操作", "成功", "失败", " recreateRegistry("+id+")", "loadRegistryData()");
}

/**
 * 加载数据
 * @param key
 */
function loadRegistryData(key) {
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

    $("#registry-data-table").dataTable({
        "filter": false,//去掉搜索框
        "ordering": false, // 是否允许排序
        "paginationType": "full_numbers", // 页码类型
        "destroy": true,
        "processing": true,
        "bPaginate": false, //是否显示（应用）分页器
        "serverSide": true,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/registry?t=" + new Date().getTime() + "&search=" + key ,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "Entname","sWidth":"7%"},
            {"data": "Name","sWidth":"8%","mRender":function (data,type, full) {
                return "<div style='word-wrap:break-word'><a  target='_self' href='/image/registry/group/list?registryName="+full["ServerDomain"]+"'>"+data+"</a></div>";
            }},
            {"data": "AuthServer","sWidth":"9%", "mRender":function (data) {
                return "<div style='word-wrap:break-word'><a  target='_blank' href="+data+">"+data+"</a></div>";
            }},
            {"data": "ClusterName","sWidth":"7%","mRender":function (data) {
                data = data.replace(/"/g,"")
                return "<a href='/base/cluster/detail/"+data+"'>"+data+"</a>";
            }},
            {"data": "CreateTime","sWidth":"9%"},
            {"data": "Access","sWidth":"20%"},
            {"data": "Status","sWidth":"6%", "mRender":function (data) {
                    if(data=="正常"){
                         return '<div class="Running"><div><i class="fa fa-circle"></i><span>&nbsp;正常</span>'
                    }
                    return '<div class="Fail"><div><i class="fa fa-circle"></i><span>&nbsp;异常</span>&nbsp;'
                }},
            {"data": "ServerId", "sWidth":"6%","mRender": function (data) {
                return '<button type="button" title="更新" onclick="addRegistry(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                '<button type="button"  title="删除" onClick="deleteRegistrySwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>&nbsp;'+
                '<button type="button"  title="重建" onClick="recreateRegistrySwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-undo"></i></button>';

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
function deleteRegistry(id) {
    var url = "/api/registry/"+id;
    var result = del({}, url);
    result = JSON.stringify(result);
    return result
}

/**
 * 2018-03-02 10:38
 * 重新部署仓库服务
 * @param id
 * @return {string}
 */
function recreateRegistry(id) {
    var url = "/api/registry/recreate";
    var result = post({ServerId:id}, url);
    result = JSON.stringify(result);
    return result
}


/**
 * 保存仓库
 */
function saveRegistry(serverId) {
    if(!serverId){
        serverId = 0
    }
    var data = get_form_data();
    data["ServerId"] = parseInt(serverId);
    if(!checkValue(data,"Name,Admin,ClusterName,Entname,Password,AuthServer,ServerDomain")){
        return
    }
    var url = "/api/registry";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadRegistryData();
    }else{
        faild(result);
    }
}
