// 添加配置
function addConfigure(configureId) {
    if (!configureId) {
        configureId = 0
    }
    var url = "/application/configure/add"
    var result = get({ClusterName: "{{.data.ClusterName}}", ConfigureId: configureId}, url)
    $("#add_configure_html").html(result)
    $("#add_post_html").modal("toggle")
    if(configureId != 0){
        $(".modal-title").html("更新配置")
        $(".modal-saveConfigure").html("保存")
    }
}

/**
 * 设置删除模板的id
 * @param id
 */
function setDeleteConfigureId(id) {
    $("#delete_configure_id").val(id)
    deleteConfigureSwal();
}

/**
 * 删除模板弹出框
 */
function deleteConfigureSwal() {
    !function ($) {
        "use strict";

        var SweetAlert = function () {
        };
        //examples
        SweetAlert.prototype.init = function () {
            // //Parameter
            // $('.delete-configure').click(function () {
            swal({
                title: '删除该配置',
                text: "",
                type: 'warning',
                showCancelButton: true,
                confirmButtonText: '确认删除',
                cancelButtonText: '不删除',
                confirmButtonClass: 'btn btn-success',
                cancelButtonClass: 'btn btn-danger m-l-10',
                buttonsStyling: false
            }).then(function () {
                var result = deleteConfigure($("#delete_configure_id").val())
                if (result.indexOf("删除成功") != -1) {
                    swal(
                        '删除成功!',
                        result,
                        'success'
                    )
                    setTimeout(function () {
                        window.location.href = "/application/configure/list"
                    }, 2000)
                } else {
                    swal(
                        '删除失败!',
                        result,
                        'error'
                    )
                }

            }, function (dismiss) {
                // dismiss can be 'cancel', 'overlay',
                // 'close', and 'timer'
                $("#delete_configure_id").val("")
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


/**
 * 加载数据
 * @param key
 */
function loadConfigureData(key) {
    if (!key) {
        key = ""
    } else {
        if (key.length < 4) {
            return
        }
    }

    $("#configure-data-table").dataTable({
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
            "url": "/api/configure?t=" + new Date().getTime() + "&search=" + key + "&cluster=" + getClusterName(),
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "Entname"},
            {"data": "ConfigureName","mRender":function (data) {
                return "<a href='/application/configure/detail/"+data+"' target='_self'>"+data+"</a>"
            }},
            {"data": "ClusterName"},
            {"data": "Description"},
            {"data": "CreateTime"},
            {"data": "LastModifyTime"},
            {
                "sWidth": "150px", "data": "ConfigureId", "mRender": function (data) {

                return '<button type="button" title="更新" onclick="addConfigure(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                    '<button type="button"  title="删除" onClick="setDeleteConfigureId(' + data + ')" class="delete-configure btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';

                return ""
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}

loadConfigureData();

/**
 * 删除模板方法
 * @param id
 * @return {*}
 */
function deleteConfigure(id) {
    var url = "/api/configure/" + id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}


/**
 * 保存模板
 */
function saveConfigure(configureId) {
    if (!configureId) {
        configureId = 0
    }
    var data = get_form_data();
    data["ConfigureId"] = parseInt(configureId)
    if (!checkValue(data, "ConfigureName,Description")) {
        return
    }
    var url = "/api/configure";
    var result = post(data, url)
    result = JSON.stringify(result)
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle")
        success(result)
        loadConfigureData()
    } else {
        faild(result)
    }
}
