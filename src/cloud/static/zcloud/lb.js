// 添加应用
function addLb(lbId) {
    if (!lbId) {
        lbId = 0
    }
    var url = "/base/network/lb/add"
    var result = get({ClusterName: getClusterName(), LbId: lbId}, url)
    $("#add_lb_html").html(result)
    $("#add_post_html").modal("toggle")
}

/**
 * 到应用详情页面
 * @param name
 */
function toLbDetail(id) {
    var url = "/base/network/lb/detail/" + id;
    window.location.href = url;
}


/**
 * 删除负载均衡方法
 * @param id
 * @return {*}
 */
function deleteLb(id) {
    var url = "/api/lb/" + id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}

// 设置负载均衡的IP地址
function setLbIp() {
    var url = "/api/cluster/nodes"
    var data = get({clusterName: getClusterName()}, url)
    var result = new Array;
    var ip;
    for (var i = 0; i < data.length; i++) {
        ip = data[i]["Ip"]
        if (ip) {
            result.push(ip)
        }
    }
    $("textarea[name='LbIp']").val(result.join(","))
}

/**
 * 删除负载服务
 * 2018-02-01 20:59
 */
function deleteLbServiceSwal(id) {
    Swal("删除该服务", "warning", "确认操作", "不操作", "成功", "失败", " deleteLbService(" + id + ")", "loadLbDetailData()");
}

/**
 * 删除负载均衡服务
 * 2018-02-01 21:01
 * @param id
 * @return {*}
 */
function deleteLbService(id) {
    var url = "/api/network/lb/service/" + id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}

/**
 * 保存资源配额
 */
function saveLbService(lbId) {
    if (!lbId) {
        lbId = 0
    }
    var data = get_form_data();
    data["ServiceId"] = parseInt(lbId);
    if (!checkValue(data, "LbName,ClusterName,ContainerPort,LbServiceId")) {
        return
    }
    data["LbMethod"] = "0";
    data["DefaultDomain"] = "0";
    if ($("input[name='DefaultDomain']").is(":checked")) {
        data["DefaultDomain"] = "on";
    }

    if ($("#LbMethod1").is(":checked")) {
        data["LbMethod"] = "service";
    }

    if ($("#LbMethod2").is(":checked")) {
        data["LbMethod"] = "pod";
    }
    if (data["CertFile"] == "请选择证书") {
        data["CertFile"] = "";
    }

    var url = "/api/network/lb/service";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        loadLbDetailData();
    } else {
        faild(result);
    }
}

/**
 * 保存负载均衡
 */
function saveLb(lbId) {
    if (!lbId) {
        lbId = 0
    }
    var data = get_form_data();
    data["LbId"] = parseInt(lbId)
    if (!checkValue(data, "LbName,ClusterName,Entname,LbDomainSuffix")) {
        return
    }
    var url = "/api/lb";
    var result = post(data, url)
    result = JSON.stringify(result)
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle")
        success(result)
        loadLbData()
    } else {
        faild(result)
    }
}
/**
 * 设置删除负载均衡的id
 * @param id
 */
function setDeletelbId(id) {
    $("#delete_lb_id").val(id)
    deleteLbSwal();
}


function loadLbData(ip) {
    if (!ip) {
        ip = ""
    } else {
        if (ip.length < 4) {
            return
        }
    }

    $("#lb-data-table").dataTable({
        "filter": false,//去掉搜索框
        "ordering": false, // 是否允许排序
        "paginationType": "full_numbers", // 页码类型
        "destroy": true,
        "processing": true,
        "bPaginate": false, //是否显示（应用）分页器
        "serverSide": true,
        "bLengthChange": false,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/lb?t=" + new Date().getTime() + "&ip=" + ip + "&cluster=" + getClusterName(),
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "LbName", "mRender": function (data, type, full) {
                return "<a href='javascript:void(0)' onclick='toLbDetail(\"" + full["LbId"] + "\")'>" + data + "</a>";
            }
            },
            {"data": "LbType"},
            {
                "data": "ClusterName", "mRender": function (data) {
                return "<a href='/base/cluster/detail/" + data + "'>" + data + "</a>"
            }
            },
            {
                "data": "Entname", "mRender": function (data) {
                return data;
            }
            },
            {"data": "CreateTime"},
            {"data": "ServiceNumber"},
            {"data": "LastModifyTime"},
            {
                "sWidth": "150px", "data": "LbId", "mRender": function (data, type, full) {
                if ("master" != full["HostType"]) {
                    return '<button type="button" title="更新" onclick="addLb(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                        '<button type="button"  title="删除" onClick="setDeletelbId(' + data + ')" class="delete-template btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
                }
                return ""
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}

loadLbData();

/**
 * 删除资源配额弹出框
 */
function deleteLbSwal() {
    !function ($) {
        "use strict";

        var SweetAlert = function () {
        };
        //examples
        SweetAlert.prototype.init = function () {
            // //Parameter
            // $('.delete-quota').click(function () {
            swal({
                title: '删除负载均衡',
                text: "",
                type: 'warning',
                showCancelButton: true,
                confirmButtonText: '确认删除',
                cancelButtonText: '不删除',
                confirmButtonClass: 'btn btn-success',
                cancelButtonClass: 'btn btn-danger m-l-10',
                buttonsStyling: false
            }).then(function () {
                var result = deleteLb($("#delete_lb_id").val())
                if (result.indexOf("删除成功") != -1) {
                    swal(
                        '删除成功!',
                        result,
                        'success'
                    )
                    setTimeout(function () {
                        window.location.href = "/base/network/lb/list"
                    }, 2000)
                } else {
                    swal(
                        '删除失败!',
                        result,
                        'error'
                    )
                }
            }, function (dismiss) {
                $("#delete_lb_id").val("")
            })
            // });
        },
            $.SweetAlert = new SweetAlert, $.SweetAlert.Constructor = SweetAlert
    }(window.jQuery),

//initializing
        function ($) {
            "use strict";
            $.SweetAlert.init()
        }(window.jQuery);
}

/*
 /**
 * Created by zhaoyun on 2018/1/5.
 */

