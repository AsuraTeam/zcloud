
// 添加环境
function addEnt(entId) {
    if(!entId){
        entId = 0;
    }
    var url = "/system/ent/add";
    var result = get({EntId:entId}, url);
    $("#add_ent_html").html(result);
    $("#add_post_html").modal("toggle")
}


/**
 * 删除环境弹出框
 * 2018-01-20 18:09
 */
function deleteEntSwal(id) {
    Swal("删除环境", "warning", "确认操作", "不操作", "成功", "失败", " deleteEnt("+id+")", "loadEntData()");
}


/**
 * 加载数据
 * @param key
 */
function loadEntData(key) {
    if (!key) {
        key = $("#search_ent_id").val();
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
            "url": "/api/ent?t=" + new Date().getTime() + "&search=" + key ,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "Entname","mRender":function (data) {
                return data;
            }},
            {"data": "Clusters", "mRender":function (data) {
                return data.split(",").join("<br>");
            }},
            {"data": "CreateTime"},
            {"data": "LastModifyTime"},
            {"data": "Description"},
            {
                "sWidth": "150px", "data": "EntId", "mRender": function (data) {
                    return '<button type="button" title="更新" onclick="addEnt(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                        '<button type="button"  title="删除" onClick="deleteEntSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}



/**
 * 删除环境方法
 * @param id
 * @return {*}
 */
function deleteEnt(id) {
    var url = "/api/ent/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}



/**
 * 保存环境
 */
function saveEnt(entId) {
    if(!entId){
        entId = 0
    }
    var data = get_form_data();
    data["EntId"] = parseInt(entId);
    if(!checkValue(data,"Entname")){
        return
    }
    var clusters = [];
    $("#undo_entname_redo_to option").each(function () {
        clusters.push($(this).val());
    });
    data["Clusters"] = clusters.join(",");
    if(!checkValue(data,"Clusters")){
        return
    }
    var url = "/api/ent";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadEntData();
    }else{
        faild(result);
    }
}
