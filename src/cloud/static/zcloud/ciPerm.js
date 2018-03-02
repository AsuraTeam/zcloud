
// 添加权限
function addCiPerm(permissionsId) {
    if(!permissionsId){
        permissionsId = 0
    }
    var url = "/ci/service/perm/add"
    var result = get({PermId:permissionsId}, url)
    $("#add_groups_html").html(result)
    $("#add_post_html").modal("toggle")
}


/**
 * 删除权限
 * 2018-02-18 18:50
 */
function deleteCiPermSwal(id) {
    Swal("删除该权限", "warning", "确认操作", "不操作", "成功", "失败", " deleteCiPerm("+id+")", "loadCiPermData()");
}


/**
 * 加载数据
 * @param key
 */
function loadCiPermData(key) {
    if (!key) {
        key = $("#search_user_id").val();
    } else {
        if (key.length < 4) {
            return
        }
    }
    key = getValue(key);
    $("#ci-perm-data-table").dataTable({
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
            "url": "/api/ci/service/perm?t=" + new Date().getTime() + "&search=" + key,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "Username","sWidth":"11%","mRender":function (data){
            if(!data){
                return ""
            }
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
            if(!data){
                return "";
            }
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
            {"data": "Datas","sWidth":"28%","mRender":function (data) {
                return "<div style='word-wrap:break-word'>" + data + "</div>";
            }},
            {"data": "LastModifyUser","sWidth":"12%"},
            {"data": "LastModifyTime","sWidth":"12%"},
            {"data": "PermId", "sWidth":"5%","mRender": function (data) {
                return '<button type="button" title="更新" onclick="addCiPerm(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                    '<button type="button"  title="删除" onClick="deleteCiPermSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
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
function deleteCiPerm(id) {
    var url = "/api/ci/service/perm/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}



/**
 * 2018-02-18 18:55
 * 保存权限
 */
function saveCiPerm(permissionsId) {
    if(!permissionsId){
        permissionsId = 0
    }
    var data = get_form_data();
    data["PermId"] = parseInt(permissionsId);
    var domains = [];
    $("#domain_redo_to option").each(function () {
        domains.push($(this).val());
    });
    data["Datas"] = domains.join(",");
    var url = "/api/ci/service/perm";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadCiPermData();
    }else{
        faild(result);
    }
}