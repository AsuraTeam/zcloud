
// 添加直接
function addTemplate(templateId) {
    if(!templateId){
        templateId = 0
    }
    var url = "/application/template/add"
    var result = post({ClusterName: "{{.data.ClusterName}}", TemplateId:templateId}, url)
    $("#add_template_html").html(result)
    $("#add_post_html").modal("toggle")
}

/**
 * 设置删除模板的id
 * @param id
 */
function setDeleteId(id) {
    $("#delete_template_id").val(id)
    deleteTemplateSwal();
}

/**
 * 删除模板弹出框
 */
function deleteTemplateSwal() {
    !function ($) {
        "use strict";

        var SweetAlert = function () {
        };
        //examples
        SweetAlert.prototype.init = function () {
            // //Parameter
            // $('.delete-template').click(function () {
                swal({
                    title: '删除该模板',
                    text: "",
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonText: '确认删除',
                    cancelButtonText: '不删除',
                    confirmButtonClass: 'btn btn-success',
                    cancelButtonClass: 'btn btn-danger m-l-10',
                    buttonsStyling: false
                }).then(function () {
                    var result = deleteTemplate($("#delete_template_id").val())
                    if (result.indexOf("删除成功") != -1) {
                        swal(
                            '删除成功!',
                            result,
                            'success'
                        )
                        setTimeout(function () {
                            window.location.href = "/application/template/list"
                        },2000)
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
                    $("#delete_template_id").val("")
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
function loadTemplateData(key) {
    if (!key) {
        key = ""
    } else {
        if (key.length < 4) {
            return
        }
    }

    $("#template-data-table").dataTable({
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
            "url": "/api/template?t=" + new Date().getTime() + "&search=" + key + "&cluster="+getClusterName(),
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "TemplateName"},
            {"data": "ResourceName"},
            {"data": "Description"},
            {"data": "CreateTime"},
            {"data": "LastModifyTime"},
            {
                "sWidth": "150px", "data": "TemplateId", "mRender": function (data, type, full) {

                    return '<button type="button" title="更新" onclick="addTemplate(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-pencil"></i></button>&nbsp;' +
                        '<button type="button"  title="删除" onClick="setDeleteId(' + data + ')" class="delete-template btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}

loadTemplateData();

/**
 * 删除模板方法
 * @param id
 * @return {*}
 */
function deleteTemplate(id) {
    var url = "/api/template/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
}


/**
 * yaml  部署模板
 */
function setTemplate() {
    if (!$("textarea[name='Yaml']").val()) {
        $("textarea[name='Yaml']").val($("#yaml_template").val());
        checkChange('Yaml','textarea');
    }
}

/**
 * 保存模板
 */
function saveTemplate(templateId) {
    if(!templateId){
        templateId = 0
    }
    var data = get_form_data();
    data["TemplateId"] = parseInt(templateId)
    if(!checkValue(data,"TemplateName,Yaml,ResourceName")){
        return
    }
    if(!checkYaml(data)){
        setInputError($("textarea[name='Yaml']"), "errmsg")
        return
    }
    var url = "/api/template";
    var result = post(data, url)
    result = JSON.stringify(result)
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle")
        success(result)
        loadTemplateData()
    }else{
        faild(result)
    }
}
