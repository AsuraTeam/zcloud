// 添加直接
function addRegistryGroup(serverId) {
    if (!serverId) {
        serverId = 0
    }
    var url = "/image/registry/group/add"
    var result = get({GroupId: serverId, ClusterName: getClusterName(1)}, url)
    $("#add_groups_html").html(result)
    $("#add_post_html").modal("toggle")
}

/**
 * 删除镜像弹出框
 * 2018-01-29 9:11
 */
function deleteRegistryGroupImageSwal(id) {
    Swal("删除该镜像", "warning", "确认操作", "不操作", "成功", "失败", " deleteRegistryImage(" + id + ")", "loadRegistryImageData()");
}

/**
 * 删除镜像
 * @param id
 * @return {string}
 */
function deleteRegistryImage(id) {
    var url = "/api/registry/group/images/"+id;
    var result = del({}, url);
    return JSON.stringify(result);
}

/**
 * 查看镜像tag
 * 2018-01-29 10:55
 * @param id
 */
function showTags(id) {
    var url = "/api/registry/group/images/"+id;
    var result = get({}, url);
    var data = result["data"];
    $("#show_image_tag_html_content").html(data["Tags"].replace(/,/g, "\n"));
    $("#tags_title").html(data["Name"]);
    $("#show_image_tag_html").modal("toggle");
}

/**
 * 删除仓库组弹出框
 * 2018-01-26 18:09
 */
function deleteRegistryGroupSwal(id) {
    Swal("删除该仓库", "warning", "确认操作", "不操作", "成功", "失败", " deleteRegistryGroup(" + id + ")", "loadRegistryGroupData()");
}


/**
 * 加载数据
 * @param key
 * @param grouptype
 */
function loadRegistryGroupData(key, grouptype) {
    if (!key) {
        key = $("#search_user_id").val();
    } else {
        if (key.length < 4) {
            return
        }
    }
    key = getValue(key);
    grouptype = getValue(grouptype);

    $("#registry-data-table").dataTable({
        "filter": false,//去掉搜索框
        "ordering": false, // 是否允许排序
        "paginationType": "full_numbers", // 页码类型
        "destroy": true,
        "processing": true,
        "bPaginate": true, //是否显示（应用）分页器
        "serverSide": true,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/registry/group?t=" + new Date().getTime() + "&search=" + key + "&groupType=" + grouptype,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "GroupName", "sWidth": "10%", "mRender": function (data,type,full) {
                return "<a target='_self' href='/image/registry/group/detail/" + full["GroupId"] + "'>" + data + "</a>";
            }
            },
            {
                "data": "ClusterName", "sWidth": "10%", "mRender": function (data) {
                if (data) {
                    data = data.replace(/"/g, "")
                }
                return "<a href='/base/cluster/detail/" + data + "'>" + data + "</a>";
            }
            },
            {"data": "ServerDomain", "sWidth": "8%"},
            {"data": "ImageNumber", "sWidth": "7%"},
            {"data": "TagNumber", "sWidth": "7%"},
            {"data": "SizeTotle", "sWidth": "6%"},
            {"data": "GroupType", "sWidth": "7%"},
            {"data": "LastModifyTime", "sWidth": "8%"},
            {
                "data": "GroupId", "sWidth": "6%", "mRender": function (data) {
                return '<button type="button" title="更新" onclick="addRegistryGroup(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                    '<button type="button"  title="删除" onClick="deleteRegistryGroupSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}


/**
 * 删除仓库组方法
 * @param id
 * @return {*}
 */
function deleteRegistryGroup(id) {
    var url = "/api/registry/group/" + id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}


/**
 * 保存仓库组
 */
function saveRegistryGroup(serverId) {
    if (!serverId) {
        serverId = 0
    }
    var data = get_form_data();
    data["GroupId"] = parseInt(serverId);
    if (!checkValue(data, "GroupName,ClusterName,ServerDomain,GroupType")) {
        return
    }
    if(data["ServerDomain"]=="--请选择--"){
        return
    }
    var url = "/api/registry/group";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        // loadRegistryGroupData();
    } else {
        faild(result);
    }
}
