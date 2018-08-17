// 添加发布服务
function addService(serviceId) {
    if (!serviceId) {
        serviceId = 0
    }
    var url = "/ci/service/add";
    var result = get({ServiceId: serviceId}, url);
    $("#add_code_html").html(result);
    $("#add_post_html").modal("toggle")
}

/**
 * 查看top
 * @param serviceId
 */
function showServiceTop(serviceId) {
    var url = "/ci/service/top/" + serviceId;
    var result = get({}, url);
    $("#add_code_html").html(result);
    $("#add_post_html").modal("toggle")
}

/**
 * 2018-02-14 09;50
 * 发布服务弹出窗口
 * @param serviceId
 */
function releaseService(serviceId) {
    var url = "/ci/service/release";
    var result = get({ServiceId: serviceId}, url);
    $("#add_code_html").html(result);
    $("#add_post_html").modal("toggle")
}

/**
 * 2018-02-18 17;20
 * 修改服务弹出窗口
 * @param serviceId
 */
function modifyHistory(historyId) {
    var url = "/ci/service/release";
    var result = get({historyId: historyId}, url);
    $("#add_code_html").html(result);
    $("#add_post_html").modal("toggle")
}


/**
 * 2018-02-16 17:34
 * 滚动更新服务弹出页面
 * @param serviceId
 */
function rollingService(serviceId, version) {
    var url = "/ci/service/rolling/" + serviceId;
    var result = get({version: version}, url);
    if (result == "资源不可用,请发布后操作"){
        faild(result);
        return;
    }
    $("#add_code_html").html(result);
    $("#add_post_html").modal("toggle")
}

/**
 * 修改发布历史数据
 * 2018-02-18 17:21
 * @param historyId
 */
function saveReleaseHistory(historyId) {
    var data = get_form_data();
    if (!checkValue(data, "Description,ReleaseTestUser,ReleaseOnlineType")) {
        return
    }
    if (data["ReleaseOnlineType"] == "需求") {
        if (!checkValue(data, "ReleaseDemandDescription,ReleaseJobPmLink")) {
            return
        }
    }
    if (data["ReleaseOnlineType"] == "项目") {
        if (!checkValue(data, "ReleaseItemDescription,ReleaseJobPmLink")) {
            return
        }
    }
    if (data["ReleaseOnlineType"] == "BUG修复") {
        if (!checkValue(data, "ReleaseBugDescription,ReleaseBugPmLink")) {
            return
        }
    }
    data["HistoryId"] = historyId;
    var url = "/api/ci/service/history/" + historyId;
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        setTimeout(function () {
            loadReleaseHistoryData();
        }, 3000);
    } else {
        faild(result);
    }
}


/**
 * 执行发布操作
 * 2018-02-14 11:36
 * @param serviceId
 */
function saveReleaseService(serviceId) {
    if (!serviceId) {
        serviceId = 0
    }
    var data = get_form_data();
    data["ServiceId"] = parseInt(serviceId);
    if (!checkValue(data, "Description,ReleaseTestUser,ServiceId,ReleaseOnlineType")) {
        return
    }
    if (data["ReleaseOnlineType"] == "需求") {
        if (!checkValue(data, "ReleaseDemandDescription,ReleaseJobPmLink")) {
            return
        }
    }
    if (data["ReleaseOnlineType"] == "项目") {
        if (!checkValue(data, "ReleaseItemDescription,ReleaseJobPmLink")) {
            return
        }
    }
    if (data["ReleaseOnlineType"] == "BUG修复") {
        if (!checkValue(data, "ReleaseBugDescription,ReleaseBugPmLink")) {
            return
        }
    }
    if (!data["ImageName"] || !data["Version"]) {
        return
    }
    data["ImageName"] = data["ImageName"] + ":" + data["Version"];
    var url = "/api/ci/service/release/" + serviceId;
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        setTimeout(function () {
            loadCiServiceData();
        }, 3000);
    } else {
        faild(result);
    }
}

/**
 * 删除发布服务弹出框
 * 2018-02-10 19:09
 */
function deleteServiceSwal(id) {
    Swal("删除发布服务", "warning", "确认操作", "不操作", "成功", "失败", " deleteService(" + id + ")", "loadCiServiceData()");
}

/**
 * 下线服务弹出框
 * 2018-02-14 18.21
 */
function offlineServiceSwal(id, version) {
    Swal("删除并下线服务,操作不可回滚<br>请确定服务稳定并不需要回滚", "warning", "确认操作", "不操作", "成功", "失败", " serviceOffline(" + id + "," + version + ")", "loadCiServiceData()");
}

/**
 * 2018-02-16 7:32
 * 回滚服务弹出框
 * @param id
 */
function rollbackServiceSwal(id) {
    Swal("回滚后,所有访问都回滚到默认,期间不会更新服务,如果老服务更新后,将无法回滚", "warning", "确认操作", "不操作", "成功", "失败", " rollbackService(" + id + ")", "loadCiServiceData()");
}

/**
 * 2018-02-18 15:45
 * 从历史页面回滚弹出
 * @param id
 */
function rollbackServiceHistorySwal(id, images) {
    Swal("回滚到该版本,回滚只能到2周以内的有效版本!<br>回滚到镜像版本为:<br>" + images, "warning", "确认操作", "不操作", "成功", "失败", " rollbackServiceHistory(" + id + ")", "loadCiServiceData()");
}

/**
 * 2018-02-14 21:26
 *  上线服务
 * @param id
 * @param version
 */
function onServiceSwal(id, version) {
    Swal("点击上线后,所有请求会请求到新部署的服务上", "warning", "确认操作", "不操作", "成功", "失败", " serviceOnline(" + id + "," + version + ")", "loadCiServiceData()");
}


/**
 * 2018-02-16 15:07
 * 切入流量，将绿版的部分服务
 * @param id
 * @param version
 */
function startFlowExecSwal(id, percent) {
    var v = $("#range_01").val();
    var msg = "同意后,负载均衡会加入部分新部署的服务,提供给用户访问";
    if (v == "0") {
        msg = "切换0流量,所有绿版服务会重负载均衡移除,负载均衡将保留蓝版服务";
    }
    if (percent == 0) {
        msg = "将清空所有切入的金丝雀流量服务器...";
    } else {
        percent = ""
    }

    Swal(msg, "warning", "确认操作", "不操作", "成功", "失败", " startFlowExec(" + id + "," + percent + ")", "loadCiServiceData()");
}


/**
 * 2018-02-15 10:52
 * 更新篮版本弹出框
 * @param id
 * @param version
 */
function updateBlueServiceSwal(id) {
    Swal("同意后,会将老版本应用更新到新版本,一般用于绿版测试完成后更新", "warning", "确认操作", "不操作", "成功", "失败", "rollingBlueService(" + id + ")", "loadCiServiceData()");
}

/**
 * 加载数据
 * @param key
 */
function loadCiServiceData(key) {
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

    $("#service-data-table").dataTable({
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
            "displayLength": 3, // 默认长度
            "ajax": { // 请求地址
                "url": "/api/ci/service?t=" + new Date().getTime() + "&search=" + key,
                "type": 'get'
            },
            "columns": [ // 数据映射
                {
                    "data": "Entname", "sWidth": "8%", "mRender": function (data, type, full) {
                    return data + "<br><a class='th-top-8'>" + full["ClusterName"] + "</a>";
                }
                },
                {
                    "data": "Domain", "sWidth": "9%", "mRender": function (data, type, full) {
                    return data + "<br><a class='th-top-8'>" + full["GroupName"] + "</a>";
                }
                },
                {
                    "data": "AppName", "sWidth": "7%", "mRender": function (data, type, full) {
                    return data + "<br><a  class='th-top-8' href='/application/service/list?name=" + full["ServiceName"] + "'>" + full["ServiceName"] + "</a>";
                }
                },
                {
                    "data": "ImageInfoGreen", "sWidth": "22%", "mRender": function (data, type, full) {
                    var version = full["LbVersion"];
                    var html = full["ImageName"] + "<br>" + full["ImageInfoBlue"];
                    var offlineGreen = "<span title='该版本当前没有在负载均衡提供服务' class='text-default m-l-15'>离线</span>";
                    var offline = "<span title='该版本当前没有在负载均衡提供服务' class='text-default th-top-8 m-l-15'>离线</span>";
                    var online = "<span title='该版本当前在负载均衡提供服务' class='text-default th-top-8 m-l-15' style='color:#f96a74 !important '>在线</span>";
                    var newVersion = "<span  title='该版本为最新版本' class='m-l-15 text-danger'>新</span>";
                    var newVersionBlue = "<span  title='该版本为最新版本' class='m-l-15 th-top-8 text-danger'>新</span>";
                    var oldVersion = "<span  title='该版本为老版本' class='m-l-15 th-top-8 text-default'>旧</span>";
                    var oldVersionGreen = "<span  title='该版本为老版本' class='m-l-15 text-default'>旧</span>";

                    if (full["NewVersion"] == "1" && full["Percent"] != "0") {
                        html += data + "<span title='金丝雀部署在负载均衡已切入流量百分比' class='m-l-5 text-danger m-r-5'>" + full["Percent"] + "&nbsp;%</span>";
                    }
                    if (version == "1") {
                        html += online;
                    } else {
                        if (full["Percent"] == "0") {
                            html += offline;
                        }
                    }
                    if (parseInt(full["Percent"]) > 0 && full["NewVersion"] == "1") {
                        html += "<span class='m-l-5' onclick='startFlowExecSwal(" + full["ServiceId"] + ",0)' title='将金丝雀的服务器清空,就是金丝雀部署加入的服务器都不接受访问了'><a>清空金丝雀</a></span>";
                    }
                    if (full["NewVersion"] == "1") {
                        html += newVersionBlue;
                    } else {

                        html += oldVersion;
                    }
                    html += "<br><br>";
                    if (data) {
                        html += data;
                        if (full["NewVersion"] == "2" && full["Percent"] != "0") {
                            html += "<span title='金丝雀部署在负载均衡已切入流量百分比' class='m-l-5 text-danger m-r-5'>" + full["Percent"] + "%在线</span>";
                        }
                        if (parseInt(full["Percent"]) > 0 && full["NewVersion"] == "2") {
                            html += "<span class='m-l-5' onclick='startFlowExecSwal(" + full["ServiceId"] + ",0)' title='将金丝雀的服务器清空,就是金丝雀部署加入的服务器都不接受访问了'><a>清空金丝雀</a></span>";
                        }
                        if (version == "2") {
                            html += "<span title='该版本当前在负载均衡提供服务' class='text-default m-l-5' style='color:#f96a74 !important '>在线</span>";
                        } else {
                            if (full["Percent"] == "0") {
                                html += offlineGreen;
                            }
                        }
                        if (full["NewVersion"] == "2") {
                            html += newVersion;
                        } else {
                            html += oldVersionGreen;
                        }
                        return html;
                    }
                    return html;
                }
                },
                {
                    "data": "LbVersion", "sWidth": "13%", "mRender": function (data, type, full) {
                    var html = "";
                    var serviceId = full["ServiceId"];
                    var blue = full["ImageInfoBlue"];
                    var green = full["ImageInfoGreen"];
                    blue = blue.split(":");
                    green = green.split(":");
                    // 蓝绿操作按钮
                    if (green.indexOf("版本不存在") != -1) {
                        var blueGreen = "<span title='将该版本上线,所有流量都接入到新版本,只改动lb,不改动svc,所以使用没有提供内部访问的服务' onclick='onServiceSwal(" + serviceId + ",1)'   class='text-default m-l-10'><a>切换蓝</a></span>";
                    } else {
                        if (full["LbVersion"] == "2") {
                            var blueGreen = "<span title='将该版本上线,所有流量都接入到新版本,只改动lb,不改动svc,所以使用没有提供内部访问的服务' onclick='onServiceSwal(" + serviceId + ",1)'   class='text-default m-l-10'><a>切换蓝</a></span>";
                        } else {
                            var blueGreen = "<span title='当前版本无需切换'  class='text-warning m-l-10'>切换蓝</span>";
                        }
                    }
                    var greenBlue = "<span title='将该版本上线,所有流量都接入到新版本,只改动lb,不改动svc,所以使用没有提供内部访问的服务' onclick='onServiceSwal(" + serviceId + ",2)'   class='text-default th-top-8 m-l-10'><a>切换绿</a></span>";

                    var delBlue = "<span onclick='offlineServiceSwal(" + serviceId + ",1)' class='text-default  m-l-10'><a>删除</a></span>";
                    if (data == "1") {
                        delBlue = "<span  class='text-warning  m-l-10'>删除</span>";
                    }

                    var blueUpdate = delBlue + "<span onclick='rollingService(" + serviceId + ",1)'  title='将蓝版滚动更新到最新版本' class='text-default m-l-10'><a>滚动</a></span>";
                    if (blue && green) {
                        blue = blue[blue.length - 1].replace(/<\/span>/g, "").replace(/<br>/, "");
                        green = green[green.length - 1].replace(/<\/span>/g, "").replace(/<br>/, "");
                    }

                    var rolling = "<span onclick='rollingService(" + serviceId + ",1)'  title='将蓝版滚动更新到最新版本' class='text-default  m-l-10'><a>滚动</a></span>";
                    if (full["NewVersion"] == "1") {
                        rolling = "<span title='最新版本,不能做滚动更新' class='text-warning  m-l-10'>滚动</span>";
                    }
                    if (blue == green) {
                        rolling = "<span onclick='faild(\"蓝绿版本一致,不能进行滚动更新!\")'   title='当前蓝绿版本一致,无法进行滚动更新' class='text-warning m-l-10'>滚动</span>";
                    }
                    // 选择是保留还是蓝绿切换,保留在版本都更新到最新时可用操作
                    if (blue.trim() == "版本不存在") {
                        html += "<span class='text-default' style='color: #4489e4a6;'>蓝<span class='m-l-10'>版本不存在</span></span><br>"
                    } else {
                        html += "<span class='text-default' style='color: #4489e4a6;'>蓝" + blueGreen + delBlue + rolling;
                        if (blue != green) {
                            if (full["NewVersion"] == "1") {
                                html += "<span onclick='startFlow(" + serviceId + ",1)' class='text-default  m-l-10'><a>金丝雀</a></span>";
                            } else {
                                html += "<span title='此版本不是最新版本,或该服务已经提供服务,不能进行金丝雀发布' class='text-warning m-l-10'>金丝雀</span>";
                            }
                        } else {
                            html += "<span title='蓝绿版本一致,无需进行金丝雀发布!' onclick='faild(\"蓝绿版本一致,不能进行金丝雀发布!\")'  class='text-warning  m-l-10'>金丝雀</span>";
                        }
                        html += "<br>";
                    }
                    if (green) {
                        var delcss = "<span onclick='offlineServiceSwal(" + serviceId + ",2)' class='text-default th-top-8 m-l-10'><a>删除</a></span>";
                        if (data == "2") {
                            delcss = "<span  class='text-warning th-top-8 m-l-10'>删除</span>";
                        }
                        var rolling = "<span onclick='rollingService(" + serviceId + ",2)'  title='将绿版滚动更新到最新版本' class='text-default th-top-8 m-l-10'><a>滚动</a></span>";
                        if (full["NewVersion"] == "2") {
                            rolling = "<span title='最新版本,不能做滚动更新' class='text-warning th-top-8 m-l-10'>滚动</span>";
                        }
                        html += "<span class='th-top-8 Running'>绿</span>" + greenBlue + delcss + rolling;

                        if (blue != green) {
                            if (full["NewVersion"] == "2" && full["LbVersion"] != "2") {
                                html += "<span onclick='startFlow(" + serviceId + ",2)' class='text-default th-top-8 m-l-10'><a>金丝雀</a></span>";
                            } else {
                                html += "<span  class='text-warning th-top-8 m-l-10'>金丝雀</span>";
                            }
                        } else {
                            html += "<span title='蓝绿版本一致,无需进行金丝雀发布!' onclick='faild(\"蓝绿版本一致,不能进行金丝雀发布!\")'  class='text-warning th-top-8 m-l-10'>金丝雀</span>";
                        }
                        html += "<br><br><a onclick='showServiceTop(" + serviceId + ")'>显示top</a>";
                        return html
                    } else {
                        html += "<span class='text-default th-top-8 Running'>绿<span class='m-l-10'>版本不存在</span></span><br>"
                    }
                    html += "<br><a onclick='showServiceTop(" + serviceId + ")'>显示top</a>";
                    return html
                }
                },
                {
                    "data": "ServiceId", "sWidth": "8%", "sClass": "tHeight", "mRender": function (data, type, full) {
                    var r = "<span class='m-l-10'>";
                    if ((full["ImageInfoGreen"] && !full["ImageInfoGreen"].indexOf("版本不存在") != -1) && (full["ImageInfoBlue"] && !full["ImageInfoBlue"].indexOf("版本不存在") != -1)) {
                        r += '<button type="button"  disabled title="发布,目前上线未完成,不能发布"  onclick="releaseService(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="mdi mdi-call-split"></i></button>&nbsp;';
                    } else {
                        r += '<button type="button" title="发布"  onclick="releaseService(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="mdi mdi-call-split"></i></button>&nbsp;';
                    }
                    var blue = full["ImageInfoBlue"];
                    var green = full["ImageInfoGreen"];
                    blue = blue.split(":");
                    green = green.split(":");
                    if (blue && green) {
                        // blue = blue[blue.length - 1].replace(/<\/span>/g, "").replace(/<br>/, "");
                        // green = green[green.length - 1].replace(/<\/span>/g, "").replace(/<br>/, "");
                        var rollback = "";
                        if (full["NewVersion"] != "") {
                            rollback = '<button type="button"  title="回滚" onClick="rollbackServiceSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper m-l-5"><i class="mdi mdi-rotate-left"></i></button>&nbsp;';
                        }
                        if (full["NewVersion"] == "") {
                            rollback = '<button type="button" disabled title="新旧版本一致,无法回滚" class="delete-groups btn btn-xs rb-btn-oper m-l-5"><i class="mdi mdi-rotate-left"></i></button>&nbsp;';
                        }
                        if (full["Percent"] != "0") {
                            rollback = '<button type="button" disabled title="请将金丝雀服务清空后进行回滚,点击清空金丝雀按钮" class="delete-groups btn btn-xs rb-btn-oper m-l-5"><i class="mdi mdi-rotate-left"></i></button>&nbsp;';
                        }
                        r += rollback;
                    }
                    r += '<button type="button"  title="到历史页面" onClick="toHistory(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper m-l-5"><i class="fa fa-history"></i></button>';
                    return r + "</span>";
                }
                }
                ,
                {
                    "data": "ServiceId", "sWidth": "6%", "mRender": function (data) {
                    return '<button type="button" title="更新" onclick="addService(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                        '<button type="button"  title="删除" onClick="deleteServiceSwal(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>' +
                        '<button type="button"  title="到操作日志页面" onClick="toLogs(' + data + ')" class="delete-groups btn btn-xs rb-btn-oper m-l-5"><i class="fa fa-history"></i></button>';
                }
                },
            ],
            "fnRowCallback": function (row, data) { // 每行创建完毕的回调
                $(row).data('recordId', data.recordId);
            }
        }
    )
    ;
}

/**
 * 2018-02-15 18:36
 * 跳转到历史页面
 * @param id
 */
function toHistory(id) {
    var url = "/ci/service/release/history?ServiceId=" + id;
    window.location.href = url;
}

/**
 * 2018-02-17 11:46
 * 跳转到操作日志页面
 * @param id
 */
function toLogs(id) {
    var url = "/ci/service/release/logs?ServiceId=" + id;
    window.location.href = url;
}

/**
 * 2018-02-16 10:57
 * 灰度金丝雀切入流量
 * @param id
 */
function startFlow(id, version) {
    var url = "/ci/service/flow/" + id;
    var result = get({version: version}, url);
    $("#add_code_html").html(result);
    $("#add_post_html").modal("toggle")
}

/**
 * 2018-02-14 18:23
 * 删除并下线服务方法
 * @param id
 * @return {*}
 */
function serviceOffline(id, version) {
    var url = "/api/ci/service/release/" + id + "?version=" + version;
    var result = del({}, url);
    result = JSON.stringify(result);
    return result
}

/**
 * 2018-02-14 21:25
 * 上线服务,新服务上线
 * @param id
 * @return {string}
 */
function serviceOnline(id, version) {
    var url = "/api/ci/service/online/" + id;
    var result = post({ServiceVersion: version}, url);
    result = JSON.stringify(result);
    return result
}

/**
 * 2018-02-16 10:51
 * 更新服务,将蓝版更新和绿版一致
 * @param id
 * @return {string}
 */
function updateBlueService(id) {
    var url = "/api/ci/service/blue/" + id;
    var result = post({}, url);
    result = JSON.stringify(result);
    return result
}

/**
 * 删除发布服务方法
 * @param id
 * @return {*}
 */
function deleteService(id) {
    var url = "/api/ci/service/" + id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}

/**
 * 2018-02-15 07:29
 * 回滚服务
 * @param id
 * @return {string}
 */
function rollbackService(id) {
    var url = "/api/ci/service/rollback/" + id;
    var result = post({}, url);
    result = JSON.stringify(result);
    return result
}

/**
 * 2018-02-18 15:48
 * 回滚服务从历史页面
 * @param id
 * @return {string}
 */
function rollbackServiceHistory(id) {
    var url = "/api/ci/service/rollback/" + id;
    var result = post({history: id}, url);
    result = JSON.stringify(result);
    return result
}


/**
 * 金丝雀流量切入操作
 * @param id
 * @return {string}
 */
function startFlowExec(id, percent) {
    if (!percent && percent != 0) {
        percent = $("#range_01").val();
    }
    var version = $('input[name="version"]').val();
    var url = "/api/ci/service/flow/" + id;
    var result = post({percent: percent, version: version}, url);
    result = JSON.stringify(result);
    return result
}

/**
 * 2018-02-16 19;04
 * 滚动更新蓝色服务
 * @param id
 * @return {string}
 */
function rollingBlueService(id) {
    var data = get_form_data();
    var url = "/api/ci/service/rolling/" + id;
    var result = post(data, url);
    result = JSON.stringify(result);
    return result
}

/**
 * 保存发布服务
 */
function saveCiService(serviceId) {
    if (!serviceId) {
        serviceId = 0
    }
    var data = get_form_data();
    data["ServiceId"] = parseInt(serviceId);
    if (!checkValue(data, "Entname,ClusterName,ServiceName,AppName,Domain")) {
        return
    }
    if ($("#selectReleaseType1").is(":checked")) {
        data["ReleaseType"] = 1;
    }
    if ($("#selectReleaseType2").is(":checked")) {
        data["ReleaseType"] = 2;
    }
    if ($("#selectReleaseType3").is(":checked")) {
        data["ReleaseType"] = 3;
    }
    var url = "/api/ci/service";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        loadCiServiceData();
    } else {
        faild(result);
    }
}
