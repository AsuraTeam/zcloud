// 添加直接
function addJob(jobId) {
    if (!jobId) {
        jobId = 0;
    }
    var url = "/ci/job/add";
    var result = get({JobId: jobId}, url);
    $("#add_job_html").html(result);
    $("#add_post_html").modal("toggle")
}

/**
 * 执行构建程序
 * @param jobId
 */
function execJob(jobId) {
    if (!jobId) {
        return;
    }
    var url = "/api/ci/job/exec/" + jobId;
    var result = get({}, url);
    result = JSON.stringify(result);
    setTimeout(function () {
        loadJobData();
    }, 3000);
    setTimeout(function () {
        loadJobData();
    }, 5000);
    jobLog(jobId);
    return result;
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
 * 2018-01-29 8:01
 * 获取构建的dockerfile
 * @param historyId
 */
function showBuildDockerfile(historyId, job) {
    var url = "/api/ci/job/dockerfile/" + historyId;
    var result = get({}, url);
    $("#show_build_dockerfile_html").html(result);
    $("#job_build_title").html(job);
    $("#show_dockerfile_html").modal("toggle")
}

/**
 * 删除job弹出框
 * 2018-01-25 10:24
 */
function deleteJobSwal(id, detail) {
    if (detail) {
        Swal("删除Job", "warning", "确认操作", "不操作", "成功", "失败", " deleteJob(" + id + ")", "loadJobData()");
    } else {
        Swal("删除Job", "warning", "确认操作", "不操作", "成功", "失败", " deleteJob(" + id + ")", "loadJobData()");
    }
}

/**
 * 执行构建弹出框
 * 2018-01-27 5:44
 * @param id
 */
function startExecJob(id) {
    Swal("是否执行该构建", "warning", "确认操作", "不操作", "成功", "失败", " execJob(" + id + ")", "loadJobData()");
}

/**
 * 加载数据
 * @param key
 */
function loadJobData(key) {
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
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/ci/job?t=" + new Date().getTime() + "&search=" + key,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "ItemName", "sWidth": "8%", "mRender": function (data) {
                return "<a href='/ci/job/detail/" + data + "'>" + data + "</a>"
            }
            },
            {"data": "TimeOut", "sWidth": "5%"},
            {"data": "ClusterName", "sWidth": "9%"},
            {
                "data": "LastTag", "sWidth": "8%", "mRender": function (data) {
                return data;
            }
            },
            {"data": "Description", "sWidth": "10%"},
            {
                "data": "BuildStatus", "sWidth": "6%", "mRender": function (data) {
                if (!data) {
                    var r = '<div class="Fail"><div><i class="fa fa-circle"></i><span>&nbsp;未构建</span></div>'
                    return r
                }
                if (data == "构建中") {
                    var r = '<div class="Running"><div><i class="fa fa-circle"></i><span>&nbsp;构建中</span></div>'
                    return r
                }
                if (data == "构建失败") {
                    var r = '<div class="Fail"><div><i class="fa fa-circle"></i><span>&nbsp;构建失败</span></div>'
                    return r
                }
                if (data == "构建成功") {
                    var r = '<div class="Running"><div><i class="fa fa-circle"></i><span>&nbsp;构建成功</span></div>'
                    return r
                }
                return data
            }
            },
            {"data": "LastModifyTime", "sWidth": "9%"},
            {
                "data": "JobId", "sWidth": "8%", "mRender": function (data) {
                return '<button type="button" title="立刻构建" onclick="startExecJob(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-transgender-alt"></i></button>&nbsp;' +
                    '<button type="button" title="最近构建日志" onclick="jobLog(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa  fa-hospital-o"></i></button>&nbsp;' +
                    '<button type="button" title="查看构建仓库组信息" onclick="toRegistryGroup()" class="btn btn-xs rb-btn-oper"><i class="mdi mdi-arrange-send-to-back"></i></button>&nbsp;';
            }
            },
            {
                "data": "JobId", "sWidth": "6%", "mRender": function (data) {
                return '<button type="button" title="更新" onclick="addJob(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                    '<button type="button"  title="删除" onClick="deleteJobSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>&nbsp;' +
                    '<button type="button"  title="构建历史" onClick="buildHistory(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-history"></i></button>';
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}


/**
 * 删除job方法
 * @param id
 * @return {*}
 */
function deleteJob(id) {
    var url = "/api/ci/job/" + id;
    var result = del({}, url);
    result = JSON.stringify(result);
    return result
}


/**
 * 保存job
 */
function saveJob(jobId) {
    if (!jobId) {
        jobId = 0
    }
    var data = get_form_data();
    console.log(data);
    data["JobId"] = parseInt(jobId);
    if (!checkValue(data, "ItemName,RegistryServer,TimeOut,ClusterName,BaseImage")) {
        return
    }
    if (!$("#selectImageTag1").is(":checked")) {
        if (!checkValue(data, "ImageTag")) {
            return
        }
    } else {
        data["ImageTa"] = "000";
    }
    var select = $("#setSelectDockerfile1").is(":checked");
    if (select){
        data["DockerFile"] = $("#select-docker-file").val();
    }else{
        data["DockerFile"] = 0
    }
    var url = "/api/ci/job";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        loadJobData();
    } else {
        faild(result);
    }
}


/**
 * 构建历史
 * @param id
 */
function buildHistory() {
    window.location.href = "/ci/job/history/list";
}


/**
 * 加载数据
 * @param key
 */
function loadJobHistoryData(key) {
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
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/ci/job/history?t=" + new Date().getTime() + "&search=" + key,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "JobName", "sWidth": "17%", "mRender": function (data) {
                return "<a href='/ci/job/detail/" + data + "'>" + data + "</a>"
            }
            },
            {"data": "ItemName", "sWidth": "9%"},
            {
                "data": "ImageTag", "sWidth": "10%", "mRender": function (data) {
                if (data == "000") {
                    return "时间戳";
                }
                return data;
            }
            },
            {"data": "CreateUser", "sWidth": "6%"},
            {
                "data": "BuildTime", "sWidth": "7%", "mRender": function (data) {
                return "<span class='left10'>" + data + "</span>"
            }
            },
            {
                "data": "BuildStatus", "sWidth": "6%", "mRender": function (data) {
                if (!data) {
                    var r = '<div class="Fail"><div><i class="fa fa-circle"></i><span>&nbsp;未构建</span></div>'
                    return r
                }
                if (data == "构建中") {
                    var r = '<div class="Running"><div><i class="fa fa-circle"></i><span>&nbsp;构建中</span></div>'
                    return r
                }
                if (data == "构建失败") {
                    var r = '<div class="Fail"><div><i class="fa fa-circle"></i><span>&nbsp;构建失败</span></div>'
                    return r
                }
                    if (data == "构建超时") {
                        var r = '<div class="Fail"><div><i class="fa fa-circle"></i><span>&nbsp;构建超时</span></div>'
                        return r
                    }
                if (data == "构建成功") {
                    var r = '<div class="Running"><div><i class="fa fa-circle"></i><span>&nbsp;构建成功</span></div>'
                    return r
                }
                return data
            }
            },
            {"data": "CreateTime", "sWidth": "8%"},
            {
                "data": "JobId", "sWidth": "10%", "mRender": function (data, type, full) {
                var r = '<button type="button"  title="构建日志" onClick="jobLog(' + data + ',' + full["HistoryId"] + ')" class="delete-groups btn btn-xs rb-btn-oper left10"><i class="fa fa-hospital-o"></i></button>&nbsp;' +
                    '<button type="button"  title="Dockerfile" onClick="showBuildDockerfile(' + full["HistoryId"] + ',\'' + full["JobName"] + '\')" class="delete-groups btn btn-xs rb-btn-oper "><i class=" mdi mdi-file-tree"></i></button>&nbsp;'
                if (full["BuildStatus"] == "构建成功") {
                    r+='<button type="button"  title="用该版本创建应用" onClick="toCreateApp(' + full["HistoryId"] + ')" class="delete-groups btn btn-xs rb-btn-oper "><i class="mdi mdi-exit-to-app"></i></button>&nbsp;';
                }
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
 * 到创建应用窗口
 * @param histroyId
 */
function toCreateApp(histroyId) {
    window.location.href = '/application/app/add?AppId=0&historyId=' + histroyId;
}

/**
 * 转到镜像仓库组
 * 2018-01-31 22:14
 */
function toRegistryGroup() {
    window.location.href = "/image/registry/group/list"
}