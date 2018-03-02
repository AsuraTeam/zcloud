
// 添加基础镜像
function addBase(baseId) {
    if(!baseId){
        baseId = 0
    }
    var url = "/image/registry/base/add"
    var result = get({BaseId:baseId,ClusterName:getClusterName(1)}, url)
    $("#add_base_html").html(result)
    $("#add_post_html").modal("toggle")
}


/**
 * 删除基础镜像
 * 2018-02-09 16:32
 */
function deleteBaseSwal(id) {
    Swal("删除该基础镜像", "warning", "确认操作", "不操作", "成功", "失败", " deleteBase("+id+")", "loadBaseData()");
}


/**
 * 加载数据
 * @param key
 */
function loadBaseData(key) {
    if (!key) {
        key = $("#search_user_id").val();
    } else {
        if (key.length < 4) {
            return
        }
    }
    key = getValue(key);
    var project = getValue($("#project-value").val());
    $("#registry-base-data-table").dataTable({
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
            "url": "/api/registry/base?t=" + new Date().getTime() + "&search=" + key +"&project="+project,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "RegistryServer","sWidth":"10%"},
            {"data": "ImageName","sWidth":"10%"},
            {"data": "ImageType","sWidth":"10%"},
            {"data": "CreateTime","sWidth":"10%"},
            {"data": "LastModifyTime","sWidth":"15%"},
            {"data": "Description","sWidth":"15%"},
            {"data": "BaseId", "sWidth":"5%","mRender": function (data) {
                return '<button type="button" title="更新" onclick="addBase(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                    '<button type="button"  title="删除" onClick="deleteBaseSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
            }
            },

        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}



/**
 * 删除基础镜像
 * @param id
 * @return {*}
 */
function deleteBase(id) {
    var url = "/api/registry/base/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}



/**
 * 保存用户
 */
function saveBase(baseId) {
    if(!baseId){
        baseId = 0
    }
    var data = get_form_data();
    data["BaseId"] = parseInt(baseId);
    if(!checkValue(data,"Description,ImageName")){
        return
    }

    if(!data["RegistryServer"]){
        data["RegistryServer"] = $("#select-registry-server-value").val();
        if(!data["RegistryServer"]){
            faild("仓库名称必须选择");
            return;
        }
    }
    
    var url = "/api/registry/base";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadBaseData();
    }else{
        faild(result);
    }
}