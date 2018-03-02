// 添加自动伸缩
function addAutoScale(scaleId) {
    if (!scaleId) {
        scaleId = 0
    }
    var url = "/monitor/scale/add";
    var result = get({ScaleId: scaleId}, url);
    $("#add_code_html").html(result);
    $("#add_post_html").modal("toggle")
}


/**
 * 删除自动伸缩弹出框
 * 2018-02-19 19:09
 */
function deleteAutoScaleSwal(id) {
    Swal("删除自动伸缩", "warning", "确认操作", "不操作", "成功", "失败", " deleteAutoScale(" + id + ")", "loadAutoScale()");
}

/**
 * 删除自动伸缩方法
 * @param id
 * @return {*}
 */
function deleteAutoScale(id) {
    var url = "/api/monitor/scale/" + id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}

/**
 * 2018-02-19 19:05
 * 保存自动伸缩
 */
function saveAutoScale(scaleId) {
    if (!scaleId) {
        scaleId = 0
    }
    var data = get_form_data();
    data["ScaleId"] = parseInt(scaleId);
    if (!checkValue(data, "Entname,ClusterName,ServiceName,AppName,Description,Step,ActionInterval")) {
        return
    }

    if($("#radio10").is(":checked")){
        data["MetricType"] = "custom"
    }
    if($("#radio20").is(":checked")){
        data["MetricType"] = "system"
    }
    data["MetricName"] = data["MetricName"].replace(/--请选择--/,"");
    if(data["MetricType"] =="custom"){
        if (!checkValue(data, "Query")) {
            return
        }
    }else {
        if (!checkValue(data, "MetricName")){
            return
        }
    }
    if(data["DataSource"] == "es"){
        if (!checkValue(data, "Es")){
            return
        }
    }
    var url = "/api/monitor/scale";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        loadAutoScale();
    } else {
        faild(result);
    }
}

/**
 * 2018-02-19 19:07
 * 加载数据
 * @param key
 */
function loadAutoScale(key) {
    if (!key) {
        key = ""
    } else {
        if (key.length < 4) {
            return
        }
    }

    $("#scale-data-table").dataTable({
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
            "url": "/api/monitor/scale?t=" + new Date().getTime() + "&search=" + key ,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "Entname", "sWidth": "8%", "mRender": function (data, type, full) {
                return data + "<br><a class='th-top-8'>" + full["ClusterName"] + "</a>";
            }
            },
            {
                "data": "AppName", "sWidth": "7%", "mRender": function (data, type, full) {
                return data + "<br><a  class='th-top-8' href='/application/service/list?name=" + full["ServiceName"] + "'>" + full["ServiceName"] + "</a>";
            }
            },
            {
                "data": "Gt", "sWidth": "8%", "mRender": function (data) {
                return "大于:&nbsp;" + data ;
            }
            },
            {"data": "ActionInterval","sWidth": "6%"},
            {
                "data": "Step", "sWidth": "6%", "mRender": function (data, type, full) {
                return "查询步长:<span class='m-l-10'>" +data + "</span><br><span class='th-top-8'>数值范围:</span><span class='m-l-10 th-top-8'>" + full["LastCount"] + "</span>";
            }
            },
            {
                "data": "IncreaseStep", "sWidth": "8%", "mRender": function (data) {
                return "扩容步长:<span class='m-l-10'>" +data + "</span>";
            }
            },
            {"data": "CreateTime","sWidth": "8%"},
            {"data": "Description","sWidth": "12%" },
            {
                "sWidth": "6%", "data": "ScaleId", "mRender": function (data) {

                return '<button type="button" title="更新" onclick="addAutoScale(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>' +
                    '<button type="button"  title="删除" onClick="deleteAutoScaleSwal(' + data + ')" class="delete-configure m-l-5 btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>'+
                    '<button type="button"  title="扩缩容日志" onClick="toScaleHistory(' + data + ')" class="delete-configure m-l-5 btn btn-xs rb-btn-oper"><i class="fa fa-history"></i></button>';

                return ""
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}

/**
 * 2018-02-20 07:50
 * 跳转到操作页面
 * @param id
 */
function toScaleHistory(id) {
    var url = "/monitor/scale/logs/?ServiceId=" + id;
    window.location.href = url;
}


/**
 * 2018-02-19 19:05
 * 保存自动伸缩
 */
function saveAutoScale(scaleId) {
    if (!scaleId) {
        scaleId = 0
    }
    var data = get_form_data();
    data["ScaleId"] = parseInt(scaleId);
    if (!checkValue(data, "Entname,ClusterName,ServiceName,AppName,Description,Step,ActionInterval")) {
        return
    }

    if($("#radio10").is(":checked")){
        data["MetricType"] = "custom"
    }
    if($("#radio20").is(":checked")){
        data["MetricType"] = "system"
    }
    data["MetricName"] = data["MetricName"].replace(/--请选择--/,"");
    if(data["MetricType"] =="custom"){
        if (!checkValue(data, "Query")) {
            return
        }
    }else {
        if (!checkValue(data, "MetricName")){
            return
        }
    }
    if(data["DataSource"] == "es"){
        if (!checkValue(data, "Es")){
            return
        }
    }
    var url = "/api/monitor/scale";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        loadAutoScale();
    } else {
        faild(result);
    }
}

/**
 * 2018-02-20 17:33
 * 加载日志数据
 * @param key
 */
function loadAutoScaleLog(key) {
    if (!key) {
        key = ""
    } else {
        if (key.length < 4) {
            return
        }
    }

    $("#scale-data-log-table").dataTable({
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
            "url": "/api/monitor/scale/logs?t=" + new Date().getTime() + "&search=" + key ,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "Entname", "sWidth": "8%", "mRender": function (data, type, full) {
                return data + "<br><a class='th-top-8'>" + full["ClusterName"] + "</a>";
            }
            },
            {
                "data": "AppName", "sWidth": "7%", "mRender": function (data, type, full) {
                return data + "<br><a  class='th-top-8' href='/application/service/list?name=" + full["ServiceName"] + "'>" + full["ServiceName"] + "</a>";
            }
            },
            {
                "data": "Gt", "sWidth": "8%", "mRender": function (data, type, full) {
                return "大于:&nbsp;" + data;
            }
            },
            {"data": "ActionInterval","sWidth": "6%"},
            {
                "data": "Step", "sWidth": "6%", "mRender": function (data, type, full) {
                return "查询步长:<span class='m-l-10'>" +data + "</span><br><span class='th-top-8'>数值范围:</span><span class='m-l-10 th-top-8'>" + full["LastCount"] + "</span>";
            }
            },
            {
                "data": "IncreaseStep", "sWidth": "8%", "mRender": function (data, type, full) {
                return "扩容步长:<span class='m-l-10'>" +data + "</span>";
            }
            },
            {"data": "CreateTime","sWidth": "8%"},
            {"data": "Replicas","sWidth": "6%" },
            {"data": "MonitorValue","sWidth": "6%" , "mRender": function (data, type, full) {
            return "触发值:" + data + "<br><span class='th-top-8'>"+full["Status"]+"</span>";
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}