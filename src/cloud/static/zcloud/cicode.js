
// 添加直接
function addCode(repostitoryId) {
    if(!repostitoryId){
        repostitoryId = 0
    }
    var url = "/ci/code/add"
    var result = get({RepostitoryId:repostitoryId}, url)
    $("#add_code_html").html(result)
    $("#add_post_html").modal("toggle")
}


/**
 * 删除用户弹出框
 * 2018-01-20 18:09
 */
function deleteCodeSwal(id) {
    Swal("删除代码仓库", "warning", "确认操作", "不操作", "成功", "失败", " deleteCode("+id+")", "loadCodeData()");
}


/**
 * 加载数据
 * @param key
 */
function loadCodeData(key) {
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

    $("#code-data-table").dataTable({
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
            "url": "/api/ci/code?t=" + new Date().getTime() + "&search=" + key ,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "CodeSource","sWidth":"8%"},
            {"data": "CodeUrl","sWidth":"18%","mRender":function (data) {
                return "<div style='word-wrap:break-word'><a  target='_blank' href="+data+">"+data+"</a></div>";
            }},
            {"data": "Username","sWidth":"10%"},
            {"data": "Type","sWidth":"7%","mRender":function (data) {
                if(data == 1){
                    return  '<div class="Running"><div><i class="fa fa-circle"></i><span>&nbsp;公开</span>';
                }else{
                    return  '<div class="Fail"><div><i class="fa fa-circle"></i><span>&nbsp;非公开</span>';
                }
            }},
            {"data": "GitlabToken","sWidth":"15%"},
            {"data": "LastModifyTime","sWidth":"10%"},
            {"data": "RepostitoryId","sWidth":"7%", "mRender": function (data) {
                    return '<button type="button" title="更新" onclick="addCode(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                        '<button type="button"  title="删除" onClick="deleteCodeSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
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
function deleteCode(id) {
    var url = "/api/ci/code/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}



/**
 * 保存用户
 */
function saveCode(repostitoryId) {
    if(!repostitoryId){
        repostitoryId = 0
    }
    var data = get_form_data();
    data["RepostitoryId"] = parseInt(repostitoryId);
    if(!checkValue(data,"CodeSource,CodeUrl")){
        return
    }
    var url = "/api/ci/code";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadCodeData();
    }else{
        faild(result);
    }
}
