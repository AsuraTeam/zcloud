// 添加直接
function addPipeline(pipelineId) {
    if (!pipelineId) {
        pipelineId = 0;
    }
    var url = "/pipeline/add";
    var result = get({PipelineId: pipelineId}, url);
    $("#add_pipeline_html").html(result);
    $("#add_post_html").modal("toggle")
}

/**
 * 执行流水线程序
 * @param pipelineId
 */
function execPipeline(pipelineId) {
    if (!pipelineId) {
        return;
    }
    var url = "/api/pipeline/exec/" + pipelineId;
    var result = get({}, url);
    result = JSON.stringify(result);
    setTimeout(function () {
        loadPipelineData();
    }, 3000);
    setTimeout(function () {
        loadPipelineData();
    }, 5000);
    return result;
}


/**
 * 删除流水线弹出框
 * 2018-02-05 10:24
 */
function deletePipelineSwal(id, detail) {
    Swal("删除流水线", "warning", "确认操作", "不操作", "成功", "失败", " deletePipeline(" + id + ")", "loadPipelineData()");

}

/**
 * 执行流水线弹出框
 * 2018-02-07 5:44
 * @param id
 */
function startExecPipeline(id) {
    Swal("是否执行该流水线", "warning", "确认操作", "不操作", "成功", "失败", " execPipeline(" + id + ")", "loadPipelineData()");
}

/**
 * 加载数据
 * @param key
 */
function loadPipelineData(key) {
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

    $("#file-data-table").dataTable({
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
            "url": "/api/pipeline?t=" + new Date().getTime() + "&search=" + key,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "PipelineName", "sWidth": "8%", "mRender": function (data) {
                return "<a href='/pipeline/detail/" + data + "'>" + data + "</a>"
            }
            },
            {"data": "ClusterName", "sWidth": "9%"},
            {
                "data": "AppName", "sWidth": "9%", "mRender": function (data,type,full) {
                    if(full["Status"] == "false"){
                        return "<span class='Fail'>"+data+"</span>"
                    }
                return "<span class='Running'>"+data+"</span>"
            }
            },
            {
                "data": "ServiceName", "sWidth": "9%", "mRender": function (data,type,full) {
                if(full["Status"] == "false"){
                    return "<span class='Fail'>"+data+"</span><br><span class='text-default FailTop5'>服务不存在</span>"
                }
                return "<span class='Running'>"+data+"</span>"
            }
            },
            {"data": "JobName", "sWidth": "9%"},
            {"data": "Description", "sWidth": "10%"},
            {"data": "LastModifyTime", "sWidth": "9%"},
            {
                "data": "PipelineId", "sWidth": "9%", "mRender": function (data, type, full) {
                return '<button type="button" title="更新" onclick="addPipeline(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                    '<button type="button" title="立刻构建" onclick="startExecPipeline(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-transgender-alt"></i></button>&nbsp;' +
                    '<button type="button" title="最近流水线日志" onclick="jobLog(' + full["JobId"] + ')" class="btn btn-xs rb-btn-oper"><i class="fa  fa-hospital-o"></i></button>&nbsp;' +
                    '<button type="button"  title="删除" onClick="deletePipelineSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>&nbsp;' +
                    '<button type="button"  title="流水线历史" onClick="buildHistory(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-history"></i></button>';
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}

/**
 * 获取构建执行日志
 * @param jobId
 */
function jobLog(jobId, history) {
    history = getValue(history);
    var url = "/ci/job/logs/" + jobId;
    var result = get({history: history}, url);
    $("#add_job_html").html(result);
    $("#add_post_html").modal("toggle")
}

/**
 * 2018-02-05 16;38
 * 通过jobname获取日志信息
 * @param jobName
 */
function jobnameLog(jobId, jobName) {
    var url = "/ci/job/logs/" + jobId+"?jobName="+jobName;
    var result = get({}, url);
    $("#add_job_html").html(result);
    $("#add_post_html").modal("toggle")
}


/**
 * 删除流水线方法
 * @param id
 * @return {*}
 */
function deletePipeline(id) {
    var url = "/api/pipeline/" + id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}


/**
 * 保存流水线
 */
function savePipeline(pipelineId) {
    if (!pipelineId) {
        pipelineId = 0
    }
    var data = get_form_data();

    data["PipelineId"] = parseInt(pipelineId);
    if (!checkValue(data, "JobName,ServiceName,AppName,PipelineName,ClusterName,ExecTime")) {
        return
    }
    if ($("#setSelectFailAction1").is(":checked")) {
        data["FailAction"] = "continue"
    } else {
        data["FailAction"] = "pause"
    }

    console.log(data);
    var url = "/api/pipeline";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        loadPipelineData();
    } else {
        faild(result);
    }
}


/**
 * 流水线历史
 * @param id
 */
function buildHistory(id) {
    window.location.href = "/pipeline/history/list";
}


/**
 * 加载数据
 * @param key
 */
function loadPipelineHistoryData(key) {
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
        "processing": true,
        "bPaginate": true, //是否显示（应用）分页器
        "serverSide": true,
        "bLengthChange": false,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/pipeline/history?t=" + new Date().getTime() + "&search=" + key,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "PipelineName", "sWidth": "8%", "mRender": function (data, type, full) {
                    var jobname = full["JobName"];
                return "<a href='/pipeline/detail/" + jobname+ "?JobName="+jobname+"'>" + data + "</a><br>" + full["ClusterName"]
            }
            },
            {
                "data": "ServiceName", "sWidth": "9%", "mRender": function (data, type, full) {
                return "<a href='/application/app/list'>" + full["AppName"] + "</a><br><a style='color:#ffa91c ;' href='/application/service/list'>" + data+"</a>";
            }
            },
            {
                "data": "JobName", "sWidth": "17%", "mRender": function (data) {
                return "<div style='word-wrap:break-word'>" + data + "</div>";
            }
            },
            {
                "data": "Status", "sWidth": "7%", "mRender": function (data) {
                if (data == "执行失败") {
                    return "<span class='Fail'>" + data + "</span>"
                }
                return "<span class='Running'>" + data + "</span>"
            }
            },
            {
                "data": "RunTime", "sWidth": "8%", "mRender": function (data) {
                return "<span class='left10'>" + data + "</span>";
            }
            },
            {"data": "StartTime", "sWidth": "8%"},
            {"data": "EndTime", "sWidth": "7%"},
            {
                "data": "JobId", "sWidth": "9%", "mRender": function (data, type,full) {
                var r = '<button type="button"  title="流水线日志" onClick="jobnameLog(' + data + ',\''+full["JobName"]+'\')" class="delete-groups btn btn-xs rb-btn-oper left10"><i class="fa fa-hospital-o"></i></button>&nbsp;' +
                    '<button type="button"  title="Dockerfile" onClick="showPipelineDockerfile(\'' + full["JobName"] + '\')" class="delete-groups btn btn-xs rb-btn-oper "><i class=" mdi mdi-file-tree"></i></button>&nbsp;'
                return r;
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}

/**
 * 2018-02-05 16:45
 * 获取构建的dockerfile
 * @param historyId
 */
function showPipelineDockerfile(jobName) {
    var url = "/api/ci/job/dockerfile/" + jobName;
    var result = get({}, url);
    $("#show_build_dockerfile_html").html(result);
    $("#job_build_title").html(jobName);
    $("#show_dockerfile_html").modal("toggle")
}