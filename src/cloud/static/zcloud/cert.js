
// 添加直接
function addCert(fileId) {
    if(!fileId){
        fileId = 0
    }
    var url = "/base/network/cert/add"
    var result = get({CertId:fileId}, url)
    $("#add_file_html").html(result)
    $("#add_post_html").modal("toggle")
}


/**
 * 删除证书弹出框
 * 2018-02-02 16:46
 */
function deleteCertSwal(id) {
    Swal("删除该证书", "warning", "确认操作", "不操作", "成功", "失败", " deleteCert("+id+")", "loadCertData()");
}

/**
 * 加载数据
 * @param key
 */
function loadCertData(key) {
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

    $("#file-data-table").dataTable({
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
            "url": "/api/network/cert?t=" + new Date().getTime() + "&search=" + key ,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "CertKey","sWidth":"15%","mRender":function (data) {
                return data;
            }},
            {"data": "Description","sWidth":"20%"},
            {"data": "LastModifyUser","sWidth":"15%"},
            {"data": "LastModifyTime","sWidth":"20%"},
            {"data": "CertId","sWidth":"12%", "mRender": function (data) {
                    return '<button type="button" title="更新" onclick="addCert(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                            '<button type="button"  title="删除" onClick="deleteCertSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
               }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}



/**
 * 删除证书方法
 * 2018-02-02 17;00
 * @param id
 * @return {*}
 */
function deleteCert(id) {
    var url = "/api/network/cert/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}



/**
 * 保存证书
 * 2018-02-02 17:05
 */
function saveCert(fileId) {
    if(!fileId){
        fileId = 0
    }
    var data = get_form_data();
    data["CertId"] = parseInt(fileId);
    if(!checkValue(data,"CertKey,CertValue,PemValue")){
        return
    }
    var url = "/api/network/cert";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadCertData();
    }else{
        faild(result);
    }
}
