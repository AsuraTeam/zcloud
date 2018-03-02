// 添加配置
function addConfigureData(dataId, configureName, configureId) {
    if (!dataId) {
        dataId = 0
    }
    var url = "/application/configure/data/add?ConfigureName="+configureName+"&ConfigureId="+configureId
    var result = post({ClusterName: "{{.data.ClusterName}}", DataId: dataId}, url)
    $("#add_configure_data_html").html(result)
    $("#add_post_html").modal("toggle")
    if(dataId != 0){
        $(".modal-title").html("更新配置项")
        $(".saveConfigure").html("&nbsp;保存配置项")
    }
}

/**
 * 设置删除配置数据的id
 * @param id
 */
function setDeleteDataId(id) {
    $("#delete_configure_id").val(id)
    deleteConfigureDataSwal();
}

/**
 * 删除配置数据弹出框
 */
function deleteConfigureDataSwal() {
    !function ($) {
        "use strict";

        var SweetAlert = function () {
        };
        //examples
        SweetAlert.prototype.init = function () {
            // //Parameter
            // $('.delete-configure').click(function () {
            swal({
                title: '删除该配置项目',
                text: "",
                type: 'warning',
                showCancelButton: true,
                confirmButtonText: '确认删除',
                cancelButtonText: '不删除',
                confirmButtonClass: 'btn btn-success',
                cancelButtonClass: 'btn btn-danger m-l-10',
                buttonsStyling: false
            }).then(function () {
                var result = deleteConfigureData($("#delete_configure_id").val())
                if (result.indexOf("删除成功") != -1) {
                    swal(
                        '删除成功!',
                        result,
                        'success'
                    )
                    setTimeout(function () {
                        window.location.reload();
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
 * 删除配置数据方法
 * @param id
 * @return {*}
 */
function deleteConfigureData(id) {
    var url = "/api/configure/data/" + id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}


/**
 * 保存配置数据
 */
function saveConfigureData(dataId) {
    if (!dataId) {
        dataId = 0
    }
    var data = get_form_data();
    data["DataId"] = parseInt(dataId)
    if (!checkValue(data, "ConfigureName,DataName,Data")) {
        return
    }
    var url = "/api/configure/data";
    var result = post(data, url)
    result = JSON.stringify(result)
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle")
        success(result)
        loadDetailData()
    } else {
        faild(result)
    }
}
