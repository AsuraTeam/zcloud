
// 添加直接
function addRegistryPerm(permissionsId) {
    if(!permissionsId){
        permissionsId = 0
    }
    var url = "/image/registry/perm/add"
    var result = get({PermissionsId:permissionsId,ClusterName:getClusterName(1)}, url)
    $("#add_groups_html").html(result)
    $("#add_post_html").modal("toggle")
}


/**
 * 删除权限
 * 2018-01-20 21:02
 */
function deleteRegistryPermSwal(id) {
    Swal("删除该权限", "warning", "确认操作", "不操作", "成功", "失败", " deleteRegistryPerm("+id+")", "loadRegistryPermData()");
}


/**
 * 加载数据
 * @param key
 */
function loadRegistryPermData(key) {
    if (!key) {
        key = $("#search_user_id").val();
    } else {
        if (key.length < 4) {
            return
        }
    }
    key = getValue(key);
    var project = getValue($("#project-value").val());
    $("#registry-perm-data-table").dataTable({
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
            "url": "/api/registry/perm?t=" + new Date().getTime() + "&search=" + key +"&project="+project,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "ServiceName","sWidth":"10%"},
            {"data": "Project","sWidth":"12%", "mRender":function (data) {
            if(data=="*"){
                return "<span type='所有项目权限' class='text-danger'>*</span>"
            }
                return data;
            }},
            {"data": "Action","sWidth":"8%","mRender":function (data) {
                return data;
            }},
            {"data": "UserName","sWidth":"11%","mRender":function (data){
                var  v = data.split(",");
                var h = "";
                for(var i=0;i<3;i++){
                    if(v[i]){
                        h += v[i] + "<br>";
                    }
                }
                return "<span title='"+data.replace(/,/g,"\n")+"'>" + h + "</span>";
            }},
            {"data": "GroupsName","sWidth":"12%","mRender":function (data){
                data = data.replace(/"/g,"");
                var  v = data.split(",");
                var h = "";
                for(var i=0;i<5;i++){
                    if(v[i] && v[i] != "null"){
                        h += v[i] + "<br>";
                    }
                }
                return "<span title='"+data.replace(/,/g,"\n")+"'>" + h + "</span>";
            }},
            {"data": "LastModifyTime","sWidth":"17%"},
            {"data": "PermissionsId", "sWidth":"5%","mRender": function (data) {
                return '<button type="button" title="更新" onclick="addRegistryPerm(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                    '<button type="button"  title="删除" onClick="deleteRegistryPermSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
            }
            },

        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}



/**
 * 删除权限方法
 * @param id
 * @return {*}
 */
function deleteRegistryPerm(id) {
    var url = "/api/registry/perm/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}



/**
 * 保存用户
 */
function saveRegistryPerm(permissionsId) {
    if(!permissionsId){
        permissionsId = 0
    }
    var data = get_form_data();
    data["PermissionsId"] = parseInt(permissionsId);
    if(!checkValue(data,"Project,ClusterName")){
        return
    }
    var users = [];
    $("#undo_contact_group_redo_to option").each(function () {
        users.push($(this).val());
    });
    data["UserName"] = users.join(",");

    var groups = [];
    $("#undo_perm_group_redo_to option").each(function () {
        groups.push($(this).val());
    });
    data["GroupsName"] = groups.join(",");


    if(!data["ServiceName"]){
        data["ServiceName"] = $("#select-registry-server-value").val();
        if(!data["ServiceName"]){
            faild("仓库名称必须选择");
            return;
        }
    }

    var action = [];
    if($("input[name='Pull']").is(":checked")){
        action.push("pull");
    }
    if($("input[name='Push']").is(":checked")){
        if($("input[name='Pull']").is(":checked")) {
            action.push("push");
        }else{
            faild("设置 push 权限必须有 pull 权限");
            return;
        }
    }
    if(action.length == 0){
        faild("至少选择一个权限");
        return
    }
    data["Action"] = action.join(",");
    var url = "/api/registry/perm";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadRegistryPermData();
    }else{
        faild(result);
    }
}